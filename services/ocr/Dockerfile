FROM python:3.11-slim

WORKDIR /app

# Install system dependencies including Tesseract OCR engine and language data
RUN apt-get update && apt-get install -y --no-install-recommends \
    tesseract-ocr \
    tesseract-ocr-eng \
    libgl1 \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . ./

EXPOSE 8000
ENV PORT=8000

CMD ["python", "main.py"]
