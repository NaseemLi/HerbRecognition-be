package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"
	"herb-recognition-be/pkg/onnx"
	"herb-recognition-be/pkg/upload"

	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
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
	RecordID   *uint   `json:"record_id"`
	HerbID     uint    `json:"herb_id"`
	HerbName   string  `json:"herb_name"`
	Confidence float32 `json:"confidence"`
	ImageURL   string  `json:"image_url"`
}

// recognizeImage 对 image.Image 执行 ONNX 推理和药材查询
// 返回: (ONNX识别输出, 药材名称, 药材ID, 错误)
func recognizeImage(img image.Image) (*onnx.RecognitionOutput, string, uint, error) {
	if !onnx.IsInitialized() {
		return nil, "", 0, errors.New("识别模型未初始化，请检查模型文件")
	}

	result, err := onnx.Predict(img)
	if err != nil {
		return nil, "", 0, fmt.Errorf("识别失败：%v", err)
	}

	topResult := result.TopResult
	herbName := strings.TrimSpace(topResult.HerbName)
	if herbName == "" {
		herbName = "未知"
	}

	var herbID uint
	if herbName != "未知" {
		var herbs []model.Herb
		err := repository.DB.Select("id", "name").Where("name = ?", herbName).Limit(1).Find(&herbs).Error
		if err != nil {
			return nil, "", 0, fmt.Errorf("查询药材详情失败：%v", err)
		}
		if len(herbs) > 0 {
			herbID = herbs[0].ID
			herbName = herbs[0].Name
		}
	}

	return result, herbName, herbID, nil
}

// Recognize 识别图片（使用 ONNX 模型）
func (s *RecognizeService) Recognize(userID uint, imageURL string) (*RecognizeResponse, error) {
	record := model.RecognitionRecord{
		UserID:   userID,
		ImageURL: imageURL,
		Status:   1,
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

	// 执行识别
	result, herbName, herbID, err := recognizeImage(img)
	if err != nil {
		record.Status = 0
		record.ErrMsg = err.Error()
		record.HerbName = "未知"
		record.Confidence = 0
		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}
		return nil, err
	}

	record.HerbName = herbName
	record.Confidence = float32(result.TopResult.Confidence)
	if herbID > 0 {
		record.HerbID = &herbID
	}

	if err := repository.DB.Create(&record).Error; err != nil {
		return nil, fmt.Errorf("保存识别记录失败：%v", err)
	}

	return &RecognizeResponse{
		RecordID:   &record.ID,
		HerbID:     herbID,
		HerbName:   herbName,
		Confidence: float32(result.TopResult.Confidence),
		ImageURL:   imageURL,
	}, nil
}

// Base64RecognizeRequest base64 识别请求
type Base64RecognizeRequest struct {
	ImageBase64 string `json:"image_base64" binding:"required"`
	SaveHistory *bool  `json:"save_history,omitempty"`
}

// RecognizeFromBase64 从 base64 图片数据进行识别
func (s *RecognizeService) RecognizeFromBase64(userID uint, base64Str string, saveHistory bool) (*RecognizeResponse, error) {
	// 解析 base64 前缀，提取纯数据部分
	base64Data := base64Str
	if idx := strings.Index(base64Str, ","); idx != -1 {
		base64Data = base64Str[idx+1:]
	}

	// base64 解码
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, errors.New("图片数据解码失败")
	}

	// 校验图片大小（最大 5MB）
	if len(data) > 5*1024*1024 {
		return nil, errors.New("图片大小不能超过 5MB")
	}

	// 检测真实 MIME 类型
	contentType := http.DetectContentType(data)
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	default:
		return nil, errors.New("不支持的图片格式，仅支持 jpg、png、gif、webp")
	}

	if saveHistory {
		// 先解码图片做识别，获取 herb_id 用于去重判断
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, errors.New("解码图片失败")
		}

		result, herbName, herbID, err := recognizeImage(img)
		if err != nil {
			return nil, err
		}

		// 去重：查询最近 30 秒内是否已保存过相同 herb_id 的记录
		if herbID > 0 {
			var recentRecords []model.RecognitionRecord
			if err := repository.DB.Where("user_id = ? AND herb_id = ? AND created_at >= ?",
				userID, herbID, time.Now().Add(-30*time.Second)).
				Order("created_at DESC").
				Limit(1).
				Find(&recentRecords).Error; err != nil {
				return nil, fmt.Errorf("查询历史记录失败：%v", err)
			}
			if len(recentRecords) > 0 {
				// 命中去重：直接返回已有记录，不保存新图片
				return &RecognizeResponse{
					RecordID:   &recentRecords[0].ID,
					HerbID:     herbID,
					HerbName:   herbName,
					Confidence: float32(result.TopResult.Confidence),
					ImageURL:   recentRecords[0].ImageURL,
				}, nil
			}
		}

		// 未命中去重：保存图片并创建新记录
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		uploadDir := "./uploads/images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, errors.New("创建上传目录失败")
		}
		filePath := filepath.Join(uploadDir, filename)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, errors.New("保存图片失败")
		}

		imageURL := "/uploads/images/" + filename

		record := model.RecognitionRecord{
			UserID:     userID,
			ImageURL:   imageURL,
			Status:     1,
			HerbID:     &herbID,
			HerbName:   herbName,
			Confidence: float32(result.TopResult.Confidence),
		}
		if err := repository.DB.Create(&record).Error; err != nil {
			return nil, fmt.Errorf("保存识别记录失败：%v", err)
		}

		return &RecognizeResponse{
			RecordID:   &record.ID,
			HerbID:     herbID,
			HerbName:   herbName,
			Confidence: float32(result.TopResult.Confidence),
			ImageURL:   imageURL,
		}, nil
	}

	// saveHistory == false: 内存中解码和识别，不保存图片和记录
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("解码图片失败")
	}

	result, herbName, herbID, err := recognizeImage(img)
	if err != nil {
		return nil, err
	}

	return &RecognizeResponse{
		RecordID:   nil,
		HerbID:     herbID,
		HerbName:   herbName,
		Confidence: float32(result.TopResult.Confidence),
		ImageURL:   "",
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
