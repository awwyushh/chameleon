const express = require('express');
const bodyParser = require('body-parser');
const chameleon = require('chameleon-sdk'); // Local link

const app = express();
const PORT = 3000;

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// --- CHAMELEON MIDDLEWARE ---
// Protects all routes below this line
app.use(chameleon());

// --- ROUTES ---

app.get('/', (req, res) => {
    res.send(`
        <h1>Vulnerable App</h1>
        <ul>
            <li><a href="/search?q=test">Search (Vulnerable to XSS/SQLi)</a></li>
            <li><a href="/login">Login (Vulnerable to Bruteforce)</a></li>
        </ul>
    `);
});

// Simulated Search Endpoint
app.get('/search', (req, res) => {
    const query = req.query.q || '';
    // Real app would query DB here. 
    // Since middleware didn't block/deceive, we show "normal" results.
    res.send(`
        <h1>Search Results</h1>
        <p>Found 0 results for: ${query}</p>
        <a href="/">Back</a>
    `);
});

// Simulated Login Endpoint
app.post('/login', (req, res) => {
    const { username, password } = req.body;
    if (username === 'admin' && password === 'supersecret') {
        res.json({ status: "success", token: "12345" });
    } else {
        res.status(401).json({ status: "failed" });
    }
});

app.listen(PORT, () => {
    console.log(`Demo App running on http://localhost:${PORT}`);
});
