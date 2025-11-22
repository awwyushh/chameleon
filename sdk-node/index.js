const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// Point to the proto file in the repo
const PROTO_PATH = path.resolve(__dirname, '../agentd/proto/chameleon.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const proto = grpc.loadPackageDefinition(packageDefinition).chameleon;

// Default client connection to UDS
const client = new proto.Chameleon(
    'unix:///tmp/chameleon.sock',
    grpc.credentials.createInsecure()
);

module.exports = function chameleonMiddleware() {
    return (req, res, next) => {
        // Combine body and query parameters for full inspection
            const inspectionText = (req.method === 'GET') 
                ? JSON.stringify(req.query) 
                : JSON.stringify(req.body) + JSON.stringify(req.query);

            const payload = {
                src_ip: req.ip || '127.0.0.1',
                path: req.path,
                method: req.method,
                body: inspectionText || "", // Now includes ?q=<script>...
            };

        client.Classify(payload, (err, decision) => {
            if (err) {
                // Fail open if agent is down
                // console.error("Chameleon Agent Error:", err.message);
                return next();
            }

            // 1. TARPIT (Delay)
            const delay = decision.delay_ms || 0;
            if (delay > 0) {
                // Simply sleep before proceeding
                setTimeout(() => processDecision(decision, res, next), delay);
            } else {
                processDecision(decision, res, next);
            }
        });
    };
};

function processDecision(decision, res, next) {
    // Log for debug
    if (decision.label !== 'benign') {
        console.log(`[Chameleon] Detected ${decision.label} (Confidence: ${decision.confidence.toFixed(2)}) -> Action: ${decision.action}`);
    }

    // 2. DECEIVE
    if (decision.action === 'DECEIVE' || decision.action === 1) {
        // If Honeypot Spawn was triggered (we map SPAWN to DECEIVE in agentd with hp details)
        if (decision.honeypot_port > 0) {
             console.log(`[Chameleon] Redirecting to Honeypot at http://${decision.honeypot_host}:${decision.honeypot_port}`);
             return res.redirect(`http://${decision.honeypot_host}:${decision.honeypot_port}`);
        }

        // Standard Deception (render template)
        // We try to detect content type from the body text
        if (decision.message.trim().startsWith('{')) {
            res.setHeader('Content-Type', 'application/json');
        } else {
            res.setHeader('Content-Type', 'text/html');
        }
        res.status(200); // Return 200 to fool scanners often
        return res.send(decision.message);
    }

    // 3. PASS
    next();
}
