package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

var (
	ErrServiceUnavailable = errors.New("识别服务不可用")
	ErrRecognitionFailed  = errors.New("识别失败")
)

type RecognitionResult struct {
	HerbName   string  `json:"herb_name"`
	HerbID     int     `json:"herb_id"`
	Confidence float64 `json:"confidence"`
}

type RecognitionResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    RecognitionResult `json:"data"`
	Error   string            `json:"error,omitempty"`
	Code    string            `json:"code,omitempty"`
}

type PythonServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPythonServiceClient(baseURL string) *PythonServiceClient {
	return &PythonServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *PythonServiceClient) RecognizeImage(imageBytes []byte, filename string) (*RecognitionResult, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败：%v", err)
	}

	if _, err := part.Write(imageBytes); err != nil {
		return nil, fmt.Errorf("写入图像数据失败：%v", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭 writer 失败：%v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/predict", body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求识别服务失败：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, ErrServiceUnavailable
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败：%v", err)
	}

	var result RecognitionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败：%v", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("%w: %s", ErrRecognitionFailed, result.Message)
	}

	return &result.Data, nil
}

func (c *PythonServiceClient) Health() (map[string]interface{}, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return nil, fmt.Errorf("健康检查失败：%v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败：%v", err)
	}

	var health map[string]interface{}
	if err := json.Unmarshal(respBody, &health); err != nil {
		return nil, fmt.Errorf("解析健康检查响应失败：%v", err)
	}

	return health, nil
}
