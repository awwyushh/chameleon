const express = require("express");
const axios = require("axios")

const app = express();
app.use(express.json());


const AI_SERVICE_URL = process.env.AI_SERVICE_URL || 'http://localhost:5001';

app.post('/v1/events', verifyToken, async (req, res) => {
  const event = req.body;
  console.log(`[Event] ${event.src_ip} -> ${event.label}`);

  // 1. Broadcast raw event immediately to UI (Real-time)
  io.emit('new_event', { ...event, timestamp: new Date().toISOString() });

  // 2. Trigger Agentic AI Analysis (Async - fire and forget)
  if (event.label !== 'benign') {
    analyzeAttack(event);
  }

  res.status(200).send({ status: 'received' });
});

async function analyzeAttack(event) {
  try {
    // Send to Python RAG service
    const response = await axios.post(`${AI_SERVICE_URL}/analyze`, {
      src_ip: event.src_ip,
      label: event.label,
      body: event.path + " " + (event.payload_excerpt || "") // Use path/payload context
    });

    const aiAnalysis = response.data.analysis;
    
    // Broadcast the AI Insight to the Dashboard
    console.log(`[AI Agent] Insight generated for ${event.src_ip}`);
    io.emit('ai_insight', {
      src_ip: event.src_ip,
      analysis: aiAnalysis,
      timestamp: new Date().toISOString()
    });

  } catch (error) {
    console.error("[AI Agent] Analysis failed:", error.message);
  }
}


app.post("/v1/events", (req, res) => {
    console.log("Event received:", req.body);
    return res.json({status: "ok"});
});

// For safety — aggregator health check
app.get("/health", (req, res) => res.json({status: "ok"}));


app.post('/v1/trivy', (req, res) => {
  const report = req.body;
  
  // Extract just the high-level summary for the UI
  const results = report.Results || [];
  let vulnerabilities = [];

  results.forEach(target => {
    if (target.Vulnerabilities) {
      vulnerabilities = [ ...vulnerabilities, ...target.Vulnerabilities ];
    }
  });

  // Sort by severity (Critical/High first)
  vulnerabilities.sort((a, b) => {
    const priority = { 'CRITICAL': 3, 'HIGH': 2, 'MEDIUM': 1, 'LOW': 0 };
    return priority[b.Severity] - priority[a.Severity];
  });

  // Take top 5 for the dashboard widget
  const topVulns = vulnerabilities.slice(0, 5).map(v => ({
    id: v.VulnerabilityID,
    pkg: v.PkgName,
    severity: v.Severity,
    title: v.Title
  }));

  console.log(`[Trivy] Processed scan. Found ${vulnerabilities.length} issues.`);
  
  // Push to frontend
  io.emit('trivy_update', {
    image: 'chp-admin:latest',
    scan_time: new Date().toISOString(),
    vulnerabilities: topVulns
  });

  res.json({ status: 'received' });
});

app.listen(4000, () => console.log("Aggregator running on port 4000"));

