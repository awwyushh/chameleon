#!/bin/bash

# 1. Configuration
IMAGE="chp-admin:latest"
API_URL="http://localhost:4000/v1/trivy"

echo "[SaaS-CI] 🚀 Starting security certification for $IMAGE..."

# 2. Run Trivy (Containerized)
# We output JSON so our dashboard can parse it
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
    aquasec/trivy:latest image \
    --format json \
    --scanners vuln \
    $IMAGE > report.json

echo "[SaaS-CI] ✅ Scan complete. Uploading report to Dashboard..."

# 3. Upload to Aggregator (Simulating cloud webhook)
curl -X POST -H "Content-Type: application/json" \
     -d @report.json \
     $API_URL

echo "\n[SaaS-CI] 🏁 Pipeline finished."
