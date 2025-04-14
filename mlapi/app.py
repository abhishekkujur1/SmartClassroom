from flask import Flask, request, jsonify
from PIL import Image
import io
import numpy as np
import joblib  # or use pickle
import os

app = Flask(__name__)

# Load the model (make sure model.pkl is in the same folder or adjust the path)
MODEL_PATH = "./mlapi/model/face_recognition_model.pkl"
model = joblib.load(MODEL_PATH)

def preprocess_image(image_bytes):
    """Convert image bytes to numpy array for prediction."""
    image = Image.open(io.BytesIO(image_bytes)).convert("L")  # convert to grayscale if needed
    image = image.resize((28, 28))  # example resize
    image_array = np.array(image).reshape(1, -1) / 255.0  # flatten and normalize
    return image_array

@app.route('/predict', methods=['POST'])
def predict():
    if 'image' not in request.files:
        return jsonify({"error": "No image part in the request"}), 400

    file = request.files['image']
    if file.filename == '':
        return jsonify({"error": "No selected image"}), 400

    try:
        img_bytes = file.read()
        input_data = preprocess_image(img_bytes)
        prediction = model.predict(input_data)
        return jsonify({"prediction": prediction.tolist()})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True)
