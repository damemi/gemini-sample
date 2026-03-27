import express from "express";
import { createProxyMiddleware } from "http-proxy-middleware";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const port = Number(process.env.PORT || 8080);
const backend =
  process.env.BACKEND_URL || "http://127.0.0.1:9090";

const app = express();
// Proxy /api before static files. Do not mount at "/api": Express strips the prefix, so the
// upstream would see "/chat" instead of "/api/chat" and the backend returns 404.
app.use(
  createProxyMiddleware({
    target: backend,
    changeOrigin: true,
    pathFilter: "/api",
  })
);
app.use(express.static(path.join(__dirname, "public")));

app.listen(port, "0.0.0.0", () => {
  console.log(`frontend on :${port} proxy -> ${backend}`);
});
