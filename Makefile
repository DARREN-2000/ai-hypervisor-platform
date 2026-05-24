.PHONY: help build test clean deploy-local deploy-staging deploy-prod docker-build docker-push kind-create kind-delete lint fmt vet

# Default target
help:
	@echo "AI Hypervisor Platform - Build Targets"
	@echo ""
	@echo "Build and Test:"
	@echo "  make build              Build all services"
	@echo "  make test               Run all tests"
	@echo "  make test-integration   Run integration tests"
	@echo "  make test-e2e           Run end-to-end tests"
	@echo "  make clean              Clean build artifacts"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint               Run linters"
	@echo "  make fmt                Format code"
	@echo "  make vet                Run go vet"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build       Build Docker images"
	@echo "  make docker-push        Push images to registry"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make kind-create        Create local KinD cluster"
	@echo "  make kind-delete        Delete local KinD cluster"
	@echo "  make deploy-local       Deploy to local cluster"
	@echo "  make deploy-staging     Deploy to staging"
	@echo "  make deploy-prod        Deploy to production"
	@echo ""
	@echo "Development:"
	@echo "  make run-api-server     Run API server locally"
	@echo "  make run-vm-manager     Run VM manager locally"
	@echo ""

# Build variables
SERVICES := api-server vm-manager gpu-orchestrator scheduler task-executor resource-monitor host-agent
REGISTRY ?= aihypervisor
VERSION ?= latest
GOFLAGS := -v -race

# Build all services
build:
	@echo "Building AI Hypervisor Platform services..."
	@mkdir -p bin
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		CGO_ENABLED=1 go build $(GOFLAGS) -o bin/$$service ./cmd/$$service || exit 1; \
	done
	@echo "Build complete!"

# Build specific service
build-%:
	@echo "Building $*..."
	@mkdir -p bin
	CGO_ENABLED=1 go build $(GOFLAGS) -o bin/$* ./cmd/$*
	@echo "Built bin/$*"

# Run tests
test:
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -race -tags=integration ./test/integration/...

test-e2e:
	@echo "Running end-to-end tests..."
	go test -v -timeout=30m -tags=e2e ./test/e2e/...

# Code quality
lint:
	@echo "Running linters..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

vet:
	@echo "Running go vet..."
	go vet ./...

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ coverage.out
	go clean ./...

# Docker
docker-build:
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		docker build --build-arg SERVICE=$$service -f deploy/docker/Dockerfile -t $(REGISTRY)/$$service:$(VERSION) . || exit 1; \
	done

docker-push:
	@echo "Pushing Docker images..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service..."; \
		docker push $(REGISTRY)/$$service:$(VERSION) || exit 1; \
	done

# Kubernetes - Kind
kind-create:
	@echo "Creating KinD cluster..."
	kind create cluster --name aihypervisor --config deploy/scripts/kind-config.yaml
	@echo "Cluster created successfully!"

kind-delete:
	@echo "Deleting KinD cluster..."
	kind delete cluster --name aihypervisor

# Deployment
deploy-local: docker-build
	@echo "Deploying to local cluster..."
	kubectl apply -f deploy/kubernetes/manifests.yaml
	@echo "Deployment complete!"

deploy-staging:
	@echo "Deploying to staging cluster..."
	kubectl config use-context staging
	kubectl apply -f deploy/kubernetes/manifests.yaml
	@echo "Staging deployment complete!"

deploy-prod:
	@echo "Deploying to production cluster..."
	kubectl config use-context production
	kubectl apply -f deploy/kubernetes/manifests.yaml
	@echo "Production deployment complete!"

# Local development
run-api-server:
	@echo "Running API server..."
	go run ./cmd/api-server --config config/sample-config.yaml

run-vm-manager:
	@echo "Running VM manager..."
	go run ./cmd/vm-manager --config config/sample-config.yaml

# Development setup
setup-dev:
	@echo "Setting up development environment..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development environment ready!"

# Database migration
db-migrate:
	@echo "Running database migrations..."
	go run ./cmd/db-migrate --config config/sample-config.yaml

# Generate API documentation
api-docs:
	@echo "Generating API documentation..."
	@echo "API specification available at: docs/api/openapi.yaml"
	@echo "Open with: https://editor.swagger.io/?url=file:///$(PWD)/docs/api/openapi.yaml"

# Watch mode for development
watch:
	@echo "Starting watch mode..."
	go run github.com/cosmtrek/air@latest -c .air.toml

# Metrics
metrics:
	@echo "Collecting project metrics..."
	@echo "Code lines:"
	@find . -name "*.go" ! -path "./vendor/*" ! -path "./.git/*" | xargs wc -l | tail -1
	@echo "Test coverage:"
	@if [ -f coverage.out ]; then go tool cover -func=coverage.out | tail -1; fi

.PHONY: $(addprefix build-,$(SERVICES))
