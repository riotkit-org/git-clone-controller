.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=linux GOARCH=amd64 go build -o bin/git-clone-operator .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t simple-kubernetes-webhook:latest .

.PHONY: coverage
coverage:
	go test -v ./... -covermode=count -coverprofile=coverage.out || true
