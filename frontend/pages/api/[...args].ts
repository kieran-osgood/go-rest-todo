import { createProxyMiddleware } from 'http-proxy-middleware';

// restream parsed body before proxying
const proxy = createProxyMiddleware({
  target: process.env.NEXT_PUBLIC_API_URL,
  changeOrigin: true,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
  onProxyReq: (proxyReq, req) => {
    if (req.body) {
      const bodyData = JSON.stringify(req.body);
      // incase if content-type is application/x-www-form-urlencoded ->
      // we need to change to application/json
      proxyReq.setHeader('Content-Type', 'application/json');
      proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
      // stream the content
      proxyReq.write(bodyData);
    }
  },

  router: {
    // when request.headers.host == 'dev.localhost:3000',
    // override target 'http://www.example.org' to 'http://localhost:8000'
    'localhost:3000': process.env.NEXT_PUBLIC_API_URL,
  },
});

export default proxy;
