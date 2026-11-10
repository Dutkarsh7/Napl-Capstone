import io
import os
import logging
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from PIL import Image

try:
    import pytesseract
except ImportError:
    pytesseract = None

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ocr-service")

app = FastAPI(
    title="Smart Pharmacy OCR Service",
    description="Microservice for extracting prescription text from uploaded images using pytesseract.",
    version="1.0.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

MAX_FILE_SIZE = 5 * 1024 * 1024  # 5 MB
ALLOWED_CONTENT_TYPES = {"image/jpeg", "image/png", "image/jpg", "application/pdf"}

@app.get("/health")
def health_check():
    tesseract_available = False
    if pytesseract:
        try:
            pytesseract.get_tesseract_version()
            tesseract_available = True
        except Exception:
            tesseract_available = False

    return {
        "status": "healthy",
        "service": "ocr-python",
        "engine": "pytesseract",
        "tesseract_installed": tesseract_available,
        "disclaimer": "OCR output is an automated transcription suggestion only. All prescriptions require licensed pharmacist review."
    }

@app.post("/ocr")
async def extract_text(file: UploadFile = File(...)):
    if not file.filename:
        raise HTTPException(status_code=400, detail="Missing file in upload request")

    # Read bytes and enforce size limit
    contents = await file.read()
    if len(contents) > MAX_FILE_SIZE:
        raise HTTPException(status_code=400, detail="File size exceeds the 5 MB limit")

    if len(contents) == 0:
        raise HTTPException(status_code=400, detail="Uploaded file is empty")

    logger.info(f"Processing prescription image OCR: {file.filename} ({len(contents)} bytes)")

    # Attempt to open and verify image with Pillow
    try:
        image = Image.open(io.BytesIO(contents))
        image.verify()  # Validate image integrity
        image = Image.open(io.BytesIO(contents))  # Reopen for processing after verify
    except Exception as img_err:
        logger.warning(f"File could not be opened as standard image: {img_err}")
        return {
            "text": "(Uploaded document is in PDF format or unreadable as standard image. Pharmacist manual review required)",
            "warning": "Non-image or PDF format provided. Manual clinical verification required."
        }

    # Execute Tesseract OCR
    if not pytesseract:
        logger.error("pytesseract package not installed")
        return {
            "text": "(OCR service unavailable: pytesseract library missing. Mandatory pharmacist manual review required)",
            "error": "OCR engine unavailable"
        }

    try:
        raw_text = pytesseract.image_to_string(image)
    except pytesseract.TesseractNotFoundError:
        logger.warning("System Tesseract binary not installed on host machine")
        return {
            "text": "(System OCR binary not installed on host. Mandatory manual transcription required by pharmacist)",
            "error": "Tesseract binary not installed"
        }
    except Exception as ocr_err:
        logger.warning(f"Tesseract OCR extraction failed: {ocr_err}")
        return {
            "text": "(OCR processing failed on this document image. Mandatory pharmacist manual review required)",
            "error": str(ocr_err)
        }

    cleaned = raw_text.strip()
    if not cleaned:
        return {
            "text": "(No readable text detected in prescription image. Pharmacist manual review required)"
        }

    # Mark machine-extracted text clearly to ensure clinical safety
    labeled_output = f"[MACHINE EXTRACTED OCR SUGGESTION - REQUIRES PHARMACIST VERIFICATION]\n\n{cleaned}"

    return {
        "text": labeled_output
    }

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host="0.0.0.0", port=port)
