package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Config 上传配置
type Config struct {
	MaxFileSize    int64    // 最大文件大小（字节）
	AllowedExts    []string // 允许的扩展名
	AllowedMimes   []string // 允许的 MIME 类型
	UploadDir      string   // 上传目录
	URLPrefix      string   // URL 前缀
}

// DefaultImageConfig 默认图片上传配置
var DefaultImageConfig = Config{
	MaxFileSize:  5 * 1024 * 1024, // 5MB
	AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	AllowedMimes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
}

// AvatarConfig 头像上传配置
var AvatarConfig = Config{
	MaxFileSize:  2 * 1024 * 1024, // 2MB
	AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	AllowedMimes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
}

// UploadFile 上传文件
func UploadFile(file *multipart.FileHeader, cfg Config) (string, error) {
	// 检查文件大小
	if cfg.MaxFileSize > 0 && file.Size > cfg.MaxFileSize {
		return "", fmt.Errorf("文件大小不能超过 %dMB", cfg.MaxFileSize/1024/1024)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if len(cfg.AllowedExts) > 0 {
		allowed := false
		for _, e := range cfg.AllowedExts {
			if e == ext {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("不支持的文件格式，仅支持 %s", strings.Join(cfg.AllowedExts, ", "))
		}
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
	if len(cfg.AllowedMimes) > 0 {
		allowed := false
		for _, m := range cfg.AllowedMimes {
			if m == contentType {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", errors.New("文件类型不匹配，请上传有效的文件")
		}
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := cfg.UploadDir
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
	return cfg.URLPrefix + filename, nil
}

// DeleteFile 删除文件
func DeleteFile(url string) error {
	if url == "" {
		return nil
	}
	filePath := "." + url
	return os.Remove(filePath)
}