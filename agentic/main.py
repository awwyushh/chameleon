# agentic/main.py
import logging
from fastapi import FastAPI
from pydantic import BaseModel
from engine import run_agent, build_threat_context
from tools.intel import query_threat_intel

logger = logging.getLogger("agentic")
logger.setLevel(logging.INFO)

app = FastAPI()

class AttackEvent(BaseModel):
    src_ip: str
    label: str
    payload: str
    path: str

@app.post("/analyze")
async def analyze_event(event: AttackEvent):
    prompt = (
        f"Incoming Alert: {event.label} attack detected from IP {event.src_ip}. "
        f"Target Path: {event.path}. "
        f"Payload Preview: {event.payload[:400]}. "
        "Decide on a response strategy and follow the ReAct format."
    )
    logger.info(f"Agent processing alert from {event.src_ip} label={event.label}")

    # Optionally attach intel history
    history = query_threat_intel(event.src_ip)
    # Build and/or include structured context (not required by ReAct but useful)
    threat_ctx = build_threat_context(event.src_ip, event.label, event.payload, event.path, history=history)

    # Prepend structured metadata to prompt to reduce hallucinations
    structured = {
        "threat": {
            "ip": threat_ctx.ip,
            "category": threat_ctx.category,
            "confidence": threat_ctx.confidence,
            "path": threat_ctx.path,
            "history": threat_ctx.history_summary or ""
        },
        "incident": {
            "label": event.label,
            "payload_preview": event.payload[:400]
        }
    }
    combined_prompt = f"STRUCTURED_CONTEXT: {structured}\n\nUSER_PROMPT: {prompt}"

    response_text = await run_agent(combined_prompt)
    return {"status": "executed", "ai_reasoning": response_text}

@app.get("/health")
def health():
    return {"status": "agentic_brain_online"}
