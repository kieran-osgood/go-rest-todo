import { createProxyMiddleware } from 'http-proxy-middleware';

// Create proxy instance outside of request handler function to avoid unnecessary re-creation
const apiProxy = createProxyMiddleware({
  target: 'http://localhost:8080',
  changeOrigin: true,
  // pathRewrite: { [`^/api/proxy`]: '' },
  secure: false,
  router: {
    // when request.headers.host == 'dev.localhost:3000',
    // override target 'http://www.example.org' to 'http://localhost:8000'
    'localhost:3000': 'http://localhost:8080',
  },
});

export default function (req, res) {
  apiProxy(req, res, (result) => {
    if (result instanceof Error) {
      throw result;
    }

    throw new Error(`Request '${req.url}' is not proxied! We should never reach here!`);
  });
};