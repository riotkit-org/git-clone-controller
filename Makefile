.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	mkdir -p .build
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o .build/git-clone-operator .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t git-clone-operator .

.PHONY: coverage
coverage:
	go test -v ./... -covermode=count -coverprofile=coverage.out || true
