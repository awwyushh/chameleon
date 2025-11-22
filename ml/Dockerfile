FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
# Install build essentials for some python libs if needed
RUN apt-get update && apt-get install -y build-essential && \
    pip install --no-cache-dir -r requirements.txt
COPY service.py .
EXPOSE 5001
CMD ["python", "service.py"]
