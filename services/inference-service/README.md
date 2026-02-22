# 中草药识别推理服务

基于 MobileNetV3 的中草药图像识别服务。

## 目录结构

```
services/inference-service/
├── app.py                 # Flask 服务主程序
├── requirements.txt       # Python 依赖
├── herb_classes.txt      # 类别名称
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

将训练好的模型文件放在项目根目录的 `models/` 文件夹：

```
HerbRecognition-be/
├── models/
│   └── model.pth  # 模型文件
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
      {
        "herb_name": "Chicken Gizzard Membrane",
        "herb_id": 6,
        "confidence": 95.5
      },
      { "herb_name": "Amomi Fructus", "herb_id": 1, "confidence": 2.3 }
    ]
  }
}
```

## 环境变量

| 变量名 | 说明     | 默认值 |
| ------ | -------- | ------ |
| PORT   | 服务端口 | 5001   |
