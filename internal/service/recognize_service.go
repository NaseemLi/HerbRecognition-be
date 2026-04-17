package service

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"os"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/pkg/onnx"
	"herb-recognition-be/pkg/upload"
)

// RecognizeService 识别服务
type RecognizeService struct{}

// UploadImage 上传图片（带安全校验）
func (s *RecognizeService) UploadImage(file *multipart.FileHeader) (string, error) {
	cfg := upload.DefaultImageConfig
	cfg.UploadDir = "./uploads/images"
	cfg.URLPrefix = "/uploads/images/"
	return upload.UploadFile(file, cfg)
}

// RecognizeRequest 识别请求
type RecognizeRequest struct {
	ImageURL string `json:"image_url"`
}

// RecognizeResponse 识别响应
type RecognizeResponse struct {
	RecordID   uint    `json:"record_id"`
	HerbID     uint    `json:"herb_id"`
	HerbName   string  `json:"herb_name"`
	Confidence float32 `json:"confidence"`
	ImageURL   string  `json:"image_url"`
}

// Recognize 识别图片（使用 ONNX 模型）
func (s *RecognizeService) Recognize(userID uint, imageURL string) (*RecognizeResponse, error) {
	record := model.RecognitionRecord{
		UserID:   userID,
		ImageURL: imageURL,
		Status:   1,
	}

	// 检查 ONNX 预测器是否已初始化
	if !onnx.IsInitialized() {
		record.Status = 0
		record.ErrMsg = "识别模型未初始化"
		record.HerbName = "未知"
		record.Confidence = 0

		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}
		return nil, errors.New("识别模型未初始化，请检查模型文件")
	}

	// 加载图像
	filePath := "." + imageURL
	imgFile, err := os.Open(filePath)
	if err != nil {
		record.Status = 0
		record.ErrMsg = "读取图片失败"
		record.HerbName = "未知"
		record.Confidence = 0

		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}
		return nil, errors.New("读取图片失败")
	}
	defer imgFile.Close()

	// 解码图像
	img, _, err := image.Decode(imgFile)
	if err != nil {
		record.Status = 0
		record.ErrMsg = "解码图片失败"
		record.HerbName = "未知"
		record.Confidence = 0

		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}
		return nil, errors.New("解码图片失败")
	}

	// 执行 ONNX 推理
	result, err := onnx.Predict(img)
	if err != nil {
		record.Status = 0
		record.ErrMsg = fmt.Sprintf("识别失败：%v", err)
		record.HerbName = "未知"
		record.Confidence = 0

		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}
		return nil, errors.New(record.ErrMsg)
	}

	topResult := result.TopResult
	herbID := uint(topResult.HerbID)
	record.HerbID = &herbID
	record.HerbName = topResult.HerbName
	record.Confidence = float32(topResult.Confidence)

	if err := repository.DB.Create(&record).Error; err != nil {
		return nil, fmt.Errorf("保存识别记录失败：%v", err)
	}

	return &RecognizeResponse{
		RecordID:   record.ID,
		HerbID:     herbID,
		HerbName:   topResult.HerbName,
		Confidence: float32(topResult.Confidence),
		ImageURL:   imageURL,
	}, nil
}

// GetHistory 获取用户识别历史
type HistoryQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

func (s *RecognizeService) GetHistory(userID uint, query *HistoryQuery) ([]model.RecognitionRecord, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 10
	}

	offset := (query.Page - 1) * query.PageSize

	var records []model.RecognitionRecord
	var total int64

	if err := repository.DB.Model(&model.RecognitionRecord{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	if err := repository.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&records).Error; err != nil {
		return nil, 0, errors.New("查询失败")
	}

	return records, total, nil
}

// DeleteHistory 删除识别记录
func (s *RecognizeService) DeleteHistory(userID, recordID uint) error {
	var record model.RecognitionRecord
	if err := repository.DB.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error; err != nil {
		return errors.New("记录不存在")
	}

	// 删除图片文件
	upload.DeleteFile(record.ImageURL)

	return repository.DB.Delete(&record).Error
}
