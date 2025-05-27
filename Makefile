
# GLOBAL ENVS
GOCACHE = $(shell go env GOCACHE)
GOMODCACHE = $(shell go env GOMODCACHE)
GOOS ?= $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

TAG ?= develop

.PHONY: lint unit integration fmt build docker-build setup docs

setup:
	go mod download

lint:
	@echo "Check project with linters"
	go tool golangci-lint run

unit: 
	go test ./... -v -cover
	go test ./... -v -cover -covermode=count -coverprofile coverage.out -json > report.json
	go tool cover -html=coverage.out -o=coverage.html

integration: 
	go test -tags=integration ./test/... -v
	
fmt:
	@echo "Format project files"
	go fmt ./...

build:
	go build \
	  -o ./build/sisu-$(TAG)-$(GOOS)-$(GOARCH) \
	  -ldflags="-X 'github.com/inditex/sisu/config/latest.Version=$(TAG)'" \
	  main.go

docker-build:
	docker buildx build \
		--progress=plain \
		--platform linux/amd64 \
		--build-arg GOCACHE=$(GOCACHE) \
		--build-arg GOMODCACHE=$(GOMODCACHE) \
		--output "type=docker" \
		-t sisu:latest \
		.

clean:
	rm coverage.html coverage.out report.json
	