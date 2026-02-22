package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"herb-recognition-be/internal/client"
	"herb-recognition-be/internal/model"
	"herb-recognition-be/internal/repository"

	"github.com/google/uuid"
)

var pythonServiceClient *client.PythonServiceClient

// 模型默认路径 (相对路径：项目根目录/models/best_herb_model.pth)
const DefaultModelPath = "./models/best_herb_model.pth"

// 初始化 Python 服务客户端
// 环境变量 PYTHON_SERVICE_URL 优先级高于默认值
// 服务位于：services/inference-service/
func init() {
	pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
	if pythonServiceURL == "" {
		pythonServiceURL = "http://localhost:5001"
	}
	pythonServiceClient = client.NewPythonServiceClient(pythonServiceURL)
}

// RecognizeService 识别服务
type RecognizeService struct{}

// 允许的图片扩展名
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 允许的图片 MIME 类型
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// 最大文件大小 5MB
const maxFileSize = 5 * 1024 * 1024

// UploadImage 上传图片（带安全校验）
func (s *RecognizeService) UploadImage(file *multipart.FileHeader) (string, error) {
	// 检查文件大小
	if file.Size > maxFileSize {
		return "", errors.New("文件大小不能超过 5MB")
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", errors.New("不支持的图片格式，仅支持 JPG、PNG、GIF、WEBP")
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("文件打开失败")
	}
	defer src.Close()

	// 读取文件头进行 MIME 类型校验
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && err != io.EOF {
		return "", errors.New("文件读取失败")
	}
	buf = buf[:n]

	// 检测真实 MIME 类型
	contentType := http.DetectContentType(buf)
	if !allowedMimeTypes[contentType] {
		return "", errors.New("文件类型不匹配，请上传有效的图片")
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := "./uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", errors.New("创建上传目录失败")
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("文件保存失败")
	}
	defer dst.Close()

	// 重置文件指针并复制内容
	src.Seek(0, 0)
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("文件保存失败")
	}

	// 返回访问 URL
	url := "/uploads/images/" + filename
	return url, nil
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

// Recognize 识别图片（调用模型并保存记录）
func (s *RecognizeService) Recognize(userID uint, imageURL string) (*RecognizeResponse, error) {
	record := model.RecognitionRecord{
		UserID:   userID,
		ImageURL: imageURL,
		Status:   1,
	}

	filePath := "." + imageURL
	imageBytes, err := os.ReadFile(filePath)
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

	filename := filepath.Base(imageURL)
	result, err := pythonServiceClient.RecognizeImage(imageBytes, filename)

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

	herbID := uint(result.HerbID)
	record.HerbID = &herbID
	record.HerbName = result.HerbName
	record.Confidence = float32(result.Confidence)

	if err := repository.DB.Create(&record).Error; err != nil {
		return nil, fmt.Errorf("保存识别记录失败：%v", err)
	}

	return &RecognizeResponse{
		RecordID:   record.ID,
		HerbID:     herbID,
		HerbName:   result.HerbName,
		Confidence: float32(result.Confidence),
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
	if record.ImageURL != "" {
		filePath := "." + record.ImageURL
		os.Remove(filePath)
	}

	return repository.DB.Delete(&record).Error
}
