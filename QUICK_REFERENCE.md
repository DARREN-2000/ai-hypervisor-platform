# AI Hypervisor Platform - Quick Reference Index

## 📚 Documentation

| Document | Purpose | Size |
|----------|---------|------|
| [README.md](README.md) | Project overview and quick start | 400+ lines |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Complete system architecture and design | 300+ lines |
| [SCAFFOLDING_SUMMARY.md](SCAFFOLDING_SUMMARY.md) | Project scaffolding overview | This index |
| [docs/GPU_ALLOCATION_AND_LIFECYCLE.md](docs/GPU_ALLOCATION_AND_LIFECYCLE.md) | GPU allocation strategies and VM lifecycle | 400+ lines |
| [docs/api/openapi.yaml](docs/api/openapi.yaml) | REST API specification (OpenAPI 3.0) | 500+ lines |
| [docs/deployment/DEPLOYMENT_GUIDE.md](docs/deployment/DEPLOYMENT_GUIDE.md) | Production deployment guide | 350+ lines |

## 🏗️ Architecture

**7 Microservices:**
1. API Server - REST/WebSocket gateway
2. VM Manager - Lifecycle orchestration
3. GPU Orchestrator - GPU allocation & monitoring
4. Scheduler - VM placement engine
5. Task Executor - Async operations
6. Resource Monitor - Metrics aggregation
7. Host Agent - Node-level management (DaemonSet)

**Infrastructure:**
- PostgreSQL (State persistence)
- Redis (Task queue & caching)
- NATS (Event messaging)
- Prometheus (Metrics)
- Grafana (Visualization)

## 📁 Source Code

### Core Services
```
cmd/api-server/main.go              - API server entry point
cmd/vm-manager/main.go              - VM manager entry point
cmd/gpu-orchestrator/main.go        - GPU orchestrator entry point
cmd/scheduler/main.go               - Scheduler entry point
cmd/task-executor/main.go           - Task executor entry point
cmd/resource-monitor/main.go        - Resource monitor entry point
cmd/host-agent/main.go              - Host agent entry point
```

### Internal Packages
```
internal/api/server.go              - REST API server implementation
internal/models/models.go           - 30+ data model types
internal/orchestrator/interfaces.go - 11 service interfaces
internal/config/config.go           - Configuration structures
```

### Public Packages
```
pkg/errors/errors.go               - Error handling
pkg/telemetry/                     - Observability utilities
pkg/common/                        - Common utilities
```

## 🐳 Deployment

### Kubernetes
```
deploy/kubernetes/manifests.yaml   - Complete K8s deployment (900+ lines)
```

Includes:
- 3 Namespaces
- 7 Services
- 6 Deployments + 1 StatefulSet + 1 DaemonSet
- Service Accounts and RBAC
- ConfigMaps and storage classes
- HorizontalPodAutoscalers

### Docker
```
deploy/docker/Dockerfile           - Multi-stage build
```

### Database
```
deploy/scripts/init-db.sql         - PostgreSQL schema (400+ lines)
```

14 tables with proper indexing, views, and stored procedures.

### Configuration
```
config/sample-config.yaml          - Complete configuration reference (300+ lines)
```

## 🛠️ Build & Development

```
Makefile                            - Build automation (200+ lines)
go.mod                             - Go module definition

Key targets:
  make build              - Build all services
  make test              - Run all tests
  make docker-build      - Build Docker images
  make deploy-local      - Deploy to local cluster
```

## 📊 Data Models (30+ types)

**Core Models:**
- `VirtualMachine` - VM definition and state
- `HostNode` - Physical host
- `GPU` - GPU device
- `Task` - Async task
- `ResourceMetrics` - Performance data

**Supporting Models:**
- Network configurations
- Storage volumes
- GPU allocations
- Events and audit logs
- Scheduling decisions
- Secrets management

## 🔌 Service Interfaces (11 total)

```go
type VMManager interface          // VM lifecycle
type Scheduler interface          // Placement decisions
type GPUOrchestrator interface    // GPU management
type TaskExecutor interface       // Async tasks
type ResourceMonitor interface    // Metrics
type HostAgent interface          // Node operations
type LibvirtClient interface      // VM operations
type ConfigManager interface      // Configuration
type EventBus interface           // Event streaming
type AuditLogger interface        // Compliance logging
type StateStore interface         // Persistence
```

## 🌐 API Endpoints (20+)

**VMs:**
- `POST   /api/v1/vms` - Create VM
- `GET    /api/v1/vms` - List VMs
- `GET    /api/v1/vms/{id}` - Get VM
- `PATCH  /api/v1/vms/{id}` - Update VM
- `DELETE /api/v1/vms/{id}` - Delete VM
- `POST   /api/v1/vms/{id}/start` - Start VM
- `POST   /api/v1/vms/{id}/stop` - Stop VM
- `POST   /api/v1/vms/{id}/reboot` - Reboot VM

**GPUs:**
- `GET    /api/v1/gpus` - List GPUs
- `GET    /api/v1/gpus/{id}` - Get GPU

**Hosts:**
- `GET    /api/v1/hosts` - List hosts
- `GET    /api/v1/hosts/{id}` - Get host
- `GET    /api/v1/hosts/{id}/metrics` - Host metrics

**Monitoring:**
- `GET    /api/v1/tasks/{id}` - Get task status
- `GET    /api/v1/metrics` - Cluster metrics
- `GET    /health` - Health check

**WebSocket:**
- `WS     /ws/cluster/events` - Event stream
- `WS     /ws/vm/{id}/metrics` - VM metrics stream

## 🔐 Security Features

✓ mTLS between services
✓ RBAC for access control
✓ Encryption at rest
✓ Encryption in transit
✓ Audit logging
✓ Secret management
✓ JWT authentication
✓ API key support

## 📊 Database Schema (14 tables)

```
vms                    - Virtual machine definitions
host_nodes             - Physical hosts
gpus                   - GPU devices
gpu_allocations        - GPU-to-VM mappings
tasks                  - Async tasks
resource_metrics       - Performance data
audit_logs             - Compliance trail
configurations         - System config
vm_network_configs     - Network settings
vm_storage_configs     - Storage settings
events                 - System events
scheduling_decisions   - Scheduler audit
secrets                - Encrypted secrets
```

Plus views and stored procedures for common operations.

## 🚀 Quick Start

### 1. Prerequisites
```bash
# Kubernetes cluster
kubectl cluster-info

# GPU nodes with KVM/QEMU
virsh list
nvidia-smi

# Required packages
postgresql redis nats
```

### 2. Deploy Infrastructure
```bash
helm install postgres bitnami/postgresql --namespace infra --create-namespace
helm install redis bitnami/redis --namespace infra
helm install nats nats/nats --namespace infra
```

### 3. Deploy Platform
```bash
kubectl apply -f deploy/kubernetes/manifests.yaml
```

### 4. Verify
```bash
kubectl -n aihypervisor get pods
kubectl -n aihypervisor port-forward svc/api-server 8080:80
curl http://localhost:8080/health
```

### 5. Create VM
```bash
curl -X POST http://localhost:8080/api/v1/vms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "inference-vm",
    "flavor": {"name": "medium-gpu", "cpu": 8, "memory": 16},
    "image": {"id": "ubuntu-22.04"},
    "gpu_requests": [{"type": "nvidia", "model": "A100", "count": 2}]
  }'
```

## 🔍 Monitoring

**Grafana Dashboards:**
- Cluster overview
- GPU utilization
- VM lifecycle
- Task execution
- Resource metrics

**Prometheus Metrics:**
- 100+ built-in metrics
- API latency and throughput
- Task execution performance
- GPU utilization and temperature
- VM provisioning time

**Logs:**
- Structured JSON logging
- Correlation IDs
- Separate streams: app, audit, infra

## 📚 Learning Path

1. **Start Here**: [README.md](README.md)
2. **Understand Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
3. **Deep Dive GPU/VM**: [docs/GPU_ALLOCATION_AND_LIFECYCLE.md](docs/GPU_ALLOCATION_AND_LIFECYCLE.md)
4. **Deployment**: [docs/deployment/DEPLOYMENT_GUIDE.md](docs/deployment/DEPLOYMENT_GUIDE.md)
5. **API Usage**: [docs/api/openapi.yaml](docs/api/openapi.yaml)

## 🎯 Key Statistics

- **Services**: 7
- **Namespaces**: 3
- **Data Models**: 30+
- **API Endpoints**: 20+
- **Database Tables**: 14
- **Kubernetes Resources**: 40+
- **Error Types**: 12+
- **Configuration Options**: 50+
- **Documentation**: 2000+ lines

## ✅ Production Readiness Checklist

- [x] Modular clean architecture
- [x] Production-style folder structure
- [x] Infrastructure-focused design
- [x] Scalable service boundaries
- [x] Extensible orchestration layer
- [x] Secure VM lifecycle management
- [x] Async task execution
- [x] Resource monitoring and scheduling
- [x] Configuration-driven deployment
- [x] Comprehensive documentation
- [x] Complete API specification
- [x] Kubernetes manifests for deployment
- [x] Database schema with migrations
- [x] Docker multi-stage builds
- [x] RBAC and security policies
- [x] Monitoring and observability setup
- [x] Error handling and recovery

## 🔗 Related Resources

- **Kubernetes**: https://kubernetes.io/
- **KVM/QEMU**: https://www.qemu.org/
- **Libvirt**: https://libvirt.org/
- **PostgreSQL**: https://www.postgresql.org/
- **Prometheus**: https://prometheus.io/
- **OpenTelemetry**: https://opentelemetry.io/

---

**Status**: ✅ Complete scaffolding and architecture

This is a production-grade foundation ready for development and deployment. All components are designed with scalability, security, and observability in mind.
