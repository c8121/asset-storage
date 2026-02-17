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
            return JSONResponse(content={"Faces": []})

        results = []

        for idx, face in enumerate(faces):
            
            #cropped_face = crop_square_face(image, face.bbox, resize=224)
            cropped_face = crop_square_face(image, face.bbox)

            embedding = face.embedding.tolist()

            results.append({
                "Index": idx,
                "Embedding": embedding,
                "Image": encode_image_to_base64(cropped_face)
            })

        return {"Faces": results}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    

def crop_square_face(image: np.ndarray, bbox: np.ndarray, resize: int | None = None) -> np.ndarray:
    """
    Crop a face region into a square format.

    :param image: Original image (H, W, C)
    :param bbox: Bounding box array [x1, y1, x2, y2]
    :param resize: Optional size to resize output (e.g., 224)
    :return: Square cropped face image
    """
    x1, y1, x2, y2 = bbox.astype(int)

    # Original width and height
    w = x2 - x1
    h = y2 - y1

    # Side length of square
    side = max(w, h)

    # Center of bounding box
    cx = x1 + w // 2
    cy = y1 + h // 2

    # New square coordinates
    new_x1 = cx - side // 2
    new_y1 = cy - side // 2
    new_x2 = new_x1 + side
    new_y2 = new_y1 + side

    # Clip to image boundaries
    new_x1 = max(0, new_x1)
    new_y1 = max(0, new_y1)
    new_x2 = min(image.shape[1], new_x2)
    new_y2 = min(image.shape[0], new_y2)

    cropped = image[new_y1:new_y2, new_x1:new_x2]

    # Optional resize
    if resize is not None:
        cropped = cv2.resize(cropped, (resize, resize))

    return cropped


