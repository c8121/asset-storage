import io
import cv2
import base64
import numpy as np
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import JSONResponse
from insightface.app import FaceAnalysis

app = FastAPI(title="InsightFace Face Extractor API")

# Initialize InsightFace model
face_app = FaceAnalysis(name="buffalo_l")
face_app.prepare(ctx_id=0, det_size=(640, 640))  # ctx_id=0 for GPU, -1 for CPU


@app.get("/status")
async def status():
    try:
        return "Ready"
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/extract")
async def extract_faces(file: UploadFile = File(...)):
    try:
        contents = await file.read()
        np_arr = np.frombuffer(contents, np.uint8)
        image = cv2.imdecode(np_arr, cv2.IMREAD_COLOR)

        if image is None:
            raise HTTPException(status_code=400, detail="Invalid image file")

        faces = face_app.get(image)

        if not faces:
            return JSONResponse(content={"Faces": []})

        results = []

        for idx, face in enumerate(faces):

            embedding = face.embedding.tolist()
            bbox_list = face.bbox.astype(int).tolist()

            results.append({
                "Index": idx,
                "Embedding": embedding,
                "Image": bbox_list
            })

        return {"Faces": results}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

