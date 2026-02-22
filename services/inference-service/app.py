import os
import io
import base64
import logging
from datetime import datetime

import numpy as np
import torch
import torch.nn as nn
from PIL import Image
from torchvision import models, transforms
from flask import Flask, request, jsonify
from flask_cors import CORS

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

DEVICE = torch.device("cuda" if torch.cuda.is_available() else "cpu")
MODEL_PATH = os.path.join(os.path.dirname(__file__), "../../models/best_herb_model.pth")
CLASS_NAMES_PATH = os.path.join(os.path.dirname(__file__), "herb_classes.txt")

IMG_SIZE = 224
CONFIDENCE_THRESHOLD = 0.3

app = Flask(__name__)
CORS(app)

class_names = []
model = None

def load_class_names():
    global class_names
    with open(CLASS_NAMES_PATH, 'r', encoding='utf-8') as f:
        class_names = [line.strip() for line in f if line.strip()]
    logger.info(f"加载类别名称：{len(class_names)} 类")

def load_model():
    global model
    
    model = models.mobilenet_v3_large(weights=None)
    in_features = model.classifier[-1].in_features
    model.classifier[-1] = nn.Linear(in_features, len(class_names))
    
    if os.path.exists(MODEL_PATH):
        checkpoint = torch.load(MODEL_PATH, map_location=DEVICE, weights_only=False)
        model.load_state_dict(checkpoint['model_state_dict'])
        logger.info(f"加载模型成功：{MODEL_PATH}")
    else:
        logger.warning(f"模型文件不存在：{MODEL_PATH}")
    
    model = model.to(DEVICE)
    model.eval()
    logger.info(f"模型运行在：{DEVICE}")

val_transform = transforms.Compose([
    transforms.Resize(int(IMG_SIZE * 1.14)),
    transforms.CenterCrop(IMG_SIZE),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                         std=[0.229, 0.224, 0.225]),
])

def preprocess_image(image_data):
    if isinstance(image_data, str) and (image_data.startswith('data:image') or len(image_data) > 1000):
        if image_data.startswith('data:image'):
            image_data = image_data.split(',')[1]
        image_bytes = base64.b64decode(image_data)
        image = Image.open(io.BytesIO(image_bytes)).convert('RGB')
    elif isinstance(image_data, str) and os.path.exists(image_data):
        image = Image.open(image_data).convert('RGB')
    else:
        image = Image.fromarray(image_data).convert('RGB')
    
    image_tensor = val_transform(image).unsqueeze(0).to(DEVICE)
    return image_tensor

@torch.no_grad()
def predict(image_tensor):
    model.eval()
    outputs = model(image_tensor)
    probabilities = torch.softmax(outputs, dim=1)
    confidence, predicted = torch.max(probabilities, 1)
    
    confidence_score = float(confidence.item())
    predicted_class = int(predicted.item())
    
    top5_confidence, top5_indices = torch.topk(probabilities[0], k=min(5, len(class_names)))
    
    results = []
    for conf, idx in zip(top5_confidence.cpu().numpy(), top5_indices.cpu().numpy()):
        if conf > CONFIDENCE_THRESHOLD:
            results.append({
                "herb_name": class_names[idx],
                "herb_id": int(idx + 1),
                "confidence": round(float(conf) * 100, 2)
            })
    
    return confidence_score, predicted_class, results

@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        "status": "ok",
        "model_loaded": model is not None,
        "device": str(DEVICE),
        "classes": len(class_names),
        "timestamp": datetime.now().isoformat()
    })

@app.route('/predict', methods=['POST'])
def predict_endpoint():
    try:
        if 'image' not in request.files and 'image' not in request.json:
            return jsonify({"error": "缺少图像数据", "code": "MISSING_IMAGE"}), 400
        
        if 'image' in request.files:
            file = request.files['image']
            image_bytes = file.read()
            image = Image.open(io.BytesIO(image_bytes)).convert('RGB')
            image_tensor = val_transform(image).unsqueeze(0).to(DEVICE)
        else:
            data = request.json
            image_data = data.get('image')
            if not image_data:
                return jsonify({"error": "图像数据为空", "code": "EMPTY_IMAGE"}), 400
            image_tensor = preprocess_image(image_data)
        
        confidence, predicted_class, top_results = predict(image_tensor)
        
        if not top_results:
            return jsonify({
                "success": True,
                "message": "置信度过低，无法识别",
                "data": {
                    "herb_name": "未知",
                    "herb_id": 0,
                    "confidence": round(confidence * 100, 2),
                    "all_results": []
                }
            }), 200
        
        top_result = top_results[0]
        
        logger.info(f"识别结果：{top_result['herb_name']} (置信度：{top_result['confidence']}%)")
        
        return jsonify({
            "success": True,
            "message": "识别成功",
            "data": {
                "herb_name": top_result["herb_name"],
                "herb_id": top_result["herb_id"],
                "confidence": top_result["confidence"],
                "all_results": top_results
            }
        }), 200
        
    except Exception as e:
        logger.error(f"预测失败：{str(e)}", exc_info=True)
        return jsonify({
            "success": False,
            "error": str(e),
            "code": "PREDICTION_ERROR"
        }), 500

@app.route('/predict/batch', methods=['POST'])
def batch_predict():
    try:
        data = request.json
        images = data.get('images', [])
        
        if not images:
            return jsonify({"error": "图像列表为空", "code": "EMPTY_IMAGES"}), 400
        
        results = []
        for i, image_data in enumerate(images):
            try:
                image_tensor = preprocess_image(image_data)
                confidence, predicted_class, top_results = predict(image_tensor)
                results.append({
                    "index": i,
                    "success": True,
                    "data": top_results[0] if top_results else None
                })
            except Exception as e:
                results.append({
                    "index": i,
                    "success": False,
                    "error": str(e)
                })
        
        return jsonify({
            "success": True,
            "count": len(results),
            "results": results
        }), 200
        
    except Exception as e:
        logger.error(f"批量预测失败：{str(e)}", exc_info=True)
        return jsonify({
            "success": False,
            "error": str(e),
            "code": "BATCH_PREDICTION_ERROR"
        }), 500

def initialize():
    logger.info("正在初始化识别服务...")
    load_class_names()
    load_model()
    logger.info("初始化完成")

if __name__ == '__main__':
    initialize()
    app.run(host='0.0.0.0', port=5001, debug=False)
