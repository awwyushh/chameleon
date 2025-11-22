import os
import logging
from flask import Flask, request, jsonify
from llama_index.core import VectorStoreIndex, Document, Settings
from llama_index.vector_stores.chroma import ChromaVectorStore
from llama_index.llms.ollama import Ollama
from llama_index.embeddings.huggingface import HuggingFaceEmbedding
import chromadb

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)

# --- CONFIGURATION ---
CHROMA_HOST = os.getenv("CHROMA_HOST", "localhost")
CHROMA_PORT = os.getenv("CHROMA_PORT", "8000")
OLLAMA_URL = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")

# 1. Setup LLM (Ollama) & Embeddings
# using a small, fast model for the hackathon
Settings.llm = Ollama(model="llama3", base_url=OLLAMA_URL, request_timeout=120.0)
Settings.embed_model = HuggingFaceEmbedding(model_name="BAAI/bge-small-en-v1.5")

# 2. Setup Vector DB (Chroma)
remote_db = chromadb.HttpClient(host=CHROMA_HOST, port=int(CHROMA_PORT))
chroma_collection = remote_db.get_or_create_collection("chameleon_attacks")
vector_store = ChromaVectorStore(chroma_collection=chroma_collection)
index = VectorStoreIndex.from_vector_store(vector_store=vector_store)

@app.route("/predict", methods=["POST"])
def predict():
    """Fast path: Quick heuristic/ML classification for Agentd blocking"""
    data = request.json
    text = data.get("text", "").lower()
    
    # Keep the simple rules for millisecond-latency blocking
    if "select" in text or "union" in text or "--" in text:
        return jsonify({"label": "sqli", "confidence": 0.95})
    if "<script>" in text or "alert(" in text:
        return jsonify({"label": "xss", "confidence": 0.92})
    if "password" in text or "login" in text:
        return jsonify({"label": "bruteforce", "confidence": 0.60})
    
    return jsonify({"label": "benign", "confidence": 0.20})

@app.route("/analyze", methods=["POST"])
def analyze():
    """Slow path: Agentic RAG Analysis for Dashboard"""
    data = request.json
    payload = data.get("body", "")
    src_ip = data.get("src_ip", "unknown")
    attack_type = data.get("label", "unknown")

    # 1. Store this attack in Vector DB for future context
    doc_text = f"Attack Type: {attack_type} | IP: {src_ip} | Payload: {payload}"
    doc = Document(text=doc_text, metadata={"ip": src_ip, "type": attack_type})
    index.insert(doc)

    # 2. Perform RAG: Ask LLM to analyze based on history
    query_engine = index.as_query_engine()
    
    prompt = (
        f"Analyze this new attack: '{payload}' from IP {src_ip}. "
        f"Check the context for previous attacks from this IP or similar payloads. "
        "Is this part of a coordinated campaign? Suggest 1 specific defensive action."
    )
    
    try:
        response = query_engine.query(prompt)
        return jsonify({"analysis": str(response), "status": "success"})
    except Exception as e:
        logging.error(f"LLM Error: {e}")
        return jsonify({"analysis": "AI Analysis unavailable (Ollama might be offline).", "status": "error"})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001)
