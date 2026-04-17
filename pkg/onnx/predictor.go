package onnx

import (
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"sort"

	"herb-recognition-be/pkg/logger"

	"github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

// Predictor ONNX 预测器
type Predictor struct {
	session *onnxruntime_go.AdvancedSession
	input   *onnxruntime_go.Tensor[float32]
	output  *onnxruntime_go.Tensor[float32]
	classes []string
	imgSize int
}

// PredictionResult 预测结果
type PredictionResult struct {
	HerbName   string  `json:"herb_name"`
	HerbID     int     `json:"herb_id"`
	Confidence float64 `json:"confidence"`
}

// RecognitionOutput 识别输出
type RecognitionOutput struct {
	TopResult  PredictionResult   `json:"top_result"`
	AllResults []PredictionResult `json:"all_results"`
}

const (
	inputSize           = 224
	resizeSize          = 255 // int(224 * 1.14)
	inputName           = "input"
	outputName          = "logits"
	confidenceThreshold = 0.3
)

var (
	mean = [3]float32{0.485, 0.456, 0.406}
	std  = [3]float32{0.229, 0.224, 0.225}

	predictor *Predictor
)

// findLibraryPath 查找 ONNX Runtime 库文件
func findLibraryPath() string {
	libNames := []string{
		"libonnxruntime.dylib",
		"libonnxruntime.so",
		"onnxruntime.dll",
	}

	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)

	searchPaths := []string{
		exeDir,
		filepath.Join(exeDir, ".."),
		filepath.Join(exeDir, "..", ".."),
		".",
		"./models/onnx",  // 新增：模型目录
		"./lib",
		"/usr/local/lib",
		"/opt/homebrew/lib",
	}

	for _, dir := range searchPaths {
		for _, name := range libNames {
			path := filepath.Join(dir, name)
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				return path
			}
		}
	}

	return ""
}

// InitPredictor 初始化 ONNX 预测器
func InitPredictor(modelPath, classesPath string) error {
	libPath := findLibraryPath()
	if libPath != "" {
		onnxruntime_go.SetSharedLibraryPath(libPath)
		logger.Infof("找到 ONNX Runtime 库: %s", libPath)
	} else {
		logger.Warn("未找到 libonnxruntime 动态库，将尝试系统默认路径")
	}

	if err := onnxruntime_go.InitializeEnvironment(); err != nil {
		return fmt.Errorf("初始化 ONNX Runtime 环境失败: %v", err)
	}

	classes, err := loadClasses(classesPath)
	if err != nil {
		return fmt.Errorf("加载类别文件失败: %v", err)
	}

	if _, err := os.Stat(modelPath); err != nil {
		return fmt.Errorf("模型文件不存在: %s", modelPath)
	}

	numClasses := len(classes)

	// 输入 Tensor: [1, 3, 224, 224]
	inputShape := onnxruntime_go.NewShape(1, 3, inputSize, inputSize)
	inputData := make([]float32, 1*3*inputSize*inputSize)
	inputTensor, err := onnxruntime_go.NewTensor(inputShape, inputData)
	if err != nil {
		return fmt.Errorf("创建输入 tensor 失败: %v", err)
	}

	// 输出 Tensor: [1, num_classes] - 使用 NewEmptyTensor
	outputShape := onnxruntime_go.NewShape(1, int64(numClasses))
	outputTensor, err := onnxruntime_go.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return fmt.Errorf("创建输出 tensor 失败: %v", err)
	}

	// 输入名称 "input"，输出名称 "logits"
	session, err := onnxruntime_go.NewAdvancedSession(
		modelPath,
		[]string{inputName},
		[]string{outputName},
		[]onnxruntime_go.Value{inputTensor},
		[]onnxruntime_go.Value{outputTensor},
		nil,
	)
	if err != nil {
		return fmt.Errorf("创建 ONNX 会话失败: %v", err)
	}

	predictor = &Predictor{
		session: session,
		input:   inputTensor,
		output:  outputTensor,
		classes: classes,
		imgSize: inputSize,
	}

	return nil
}

// loadClasses 加载类别文件
func loadClasses(path string) ([]string, error) {
	possiblePaths := []string{
		path,
		"./models/onnx/classes.txt",  // 新增
		"./classes.txt",
		"../classes.txt",
		"../../classes.txt",
	}

	var data []byte
	var err error
	for _, p := range possiblePaths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var classes []string
	for _, line := range splitLines(string(data)) {
		if line != "" {
			// 处理格式，提取名称部分
			parts := splitLine(line, '\t')
			if len(parts) >= 2 {
				// 取最后一部分作为名称（去掉编号）
				classes = append(classes, parts[len(parts)-1])
			} else {
				// 没有制表符，直接添加
				classes = append(classes, line)
			}
		}
	}
	return classes, nil
}

// splitLine 按分隔符分割字符串
func splitLine(s string, sep rune) []string {
	var parts []string
	start := 0
	for i, c := range s {
		if c == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// Predict 执行预测
func Predict(img image.Image) (*RecognitionOutput, error) {
	if predictor == nil {
		return nil, fmt.Errorf("预测器未初始化")
	}

	// 预处理图像
	inputData := predictor.preprocess(img)

	// 复制数据到输入 tensor
	tensorData := predictor.input.GetData()
	for i := 0; i < len(inputData) && i < len(tensorData); i++ {
		tensorData[i] = inputData[i]
	}

	// 执行推理
	if err := predictor.session.Run(); err != nil {
		return nil, fmt.Errorf("推理失败: %v", err)
	}

	// 解析输出
	return predictor.postprocess(), nil
}

// preprocess 图像预处理
func (p *Predictor) preprocess(img image.Image) []float32 {
	// 1. Resize(255) 保持宽高比
	resized := resizeKeepRatio(img, resizeSize)

	// 2. CenterCrop(224)
	cropped := centerCrop(resized, inputSize)

	// 3. ToTensor + Normalize
	return imageToCHWFloat32(cropped)
}

// resizeKeepRatio 保持宽高比 resize，短边为目标尺寸
func resizeKeepRatio(img image.Image, targetShort int) image.Image {
	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()

	if srcW <= 0 || srcH <= 0 {
		return img
	}

	var newW, newH int
	if srcW < srcH {
		newW = targetShort
		newH = int(float64(srcH) * float64(targetShort) / float64(srcW))
	} else {
		newH = targetShort
		newW = int(float64(srcW) * float64(targetShort) / float64(srcH))
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

// centerCrop 中心裁剪
func centerCrop(img image.Image, cropSize int) image.Image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	if w < cropSize || h < cropSize {
		return img
	}

	left := (w - cropSize) / 2
	top := (h - cropSize) / 2

	dst := image.NewRGBA(image.Rect(0, 0, cropSize, cropSize))
	draw.Draw(dst, dst.Bounds(), img, image.Point{X: left, Y: top}, draw.Src)
	return dst
}

// imageToCHWFloat32 转换为 CHW 格式并归一化
func imageToCHWFloat32(img image.Image) []float32 {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	data := make([]float32, 3*w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// RGBA() 返回 0~65535，转成 0~1
			rf := float32(r) / 65535.0
			gf := float32(g) / 65535.0
			bf := float32(b) / 65535.0

			// Normalize
			rf = (rf - mean[0]) / std[0]
			gf = (gf - mean[1]) / std[1]
			bf = (bf - mean[2]) / std[2]

			// CHW 格式
			idx := y*w + x
			data[idx] = rf
			data[w*h+idx] = gf
			data[2*w*h+idx] = bf
		}
	}

	return data
}

// postprocess 后处理输出
func (p *Predictor) postprocess() *RecognitionOutput {
	logits := p.output.GetData()

	// Softmax
	probs := softmax(logits)

	// Top 5
	top5 := getTop5(logits, probs, p.classes)

	return &RecognitionOutput{
		TopResult:  top5[0],
		AllResults: top5,
	}
}

// softmax
func softmax(logits []float32) []float64 {
	if len(logits) == 0 {
		return nil
	}

	maxLogit := logits[0]
	for _, v := range logits {
		if v > maxLogit {
			maxLogit = v
		}
	}

	exps := make([]float64, len(logits))
	var sum float64
	for i, v := range logits {
		e := math.Exp(float64(v - maxLogit))
		exps[i] = e
		sum += e
	}

	probs := make([]float64, len(logits))
	for i := range exps {
		probs[i] = exps[i] / sum
	}
	return probs
}

// getTop5 获取 Top 5 结果
func getTop5(logits []float32, probs []float64, classes []string) []PredictionResult {
	type item struct {
		index int
		prob  float64
		logit float32
	}

	items := make([]item, len(logits))
	for i := range logits {
		items[i] = item{i, probs[i], logits[i]}
	}

	// 按概率排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].prob > items[j].prob
	})

	results := []PredictionResult{}
	k := 5
	if k > len(items) {
		k = len(items)
	}

	for i := 0; i < k; i++ {
		if items[i].prob > confidenceThreshold {
			results = append(results, PredictionResult{
				HerbName:   classes[items[i].index],
				HerbID:     items[i].index + 1,
				Confidence: round(items[i].prob*100, 2),
			})
		}
	}

	return results
}

// round 四舍五入
func round(x float64, decimals int) float64 {
	factor := 1.0
	for i := 0; i < decimals; i++ {
		factor *= 10
	}
	return float64(int(x*factor+0.5)) / factor
}

// IsInitialized 检查预测器是否已初始化
func IsInitialized() bool {
	return predictor != nil
}

// GetClassCount 获取类别数量
func GetClassCount() int {
	if predictor == nil {
		return 0
	}
	return len(predictor.classes)
}

// Close 关闭预测器并释放资源
func Close() {
	if predictor != nil {
		if predictor.session != nil {
			predictor.session.Destroy()
		}
		if predictor.input != nil {
			predictor.input.Destroy()
		}
		if predictor.output != nil {
			predictor.output.Destroy()
		}
		predictor = nil
	}
	onnxruntime_go.DestroyEnvironment()
}
