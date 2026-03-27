IMAGE ?= docker.io/mikeodigos/gemini-sample-backend
TAG ?= latest

.PHONY: build push run

build:
	docker build -t $(IMAGE):$(TAG) .

push: build
	docker push $(IMAGE):$(TAG)

run:
	docker run --rm -e GEMINI_API_KEY -p 9090:9090 $(IMAGE):$(TAG)
