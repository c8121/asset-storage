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


def encode_image_to_base64(image: np.ndarray) -> str:
    _, buffer = cv2.imencode(".jpg", image)
    return base64.b64encode(buffer).decode("utf-8")


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
            return JSONResponse(content={"faces": []})

        results = []

        for idx, face in enumerate(faces):
            bbox = face.bbox.astype(int)
            x1, y1, x2, y2 = bbox

            cropped_face = image[y1:y2, x1:x2]

            embedding = face.embedding.tolist()

            results.append({
                "face_index": idx,
                "embedding": embedding,
                "face_image_base64": encode_image_to_base64(cropped_face)
            })

        return {"faces": results}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

