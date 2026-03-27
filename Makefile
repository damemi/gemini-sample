# Images
BACKEND_IMAGE ?= docker.io/mikeodigos/gemini-sample-backend
FRONTEND_IMAGE ?= docker.io/mikeodigos/gemini-sample-frontend
TAG ?= latest

# Empty = default namespace; set NS=my-namespace for make targets
NS ?=
KUBECTL = kubectl$(if $(strip $(NS)), -n $(NS),)

.PHONY: help build build-backend build-frontend push push-backend push-frontend deploy secret apply port-forward

help:
	@echo "Targets:"
	@echo "  make build          - docker build both images"
	@echo "  make push           - build and push both images"
	@echo "  make deploy         - Secret from GEMINI_API_KEY, then kubectl apply -f k8s/"
	@echo "  make secret         - Secret only (needs GEMINI_API_KEY)"
	@echo "  make apply          - kubectl apply -f k8s/ (Secret must already exist)"
	@echo "  make port-forward   - kubectl port-forward to frontend :3000"
	@echo ""
	@echo "Namespace: omit for default; or NS=my-ns make deploy (same for apply, port-forward)"
	@echo ""
	@echo "Plain kubectl (no make), default namespace, after export GEMINI_API_KEY:"
	@echo "  kubectl create secret generic gemini-credentials --from-literal=GEMINI_API_KEY=\"\$$GEMINI_API_KEY\" --dry-run=client -o yaml | kubectl apply -f -"
	@echo "  kubectl apply -f k8s/"
	@echo "Other namespace: add -n my-ns to both kubectl lines (and use the same -n for apply)."

build: build-backend build-frontend

build-backend:
	$(MAKE) -C backend build IMAGE=$(BACKEND_IMAGE) TAG=$(TAG)

build-frontend:
	$(MAKE) -C frontend build IMAGE=$(FRONTEND_IMAGE) TAG=$(TAG)

push: push-backend push-frontend

push-backend:
	$(MAKE) -C backend push IMAGE=$(BACKEND_IMAGE) TAG=$(TAG)

push-frontend:
	$(MAKE) -C frontend push IMAGE=$(FRONTEND_IMAGE) TAG=$(TAG)

secret:
	@test -n "$$GEMINI_API_KEY" || (echo "error: set GEMINI_API_KEY in the environment" >&2 && exit 1)
	$(KUBECTL) create secret generic gemini-credentials \
		--from-literal=GEMINI_API_KEY="$$GEMINI_API_KEY" \
		--dry-run=client -o yaml | $(KUBECTL) apply -f -

apply:
	$(KUBECTL) apply -f k8s/

deploy: secret apply

port-forward:
	$(KUBECTL) port-forward svc/gemini-sample-frontend 3000:3000
