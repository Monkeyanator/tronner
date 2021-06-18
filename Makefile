GO=go
IMAGE=tronner
TAG=latest

.PHONY: docker.build
docker.build:
	docker build -t $(IMAGE):$(TAG) .

.PHONY: server.build
server.build:
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -o server ./cmd/server

.PHONY: wasm.build
wasm.build:
	CGO_ENABLED=0 GOOS=js GOARCH=wasm $(GO) build -a -o client.wasm ./cmd/client
	mv ./client.wasm ./static

.PHONY: serve
serve:
	$(GO) run ./cmd/server/main.go

.PHONY: test
test:
	$(GO) test ./...
