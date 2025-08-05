
# GLOBAL ENVS
GOCACHE = $(shell go env GOCACHE)
GOMODCACHE = $(shell go env GOMODCACHE)
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

TAG ?= develop

.PHONY: setup
setup:
	go mod download

.PHONY: lint
lint:
	@echo "Check project with linters"
	go tool golangci-lint run

.PHONY: unit
unit:
	go test ./... -v -cover
	go test ./... -v -cover -covermode=count -coverprofile coverage.out -json > report.json
	go tool cover -html=coverage.out -o=coverage.html

.PHONY: integration
integration:
	go test -tags=integration ./test/... -v

.PHONY: fmt
fmt:
	@echo "Format project files"
	go fmt ./...

.PHONY: build
build:
	go build \
	  -o ./build/go-agents-$(TAG)-$(GOOS)-$(GOARCH) \
	  main.go

.PHONY: docker-build
docker-build:
	docker buildx build \
		--progress=plain \
		--platform linux/amd64 \
		--build-arg GOCACHE=$(GOCACHE) \
		--build-arg GOMODCACHE=$(GOMODCACHE) \
		--output "type=docker" \
		-t go-agents:latest \
		.

.PHONY: clean
clean:
	rm coverage.html coverage.out report.json

.PHONY: gen-proto
gen-proto:
	buf lint
	buf generate

.PHONY: generate
generate: gen-proto
