# 中草药识别推理服务

基于 MobileNetV3 的中草药图像识别服务。

## 目录结构

```
services/inference-service/
├── app.py                 # Flask 服务主程序
├── requirements.txt       # Python 依赖
├── herb_classes.txt      # 类别名称 (20 种中草药)
├── run.sh                # 启动脚本
├── .gitignore           # Git 忽略配置
└── README.md            # 本文档
```

## 快速开始

### 1. 创建 Conda 环境

```bash
conda create -n herb-inference python=3.10 -y
conda activate herb-inference
```

### 2. 安装依赖

```bash
pip install -r requirements.txt
```

### 3. 准备模型文件

将训练好的模型文件 `best_herb_model.pth` 放在项目根目录的 `models/` 文件夹：

```
HerbRecognition-be/
├── models/
│   └── best_herb_model.pth  # 模型文件
└── services/
    └── inference-service/
```

### 4. 启动服务

```bash
./run.sh
```

服务默认运行在 `http://localhost:5001`

## API 接口

### 健康检查

```bash
curl http://localhost:5001/health
```

响应：
```json
{
  "status": "ok",
  "model_loaded": true,
  "device": "cpu",
  "classes": 20
}
```

### 图像识别

#### 方式 1：文件上传

```bash
curl -X POST http://localhost:5001/predict \
  -F "image=@path/to/image.jpg"
```

#### 方式 2：Base64

```bash
curl -X POST http://localhost:5001/predict \
  -H "Content-Type: application/json" \
  -d '{"image": "data:image/jpeg;base64,/9j/4AAQSkZJRg..."}'
```

响应：
```json
{
  "success": true,
  "message": "识别成功",
  "data": {
    "herb_name": "Chicken Gizzard Membrane",
    "herb_id": 6,
    "confidence": 95.5,
    "all_results": [
      {"herb_name": "Chicken Gizzard Membrane", "herb_id": 6, "confidence": 95.5},
      {"herb_name": "Amomi Fructus", "herb_id": 1, "confidence": 2.3}
    ]
  }
}
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | 服务端口 | 5001 |

## 与 Go 后端集成

Go 后端会通过 HTTP 调用此服务进行识别。

设置环境变量：
```bash
export PYTHON_SERVICE_URL=http://localhost:5001
```

## 识别类别

支持 20 种中草药识别：

1. Amomi Fructus (砂仁)
2. Amomi Fructus Rotundus (豆蔻)
3. Apricot Kernels (杏仁)
4. Aster Root (紫菀)
5. Chaenomelis Fructus (木瓜)
6. Chicken Gizzard Membrane (鸡内金)
7. Cinnamon (肉桂)
8. Corni Fructus (山茱萸)
9. Foeniculi Fructus (小茴香)
10. Gardeniae Fructus (栀子)
11. Kochiae Fructus (地肤子)
12. Lycii Fructus (枸杞)
13. Mume Fructus (乌梅)
14. Psoraleae Fructus (补骨脂)
15. Rosae Laevigatae Frucyus (金樱子)
16. Rubi Fructus (覆盆子)
17. Schisandrae Chinensis Fructus (五味子)
18. Toosendan Fructus (川楝子)
19. Trichosanthis Pericarpopium (瓜蒌皮)
20. Turtle Shell (龟板)

## 开发说明

### 修改端口

编辑 `app.py` 末尾：
```python
app.run(host='0.0.0.0', port=5001, debug=False)
```

### 修改置信度阈值

编辑 `app.py`：
```python
CONFIDENCE_THRESHOLD = 0.3  # 低于 30% 置信度返回"未知"
```

### 生产部署

使用 Gunicorn：
```bash
gunicorn -w 4 -b 0.0.0.0:5001 app:app
```

## 注意事项

- 模型文件不会被提交到 Git（见 `.gitignore`）
- 训练好的模型请放在 `../../models/best_herb_model.pth`
- 服务默认使用 CPU，支持 CUDA 的设备会自动使用 GPU
