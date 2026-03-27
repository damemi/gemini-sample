IMAGE ?= docker.io/mikeodigos/gemini-sample-frontend
TAG ?= latest

.PHONY: build push run

build:
	docker build -t $(IMAGE):$(TAG) .

push: build
	docker push $(IMAGE):$(TAG)

run:
	docker run --rm -e BACKEND_URL=http://host.docker.internal:9090 -p 8080:8080 $(IMAGE):$(TAG)
