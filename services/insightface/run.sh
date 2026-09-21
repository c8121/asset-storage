#!/bin/bash

WORK_DIR=$(realpath "$(dirname "$0")")

PYTHON_VENV_ACTIVATE=~/venv/bin/activate

if [ ! -f "$PYTHON_VENV_ACTIVATE" ] ; then
	echo "Create isolated python virtual environment"
	cd ~
	python3 -m venv venv

	source $PYTHON_VENV_ACTIVATE

	echo "Install modules"
	pip install fastapi uvicorn insightface onnxruntime opencv-python numpy python-multipart

else

	source $PYTHON_VENV_ACTIVATE

fi

cd "$WORK_DIR"
uvicorn service:app --host 127.0.0.1 --port 8000
