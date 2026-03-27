# Gemini sample chat (Kubernetes)

Small demo: a **Node.js** frontend and **Go** backend in separate pods. The browser talks to the frontend; the frontend proxies `/api` to the backend; the backend calls the [Google Gemini API](https://ai.google.dev/gemini-api/docs) (`generateContent`).

Published images (default in manifests):

- `docker.io/mikeodigos/gemini-sample-frontend`
- `docker.io/mikeodigos/gemini-sample-backend`

## Layout

| Path | Purpose |
|------|---------|
| `frontend/` | Express static UI + HTTP proxy to the backend |
| `backend/` | Go service: `POST /api/chat`, `GET /healthz` (uses [`google.golang.org/genai`](https://pkg.go.dev/google.golang.org/genai)) |
| `k8s/` | Deployments and Services (no namespace in YAML — uses default unless you pass `-n`) |
| `Makefile` | Build, push, deploy, port-forward |

## Prerequisites

- Docker (for images)
- `kubectl` configured for your cluster
- A [Gemini API key](https://aistudio.google.com/apikey) in `GEMINI_API_KEY`

## Build and push images

From the repo root:

```bash
make build
make push
```

Override image tag:

```bash
make push TAG=v1
```

Per service:

```bash
make -C backend push
make -C frontend push
```

Log in to Docker Hub (or your registry) before `make push` if images are not public.

## Deploy

Kubernetes does not substitute shell variables into manifests, so the API key is stored in a **Secret** named `gemini-credentials` (key `GEMINI_API_KEY`).

**Recommended (one step):**

```bash
export GEMINI_API_KEY='your-key-here'
make deploy
```

That creates or updates the Secret, then applies everything under `k8s/`. Resources land in the **default** namespace unless you set `NS`:

```bash
NS=my-namespace make deploy
```

**Without Make**, after `export GEMINI_API_KEY` — **default** namespace:

```bash
kubectl create secret generic gemini-credentials \
  --from-literal=GEMINI_API_KEY="$GEMINI_API_KEY" \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f k8s/
```

Use another namespace by adding the same `-n my-namespace` to both commands (and to any later `kubectl apply -f k8s/`).

To reapply manifests only (Secret already exists):

```bash
make apply
# or: NS=staging make apply
```

## Use the UI

- **NodePort:** frontend Service uses node port **30080** — open `http://<node-ip>:30080`.
- **Port-forward:** `make port-forward` (or `NS=my-namespace make port-forward`), then open [http://127.0.0.1:3000](http://127.0.0.1:3000).

## Configuration

| Variable | Where | Purpose |
|----------|--------|---------|
| `GEMINI_API_KEY` | Secret → backend pod | Gemini API authentication |
| `GEMINI_MODEL` | `k8s/backend-deployment.yaml` | Model id (default `gemini-2.0-flash`) |
| `BACKEND_URL` | Frontend Deployment | In-cluster URL of backend Service (default `http://gemini-sample-backend:8080`) |

## API shape

`POST /api/chat` (JSON):

```json
{
  "messages": [
    { "role": "user", "content": "Hello" },
    { "role": "model", "content": "Hi there." }
  ]
}
```

Roles must be `user` or `model`. Response:

```json
{ "reply": "…" }
```

## Optional: secret template with envsubst

If you prefer generating the Secret from a file:

```bash
export GEMINI_API_KEY='your-key'
envsubst '$GEMINI_API_KEY' < k8s/templates/secret.envsubst | kubectl apply -f -
kubectl apply -f k8s/
# With a namespace: pipe to `kubectl -n my-ns apply -f -` and use the same `-n` on the second line.
```

(`envsubst` is from gettext; on macOS it is often `brew install gettext` and add it to `PATH`.)

## Makefile reference

Run `make help` for the current target list.
