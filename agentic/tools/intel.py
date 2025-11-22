# agentic/tools/intel.py
import os
import logging
import chromadb
from llama_index.core.tools import FunctionTool

logger = logging.getLogger("tools.intel")
logger.setLevel(logging.INFO)

CHROMA_HOST = os.getenv("CHROMA_HOST", "chroma")
CHROMA_PORT = int(os.getenv("CHROMA_PORT", "8000"))

try:
    chroma_client = chromadb.HttpClient(host=CHROMA_HOST, port=CHROMA_PORT)
    collection = chroma_client.get_or_create_collection("chameleon_attacks")
except Exception:
    chroma_client = None
    collection = None
    logger.exception("Chroma init failed")

def query_threat_intel(ip: str) -> str:
    """
    Query Chroma for past events for this IP.
    """
    try:
        if collection is None:
            return "No prior history found (chroma unavailable)."
        results = collection.query(query_texts=[ip], n_results=5, where={"ip": ip})
        docs = results.get("documents", [[]])[0]
        if not docs:
            return "No prior history found for this IP."
        history = "\n".join(docs)
        logger.info(f"Queried threat intel for {ip}: {len(docs)} docs")
        return f"History for {ip}:\n{history}"
    except Exception as e:
        logger.exception("query_threat_intel failed")
        return f"Failed to query threat intel: {str(e)}"

intel_tools = [
    FunctionTool.from_defaults(fn=query_threat_intel, name="query_threat_intel")
]
