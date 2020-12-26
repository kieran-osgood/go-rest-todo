import { createProxyMiddleware } from 'http-proxy-middleware';

// Create proxy instance outside of request handler function to avoid unnecessary re-creation
const apiProxy = createProxyMiddleware({
  target: process.env.NEXT_PUBLIC_API_URL,
  changeOrigin: true,
  // pathRewrite: { [`^/api/proxy`]: '' },
  secure: false,
  router: {
    // when request.headers.host == 'dev.localhost:3000',
    // override target 'http://www.example.org' to 'http://localhost:8000'
    'localhost:3000': process.env.NEXT_PUBLIC_API_URL,
  },
});

export default (req, res) => {
  apiProxy(req, res, (result) => {
    if (result instanceof Error) {
      throw result;
    }

    throw new Error(`Request '${req.url}' is not proxied! We should never reach here!`);
  });
};
