REGISTRY_PORT=50161
REGISTRY=gco-registry.localhost

.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./... -covermode=count -coverprofile=coverage.out

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	mkdir -p .build
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o .build/git-clone-controller .
