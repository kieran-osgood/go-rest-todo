import { createProxyMiddleware } from 'http-proxy-middleware';

const proxy = createProxyMiddleware({
  target: process.env.NEXT_PUBLIC_API_URL,
  changeOrigin: true,
  secure: false,
  router: {
    // when request.headers.host == 'dev.localhost:3000',
    // override target 'http://www.example.org' to 'http://localhost:8000'
    'localhost:3000': process.env.NEXT_PUBLIC_API_URL,
  },
});

export default proxy;
