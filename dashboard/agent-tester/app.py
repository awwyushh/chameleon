import streamlit as st
import requests
import json
import pandas as pd
from datetime import datetime
import os

# Configuration
# If running inside Docker, use the container name. If local, use localhost.
# We default to localhost for easy testing if you run 'streamlit run app.py' manually.

#AGENT_URL = "http://localhost:8001/analyze" 
DEFAULT_URL = "http://localhost:8001/analyze"
AGENT_URL = os.getenv("AGENT_URL", DEFAULT_URL) 

st.set_page_config(page_title="Chameleon Brain Surgeon", layout="wide", page_icon="🧠")

st.title("🧠 Agentic AI: Tool Selection Tester")
st.markdown("Submit mock attack events to see how the **LlamaIndex ReAct Agent** reasons and selects tools.")
st.caption(f"Connecting to Agent at: `{AGENT_URL}`")

# --- Sidebar: Simulation Controls ---
st.sidebar.header("Simulate Attack")
attack_type = st.sidebar.selectbox(
    "Attack Type", 
    ["sqli", "xss", "bruteforce", "path_traversal", "scanning"]
)

cve_id = st.sidebar.text_input("Specific CVE (Optional)", "CVE-2023-NONE")
ip_addr = st.sidebar.text_input("Source IP", "192.168.1.55")
target_path = st.sidebar.text_input("Target Path", "/admin/login")

# Dynamic payload based on type
default_payloads = {
    "sqli": "' OR 1=1 --",
    "xss": "<script>alert('pwned')</script>",
    "bruteforce": "admin / password123",
    "scanning": "GET /wp-admin.php",
    "path_traversal": "../../../etc/passwd"
}
payload = st.sidebar.text_area("Raw Payload", default_payloads.get(attack_type, ""))

# --- Main Execution ---
if st.sidebar.button("🚨 Send Alert to Agent", type="primary"):
    
    # Construct the event payload matching main.py:AttackEvent
    event_data = {
        "src_ip": ip_addr,
        "label": attack_type,
        "payload": payload,
        "path": target_path
    }

    col1, col2 = st.columns([1, 1])

    with col1:
        st.subheader("📤 Outgoing Event")
        st.json(event_data)

    # Call the Agentic Service
    try:
        with st.spinner("Agent is thinking... (Connecting to Ollama )"):
            # Note: In docker-compose, update the URL to http://agentic-ai:8000/analyze
            response = requests.post(AGENT_URL, json=event_data, timeout=120)
            
        if response.status_code == 200:
            result = response.json()
            
            with col2:
                st.subheader("📥 Agent Decision")
                st.success(f"Status: {result.get('status')}")
            
            # --- Parsing and Visualization of Reasoning ---
            reasoning = result.get("ai_reasoning", "")
            
            st.divider()
            st.subheader("🧠 Chain of Thought (ReAct Trace)")
            
            # Parse lines for display and to find the tool used
            lines = reasoning.split('\n')
            
            detected_tool = "None" # Default
            
            for line in lines:
                if "Thought:" in line:
                    st.info(line)
                elif "Action:" in line:
                    st.warning(f"🛠️ **TOOL SELECTED:** {line}") # Highlight Tool Use
                    # Extract tool name for the history table
                    parts = line.split("Action:")
                    if len(parts) > 1:
                        detected_tool = parts[1].strip()
                elif "Action Input:" in line:
                    st.code(line, language="json")
                elif "Observation:" in line:
                    st.success(f"👁️ Tool Output: {line}")
                else:
                    st.write(line)

            # Update History Data (Simulated DB update)
            # In a real app, you'd append this to a list or DB
            new_row = {
                "Timestamp": datetime.now().strftime("%H:%M:%S"),
                "IP": ip_addr,
                "Attack": attack_type,
                "Tool Used": detected_tool
            }
            
            # Simple way to show history updates in Streamlit without a DB for the demo
            if 'history' not in st.session_state:
                st.session_state['history'] = []
            st.session_state['history'].insert(0, new_row)
                    
        else:
            st.error(f"Agent Error {response.status_code}: {response.text}")
                

    except requests.exceptions.ConnectionError:
        st.error(f"❌ Could not connect to Agentic Service at `{AGENT_URL}`. Is it running?")

# --- History (Mock) ---
st.divider()
st.subheader("📜 Recent Decisions")

if 'history' not in st.session_state:
    # Default Initial State
    st.session_state['history'] = [{
        "Timestamp": datetime.now().strftime("%H:%M:%S"),
        "IP": "127.0.0.1",
        "Attack": "Simulation",
        "Tool Used": "Waiting..."
    }]

st.table(pd.DataFrame(st.session_state['history']))