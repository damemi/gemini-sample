import express from "express";
import { createProxyMiddleware } from "http-proxy-middleware";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const port = Number(process.env.PORT || 3000);
const backend = process.env.BACKEND_URL || "http://127.0.0.1:8080";

const app = express();
app.use(express.static(path.join(__dirname, "public")));

app.use(
  "/api",
  createProxyMiddleware({
    target: backend,
    changeOrigin: true,
  })
);

app.listen(port, "0.0.0.0", () => {
  console.log(`frontend on :${port} proxy -> ${backend}`);
});
