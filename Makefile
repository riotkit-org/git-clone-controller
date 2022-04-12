.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t simple-kubernetes-webhook:latest .

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster --config dev/manifests/kind/kind.cluster.yaml

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing admission-webhook image into Kind's Docker daemon..."
	kind load docker-image simple-kubernetes-webhook:latest

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	kubectl apply -f dev/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f dev/manifests/cluster-config/

.PHONY: deploy
deploy: push delete deploy-config
	@echo "\n🚀 Deploying simple-kubernetes-webhook..."
	kubectl apply -f dev/manifests/webhook/

.PHONY: delete
delete:
	@echo "\n♻️  Deleting simple-kubernetes-webhook deployment if existing..."
	kubectl delete -f dev/manifests/webhook/ || true

.PHONY: logs
logs:
	@echo "\n🔍 Streaming simple-kubernetes-webhook logs..."
	kubectl logs -l app=simple-kubernetes-webhook -f
