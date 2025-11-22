FROM python:3.11-slim-bookworm

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

WORKDIR /app

# Install system deps for building python extensions if needed
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Create Virtual Environment
ENV VIRTUAL_ENV=/opt/venv
RUN uv venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Copy dependency definition
COPY pyproject.toml .
# Install dependencies
RUN uv pip install -r pyproject.toml

# Copy code
COPY . .

# Expose port
EXPOSE 8000

# Run
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
