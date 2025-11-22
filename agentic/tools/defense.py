# agentic/tools/defense.py
import docker
import redis
import os
import logging
from llama_index.core.tools import FunctionTool

logger = logging.getLogger("tools.defense")
logger.setLevel(logging.INFO)

REDIS_HOST = os.getenv("REDIS_HOST", "redis")
redis_client = redis.Redis(host=REDIS_HOST, port=6379, db=0, decode_responses=True)

try:
    docker_client = docker.from_env()
except Exception:
    docker_client = None  # Docker not available in this environment

def deploy_honeypot(type: str, sticky: bool = True) -> str:
    """
    Spins up a honeypot container (mysql, admin, wordpress) and returns deployed port/ID.
    """
    if not docker_client:
        return "Error: Docker socket not connected."

    image_map = {
        "mysql": os.getenv("HP_IMAGE_MYSQL", "chm-mysql-admin:latest"),
        "admin": os.getenv("HP_IMAGE_ADMIN", "chp-admin:latest"),
        "wordpress": os.getenv("HP_IMAGE_WORDPRESS", "chp-wordpress:latest")
    }
    image = image_map.get(type, image_map["admin"])

    try:
        container = docker_client.containers.run(
            image,
            detach=True,
            ports={"80/tcp": None},
            labels={"chameleon.type": "honeypot"}
        )
        container.reload()
        host_port = container.ports["80/tcp"][0]["HostPort"]
        logger.info(f"Deployed {type} honeypot at port {host_port} id={container.short_id}")
        return f"Success: Deployed {type} honeypot on port {host_port}. ID: {container.short_id}"
    except Exception as e:
        logger.exception("deploy_honeypot failed")
        return f"Failed to deploy honeypot: {str(e)}"

def escalate_tarpit(ip: str, level: int) -> str:
    """
    Set tarpit penalty level for an IP (1-5) in Redis. Agentd should read this key.
    """
    try:
        key = f"tarpit:penalty:{ip}"
        level = max(1, min(5, int(level)))
        redis_client.setex(key, 300, level)  # 5 minutes
        logger.info(f"Set tarpit for {ip} -> {level}")
        return f"Success: IP {ip} is now tarpitted at level {level} for 5 minutes."
    except Exception as e:
        logger.exception("escalate_tarpit failed")
        return f"Failed to set tarpit: {str(e)}"

def enable_labyrinth_mode(ip: str) -> str:
    """
    Enable labyrinth mode in Redis for the IP (returns fake content for scanners).
    """
    try:
        key = f"labyrinth:{ip}"
        redis_client.setex(key, 600, "true")  # 10 minutes
        logger.info(f"Enabled labyrinth mode for {ip}")
        return f"Success: Labyrinth Mode active for {ip}. Scanner will receive false positives."
    except Exception as e:
        logger.exception("enable_labyrinth_mode failed")
        return f"Failed to enable labyrinth mode: {str(e)}"

def inject_canary(type: str) -> str:
    """
    Return a template id to serve a canary token.
    types: 'aws_keys', 'fake_flag', 'tracking_pixel'
    """
    try:
        map_type = {
            "aws_keys": "fake_aws_creds",
            "fake_flag": "fake_etc_passwd",
            "tracking_pixel": "html_pixel_trap"
        }
        tid = map_type.get(type, "default_generic")
        logger.info(f"Generated canary template {tid} for type {type}")
        return f"Action Required: Serve template '{tid}' to the user."
    except Exception as e:
        logger.exception("inject_canary failed")
        return f"Failed to generate canary: {str(e)}"

# Expose as LlamaIndex FunctionTool list for engine to import
defense_tools = [
    FunctionTool.from_defaults(fn=deploy_honeypot, name="deploy_honeypot"),
    FunctionTool.from_defaults(fn=escalate_tarpit, name="escalate_tarpit"),
    FunctionTool.from_defaults(fn=enable_labyrinth_mode, name="enable_labyrinth_mode"),
    FunctionTool.from_defaults(fn=inject_canary, name="inject_canary"),
]
