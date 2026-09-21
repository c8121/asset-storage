# InsightFace



## Install requirements for InsightFace FastAPI-Script

InsightFace uses python, so at first you need to install python.

    sudo apt install python3
    sudo apt install python3-pip
    sudo apt install python3-venv

### Linux

    pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

### Windows

    py -m pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

## Run the API

### Linux

Create virtual environment for python

    cd ~
    python3 -m venv venv

Run the service

    cd services/insightface
    source ~/venv/bin/activate
    uvicorn service:app --host 127.0.0.1 --port 8000

### Windows

    cd services/insightface
    py -m uvicorn service:app --host 127.0.0.1 --port 8000