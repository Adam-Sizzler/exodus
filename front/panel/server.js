import express from 'express';
import { createServer as createViteServer } from 'vite';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
const PORT = process.env.PORT || 9000;
const API_URL = process.env.VITE_API_URL || 'http://localhost:9243';
const apiTarget = new URL(API_URL);

async function startServer() {
  // Create Vite server in middleware mode for development
  const vite = await createViteServer({
    server: { middlewareMode: true },
    appType: 'spa'
  });

  app.use(vite.middlewares);

  // Add CORS middleware
  app.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Methods', 'GET, POST, PATCH, DELETE, OPTIONS');
    res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
    if (req.method === 'OPTIONS') {
      return res.sendStatus(204);
    }
    next();
  });

  app.use(express.json());

  // API proxy for all API requests
  app.all('/api/*', async (req, res) => {
    try {
      const http = await import('http');
      const apiPath = req.originalUrl;
      const method = req.method || 'GET';
      const body = method !== 'GET' && method !== 'HEAD' && req.body
        ? JSON.stringify(req.body)
        : undefined;

      const options = {
        hostname: apiTarget.hostname,
        port: apiTarget.port || 80,
        path: apiPath,
        method: method,
        headers: {
          'Content-Type': 'application/json',
          'Content-Length': body ? Buffer.byteLength(body) : 0,
          Authorization: req.headers.authorization || '',
        },
      };

      const proxyReq = http.request(options, (proxyRes) => {
        let data = '';
        proxyRes.on('data', chunk => data += chunk);
        proxyRes.on('end', () => {
          try {
            // Copy status code and headers
            res.status(proxyRes.statusCode);
            res.header('Content-Type', proxyRes.headers['content-type']);
            
            if (data) {
              res.json(JSON.parse(data));
            } else {
              res.send();
            }
          } catch (e) {
            res.status(500).json({ error: 'Invalid JSON from API', details: e.message });
          }
        });
      });

      proxyReq.on('error', (err) => {
        console.error('API Error:', err.message);
        res.status(500).json({ error: 'Failed to fetch data from API', details: err.message });
      });

      if (body) {
        proxyReq.write(body);
      }
      proxyReq.end();
    } catch (error) {
      console.error('Proxy Error:', error.message);
      res.status(500).json({ error: 'Proxy error', details: error.message });
    }
  });

  // Health check
  app.get('/api/health', (req, res) => {
    res.json({ status: 'ok' });
  });

  app.listen(PORT, () => {
    console.log(`V2Ray Stat dashboard running at http://localhost:${PORT}`);
    console.log(`Proxying API requests to ${API_URL}`);
  });
}

startServer();
