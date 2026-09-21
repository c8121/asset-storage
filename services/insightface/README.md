# InsightFace

Use `services/insightface/run.sh` (Linux) or `services/insightface/run.cmd` (Windows) to run the InsightFace REST-Service.

These scripts will install the required modules and then run the http-service. For details see below.


## Requirements for InsightFace FastAPI-Script

InsightFace uses python, so at first you need to install python.


### Linux

    sudo apt install python3
    sudo apt install python3-pip
    sudo apt install python3-venv

Create virtual environment for python:

    cd ~
    python3 -m venv venv

Install modules:

    pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

### Windows

First install python from official installer.

Then install modules:

    python -m pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

## Running the API

The python-script `service.py` uses the python module `uvicorn` to run a http-server.

### Linux

    cd services/insightface
    source ~/venv/bin/activate
    uvicorn service:app --host 127.0.0.1 --port 8000

### Windows

    cd services/insightface
    python -m uvicorn service:app --host 127.0.0.1 --port 8000