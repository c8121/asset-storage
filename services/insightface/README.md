# InsightFace



## Install requirements for InsightFace FastAPI-Script

InsightFace uses python, so at first you need to install python.

### Linux

    pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

### Windows

    py -m pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

## Run the API

### Linux

    cd services/insightface
    uvicorn service:app --host 127.0.0.1 --port 8000

### Windows

    cd services/insightface
    py -m uvicorn service:app --host 127.0.0.1 --port 8000