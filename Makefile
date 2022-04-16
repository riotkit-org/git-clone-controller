REGISTRY_PORT=50161
REGISTRY=gco-registry.localhost

.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	mkdir -p .build
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o .build/git-clone-operator .

.PHONY: coverage
coverage:
	go test -v ./... -covermode=count -coverprofile=coverage.out || true

.PHONY: minikube-build
k3d-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t ${REGISTRY}:${REGISTRY_PORT}/git-clone-operator:master .
	docker push ${REGISTRY}:${REGISTRY_PORT}/git-clone-operator:master

.PHONY: minikube-promote
k3d-promote: build k3d-build
	helm uninstall gitc -n git-clone-operator --wait || true
	cd helm/git-clone-operator && helm upgrade --install gitc . -n git-clone-operator --create-namespace --set image.repository=k3d-${REGISTRY}:${REGISTRY_PORT}/git-clone-operator --set image.tag=master

.PHONY: k3d
k3d:
	k3d registry create ${REGISTRY} --port ${REGISTRY_PORT}
	k3d cluster create riotkit --registry-use k3d-${REGISTRY}:${REGISTRY_PORT} --agents 1 -p "30080:30080@agent:0" -p "30081:30081@agent:0" -p "30050:30050@agent:0"
