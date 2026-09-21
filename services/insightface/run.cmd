python -m pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart


python -m uvicorn service:app --host 127.0.0.1 --port 8000