const https = require('https');
const crypto = require('crypto');
const { URL } = require('url');

const CLIENT_ID = process.argv[2];
const CLIENT_SECRET = process.argv[3];
const PORT = 8888;
const REDIRECT_URI = `https://127.0.0.1:${PORT}/callback`;

if (!CLIENT_ID || !CLIENT_SECRET) {
  console.error('Usage: node get_token.js <CLIENT_ID> <CLIENT_SECRET>');
  process.exit(1);
}

// Generate a self-signed certificate for local HTTPS
const { privateKey, publicKey } = crypto.generateKeyPairSync('ec', {
  namedCurve: 'prime256v1'
});

const cert = createSelfSignedCert(privateKey, publicKey);

function createSelfSignedCert(privKey, pubKey) {
  // Use Node's built-in X509Certificate support via a child process with openssl,
  // or fall back to the simple approach using tls options directly.
  // Node 15+ supports createSecureContext with key/cert generated on the fly.
  // We'll use a minimal approach with the SubtleCrypto-free method.
  const { execSync } = require('child_process');
  const fs = require('fs');
  const os = require('os');
  const path = require('path');

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'navisync-'));
  const keyFile = path.join(tmpDir, 'key.pem');
  const certFile = path.join(tmpDir, 'cert.pem');

  try {
    execSync(
      `openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 ` +
      `-keyout "${keyFile}" -out "${certFile}" -days 1 -nodes -subj "/CN=127.0.0.1" ` +
      `-addext "subjectAltName=IP:127.0.0.1" 2>/dev/null`
    );
    return {
      key: fs.readFileSync(keyFile),
      cert: fs.readFileSync(certFile)
    };
  } finally {
    try { fs.unlinkSync(keyFile); } catch {}
    try { fs.unlinkSync(certFile); } catch {}
    try { fs.rmdirSync(tmpDir); } catch {}
  }
}

const server = https.createServer({ key: cert.key, cert: cert.cert }, async (req, res) => {
  const url = new URL(req.url, `https://127.0.0.1:${PORT}`);

  if (url.pathname === '/callback') {
    const code = url.searchParams.get('code');
    if (code) {
      try {
        const tokenRes = await fetch('https://accounts.spotify.com/api/token', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Authorization': 'Basic ' + Buffer.from(CLIENT_ID + ':' + CLIENT_SECRET).toString('base64')
          },
          body: new URLSearchParams({
            grant_type: 'authorization_code',
            code: code,
            redirect_uri: REDIRECT_URI
          })
        });

        const data = await tokenRes.json();

        if (data.error) {
          res.writeHead(400, { 'Content-Type': 'text/html' });
          res.end(`<h1>Error</h1><p>${data.error}: ${data.error_description}</p>`);
          console.error(`Error: ${data.error}: ${data.error_description}`);
          server.close();
          process.exit(1);
        }

        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(`
          <h1>Success!</h1>
          <p>Copy this Refresh Token:</p>
          <textarea style="width: 500px; height: 100px;">${data.refresh_token}</textarea>
          <p>You can close this window.</p>
        `);

        console.log('\n--- REFRESH TOKEN ---');
        console.log(data.refresh_token);
        console.log('---------------------\n');

        server.close();
        process.exit(0);
      } catch (err) {
        res.writeHead(500, { 'Content-Type': 'text/html' });
        res.end('Error: ' + err.message);
      }
    }
  }
});

server.listen(PORT, '127.0.0.1', () => {
  const authUrl = `https://accounts.spotify.com/authorize?client_id=${CLIENT_ID}&response_type=code&redirect_uri=${encodeURIComponent(REDIRECT_URI)}&scope=playlist-read-private`;
  console.log('1. Add this redirect URI to your Spotify app settings:');
  console.log(`   ${REDIRECT_URI}\n`);
  console.log('2. Open this URL in your browser:');
  console.log(`   ${authUrl}\n`);
  console.log('(Your browser will warn about the self-signed certificate — click "Advanced" and proceed.)');
});
