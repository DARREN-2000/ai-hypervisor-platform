# Developer Onboarding

Welcome to the AI Hypervisor Platform! This guide will help you set up your environment and understand our development practices.

## Prerequisites

- Go 1.21+
- Docker and Docker Compose (for local dependencies)
- Make
- pnpm (for frontend development)

## Local Setup

1. **Clone the repo:**
   ```bash
   git clone https://github.com/ai-hypervisor/platform.git
   cd platform
   ```

2. **Install Go tools:**
   ```bash
   make setup-dev
   ```

3. **Start local infrastructure:**
   ```bash
   docker-compose up -d postgres redis nats
   ```

4. **Initialize the database:**
   ```bash
   psql -h localhost -U postgres -d aihypervisor -f deploy/scripts/init-db.sql
   ```

5. **Run the services:**
   You can run services individually using `go run cmd/<service>/main.go` or build them using `make build`.

## Coding Conventions

- We follow standard Go formatting (`gofmt`).
- Use `make lint` before submitting a PR.
- Document all exported functions and types.
- We use the `pkg/errors` package for consistent error handling and wrapping.

## Testing

We require high test coverage for core components.

- **Run all tests:** `make test`
- **Run with coverage:** `go test -coverprofile=coverage.out ./...`

*Note: Ensure you do not commit `coverage.out`.*

## CI/CD and GitHub Actions

Our CI pipeline automatically runs tests, linters, and builds Docker images on every push.
See `.github/workflows/` for definitions.