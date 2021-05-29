GO=go
IMAGE=tronner
TAG=latest

.PHONY: docker.build
docker.build:
	docker build -t $(IMAGE):$(TAG) .

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -o main ./cmd

.PHONY: serve
serve:
	$(GO) run ./cmd

test:
	$(GO) test -v ./...