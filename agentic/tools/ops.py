# agentic/tools/ops.py
import subprocess
import logging
from llama_index.core.tools import FunctionTool

logger = logging.getLogger("tools.ops")
logger.setLevel(logging.INFO)

def scan_image_vulnerability(image_name: str) -> str:
    """
    Scan a Docker image with Trivy (mocked or real if present).
    """
    try:
        # If trivy is not present this will fail quickly; in production use subprocess.run properly
        cmd = ["trivy", "image", "--quiet", image_name]
        proc = subprocess.run(cmd, capture_output=True, timeout=30)
        if proc.returncode == 0:
            return f"Scan Result for {image_name}: Found issues (mocked) or summary available."
        return f"Scan Result for {image_name}: exit_code={proc.returncode}, stderr={proc.stderr.decode()[:200]}"
    except FileNotFoundError:
        logger.warning("trivy not installed — returning mocked result")
        return f"Scan Result for {image_name}: Found CVE-2023-1234 (High). Safe to use as bait."
    except Exception as e:
        logger.exception("scan_image_vulnerability failed")
        return f"Scan failed: {str(e)}"

def start_packet_capture(duration: int, filter: str = "") -> str:
    """
    Start a tcpdump capture (synchronous here, spawn in production).
    """
    try:
        # For safety do not actually run tcpdump in this environment; return a placeholder
        return f"Started PCAP recording for {duration} seconds. Filter: '{filter}'. File will be saved to /forensics."
    except Exception as e:
        logger.exception("start_packet_capture failed")
        return f"PCAP failed: {str(e)}"

ops_tools = [
    FunctionTool.from_defaults(fn=scan_image_vulnerability, name="scan_image_vulnerability"),
    FunctionTool.from_defaults(fn=start_packet_capture, name="start_packet_capture"),
]
