🏗 System Architecture

The system operates as a distributed microservices cluster connected via a private Docker network (chameleon-net).
1. Agentd (The Reflex System)

    Tech: Go (Golang), gRPC (Unix Socket), Zerolog.

    Role: High-performance decision gateway. It sits inline with traffic.

    Logic:

        Fast Path: Checks Regex rules (/rules/*.yaml) and cached ML scores.

        Tarpit: Uses Redis to track IP windows and calculates variable delay (logarithmic growth) to slow down brute-force attacks.

        Decision: Returns PASS, TARPIT, DECEIVE (injects fake SQL errors), or SPAWN (spins up a honeypot).

        Telemetry: Asynchronously pushes event data to the Aggregator via JWT-signed requests.

2. Agentic AI (The Brain)

    Tech: Python 3.11, FastAPI, LlamaIndex, Google Gemini Pro.

    Role: Long-term reasoning and active defense orchestration.

    Workflow:

        Receives high-confidence alerts.

        Uses a ReAct (Reasoning + Acting) loop to analyze intent.

        Tools:

            deploy_honeypot: Spawns Docker containers (e.g., mysql-admin) dynamically.

            escalate_tarpit: Increases Redis penalty keys for specific IPs.

            enable_labyrinth_mode: Flags an IP to receive infinite fake content.

            query_threat_intel: Queries Vector DB (Chroma) for historical attacker behavior.

3. ML Service (The Classifier)

    Tech: Python, Flask, ChromaDB (Vector Store), Ollama (Local LLM).

    Role:

        Predict: Provides probability scores for payloads (SQLi, XSS, Bruteforce).

        Analyze: Uses RAG to compare current attacks against vector embeddings of past incidents.

4. The Application (Victim)

    Tech: Node.js, Express.

    Integration: Uses a custom Chameleon SDK middleware.

    Function: Captures request metadata, sends it to agentd via gRPC, and executes the returned decision (e.g., sleep for 500ms, render fake error).

5. Visualization Layer

    Frontend: Next.js (React 19), Tailwind CSS, Shadcn UI, Three.js (Particle effects).

    Aggregator: Node.js service that ingests telemetry from agentd and broadcasts via WebSockets.

    Agent Tester: Streamlit app for debugging the Agentic AI's "Chain of Thought" without generating real network traffic.
