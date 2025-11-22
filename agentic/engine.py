# agentic/engine.py
import os
import json
import logging
from dataclasses import dataclass, asdict
from typing import Optional

from llama_index.core.agent.workflow import ReActAgent, AgentStream, ToolCallResult
from llama_index.core.workflow import Context
from llama_index.core import PromptTemplate

# Gemini / google-genai
import google.genai as genai
from llama_index.llms.gemini import Gemini

# Import tools (they are FunctionTool lists)
from tools.defense import defense_tools
from tools.intel import intel_tools
from tools.ops import ops_tools

logger = logging.getLogger("agentic.engine")
logger.setLevel(logging.INFO)

# -------------------------
# ThreatContext
# -------------------------
@dataclass
class ThreatContext:
    category: str
    confidence: float
    ip: str
    path: str
    payload: str
    sanitized_payload: str
    mode: str = "forensic"
    history_summary: Optional[str] = None

def sanitize_payload(payload: str) -> str:
    if not isinstance(payload, str):
        return ""
    p = payload[:2000]
    p = "".join(ch for ch in p if 32 <= ord(ch) <= 126)
    for m in ["{{", "{%", "<svg", "onerror", "onload"]:
        p = p.replace(m, "")
    return p

# -------------------------
# Gemini client and LLM wrapper
# -------------------------
# Hardcoded API key as requested — replace the string with your real key
GEMINI_API_KEY = "AIzaSyAE6F8WXuFQe6AzVvzq1_AWXTudnqK3uL8"
MODEL_NAME = os.getenv("GEMINI_MODEL", "models/gemini-2.5-pro")

# Initialize genai client (optional, for direct usage if needed)
client = genai.Client(api_key=GEMINI_API_KEY)

# LlamaIndex Gemini wrapper config
llm = Gemini(
    model=MODEL_NAME,
    api_key=GEMINI_API_KEY,
    temperature=0.15,
    max_output_tokens=4096
)
logger.info(f"Initialized Gemini LLM {MODEL_NAME}")

# -------------------------
# Aggregate tools
# -------------------------
# defense_tools, intel_tools, ops_tools are lists of FunctionTool objects
all_tools = []
all_tools.extend(defense_tools)
all_tools.extend(intel_tools)
all_tools.extend(ops_tools)

# -------------------------
# Create ReAct Agent
# -------------------------
agent = ReActAgent(tools=all_tools, llm=llm)

custom_header = """\
You are the Chameleon Cyber Defense AI.
Your job is to analyze incoming web requests and execute active defense strategies.

Please respond strictly with the ReAct format:

Thought: <short reasoning>
Action: <tool name or "none">
Action Input: <JSON kwargs>

If you call a tool, await the Observation then produce a Final Answer.

Guidelines:
- If an IP is scanning (Gobuster/Dirbuster), call 'enable_labyrinth_mode'.
- If attacker probes cloud data (AWS), call 'inject_canary' with 'aws_keys'.
- If attacker is persistent, call 'query_threat_intel' then consider 'deploy_honeypot'.
- Prefer gathering intel before destructive actions.
"""

react_system_prompt = PromptTemplate(custom_header)
agent.update_prompts({"react_header": react_system_prompt})

# -------------------------
# run_agent: execute and collect ReAct trace
# -------------------------
async def run_agent(prompt: str) -> str:
    ctx = Context(agent)
    handler = agent.run(prompt, ctx=ctx)

    trace_output = []
    try:
        async for ev in handler.stream_events():
            if isinstance(ev, AgentStream):
                # stream deltas (partial thinking)
                if ev.delta:
                    trace_output.append(f"Stream: {ev.delta}")

            if isinstance(ev, ToolCallResult):
                trace_output.append("Thought: I need to use a tool to mitigate this threat.")
                trace_output.append(f"Action: {ev.tool_name}")
                trace_output.append(f"Action Input: {json.dumps(ev.tool_kwargs)}")
                # ev.tool_output typically has .content
                content = getattr(ev.tool_output, "content", str(ev.tool_output))
                trace_output.append(f"Observation: {content}")
                trace_output.append("-" * 20)

        response = await handler

    except Exception as e:
        logger.exception("Agent run failed")
        trace_output.append(f"Final Answer: Agent execution failed: {str(e)}")
        return "\n".join(trace_output)

    if not trace_output:
        trace_output.append("Thought: No specific tools were required, providing analysis directly.")

    trace_output.append(f"Final Answer: {str(response)}")
    return "\n".join(trace_output)

# -------------------------
# build_threat_context helper (optional usage)
# -------------------------
def build_threat_context(src_ip: str, label: str, payload: str, path: str, mode: str = "forensic", history: Optional[str] = None) -> ThreatContext:
    sanitized = sanitize_payload(payload)
    mapping = {
        "sqli": ("sql_injection", 0.95),
        "xss": ("xss", 0.95),
        "bruteforce": ("auth_bruteforce", 0.9),
        "scanning": ("scanner", 0.9),
        "path_traversal": ("path_traversal", 0.95)
    }
    category, conf = mapping.get(label, ("unknown", 0.5))
    return ThreatContext(category=category, confidence=conf, ip=src_ip, path=path, payload=payload, sanitized_payload=sanitized, mode=mode, history_summary=history)
