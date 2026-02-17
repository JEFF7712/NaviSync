const http = require('http');
const { URL } = require('url');
const { exec } = require('child_process');

const CLIENT_ID = process.argv[2];
const CLIENT_SECRET = process.argv[3];
const PORT = 8888;
const REDIRECT_URI = `http://localhost:${PORT}/callback`;

if (!CLIENT_ID || !CLIENT_SECRET) {
  console.error('Usage: node get_token.js <CLIENT_ID> <CLIENT_SECRET>');
  process.exit(1);
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://localhost:${PORT}`);
  
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
        
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(`
          <h1>Success!</h1>
          <p>Copy this Refresh Token:</p>
          <textarea style="width: 500px; height: 100px;">${data.refresh_token}</textarea>
        `);
        
        console.log('\n--- REFRESH TOKEN ---');
        console.log(data.refresh_token);
        console.log('---------------------\n');
        
        server.close();
        process.exit(0);
      } catch (err) {
        res.end('Error: ' + err.message);
      }
    }
  }
});

server.listen(PORT, () => {
  const authUrl = `https://accounts.spotify.com/authorize?client_id=${CLIENT_ID}&response_type=code&redirect_uri=${encodeURIComponent(REDIRECT_URI)}&scope=playlist-read-private`;
  console.log('Open this URL in your browser:');
  console.log(authUrl);
});
