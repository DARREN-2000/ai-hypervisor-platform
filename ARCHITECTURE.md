# AI Hypervisor Platform - System Architecture

## Executive Overview

The AI Hypervisor Platform is a distributed, cloud-native virtualization system designed specifically for GPU-accelerated AI inference workloads. It orchestrates KVM/QEMU virtual machines across a Kubernetes cluster, providing automated provisioning, GPU allocation, workload isolation, and real-time monitoring.

### Core Design Principles

- **Cloud-Native First**: Kubernetes-integrated, containerized components, distributed deployment
- **Infrastructure as Code**: All infrastructure defined declaratively in Kubernetes manifests
- **Async-Driven**: Event-driven architecture with distributed task processing
- **Observable by Default**: Comprehensive metrics, logs, and traces built into every component
- **Secure by Design**: Workload isolation, RBAC, encrypted communication, audit logging
- **Scalable Architecture**: Horizontal scaling of control plane and VM hosts

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │          Control Plane Services (Namespace: aihypervisor)│   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │                                                            │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌──────────────────┐ │   │
│  │  │ API Server  │  │  VM Manager │  │ GPU Orchestrator│ │   │
│  │  │  (REST/WS)  │  │   Service   │  │    Service       │ │   │
│  │  └─────────────┘  └─────────────┘  └──────────────────┘ │   │
│  │                                                            │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │   │
│  │  │ Scheduler    │  │ Task Executor│  │ Resource     │   │   │
│  │  │ Service      │  │ Service      │  │ Monitor      │   │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘   │   │
│  │                                                            │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │   │
│  │  │ Secret       │  │ Config       │  │ Audit        │   │   │
│  │  │ Manager      │  │ Manager      │  │ Logger       │   │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘   │   │
│  │                                                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              VM Host Agents (Daemonset)                  │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │                                                            │   │
│  │  ┌────────────┐  ┌────────────┐  ┌───────────────────┐  │   │
│  │  │ Libvirt    │  │ GPU        │  │ VM Lifecycle      │  │   │
│  │  │ Client     │  │ Monitor    │  │ Manager           │  │   │
│  │  └────────────┘  └────────────┘  └───────────────────┘  │   │
│  │                                                            │   │
│  │  ┌────────────┐  ┌────────────┐  ┌───────────────────┐  │   │
│  │  │ Storage    │  │ Network    │  │ Telemetry         │  │   │
│  │  │ Adapter    │  │ Manager    │  │ Collector         │  │   │
│  │  └────────────┘  └────────────┘  └───────────────────┘  │   │
│  │                                                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         Infrastructure Services (Namespace: infra)       │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │                                                            │   │
│  │  PostgreSQL (State) │ Redis (Tasks) │ NATS (Events)      │   │
│  │  Prometheus         │ Grafana       │ Alertmanager       │   │
│  │                                                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │      AI Workload Tenants (Namespace: ai-workloads)       │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │   Ollama/vLLM Pods │ Inference Endpoints │ Model Cache   │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Service Architecture

### 1. API Server Service
**Responsibility**: REST/WebSocket API gateway for VM lifecycle management
- Exposes REST endpoints for VM CRUD operations
- WebSocket support for real-time VM status updates
- Request validation and authentication
- Rate limiting and quota enforcement
- Backward-compatible versioning

**Dependencies**: PostgreSQL, Redis, VM Manager Service
**Deployment**: Multi-replica StatefulSet with load balancing

### 2. VM Manager Service
**Responsibility**: Central VM lifecycle orchestration
- Manages VM state machine transitions
- Coordinates VM provisioning and teardown
- Enforces resource constraints
- Handles VM-to-node scheduling
- Manages VM metadata and configurations

**Dependencies**: Libvirt, GPU Orchestrator, Scheduler, PostgreSQL
**Deployment**: Single-replica StatefulSet with leader election

### 3. GPU Orchestrator Service
**Responsibility**: Intelligent GPU allocation and management
- Tracks GPU availability across cluster
- Implements GPU allocation algorithms
- Enforces GPU isolation and sharing policies
- Collects GPU telemetry
- Handles GPU affinity for NUMA systems

**Dependencies**: Prometheus, VM Host Agents, PostgreSQL
**Deployment**: Multi-replica Deployment with shared state

### 4. Scheduler Service
**Responsibility**: VM placement and resource scheduling
- Bin-packing algorithm for VM placement
- Network affinity optimization
- Storage locality awareness
- Anti-affinity policies
- Fallback scheduling strategies

**Dependencies**: PostgreSQL, GPU Orchestrator, Resource Monitor
**Deployment**: Multi-replica Deployment with consistent scheduling

### 5. Task Executor Service
**Responsibility**: Async background task execution
- Processes VM provisioning tasks
- Executes maintenance operations
- Handles retry logic with exponential backoff
- Task completion tracking
- Dead letter queue management

**Dependencies**: Redis/NATS, PostgreSQL, Libvirt
**Deployment**: Horizontally scalable Deployment

### 6. Resource Monitor Service
**Responsibility**: Real-time resource tracking and metrics
- Aggregates host resource metrics
- Tracks VM resource utilization
- Monitors GPU utilization
- Detects resource contention
- Provides scheduling hints

**Dependencies**: Prometheus, VM Host Agents, PostgreSQL
**Deployment**: Multi-replica Deployment with metric aggregation

### 7. VM Host Agent (Daemonset)
**Responsibility**: Node-level VM and infrastructure management
- Libvirt daemon management
- GPU monitoring on host
- VM lifecycle implementation (KVM/QEMU)
- Local storage management
- Network interface management
- Telemetry collection and reporting

**Dependencies**: Libvirt, KVM/QEMU, Host OS
**Deployment**: DaemonSet on all worker nodes with GPU support

### 8. Secret Manager Service
**Responsibility**: Secure credential and configuration management
- VM image authentication
- SSH key management
- API credential rotation
- Encryption key management
- Audit trail of secret access

**Dependencies**: PostgreSQL, Kubernetes Secrets API
**Deployment**: Multi-replica Deployment with consensus

### 9. Config Manager Service
**Responsibility**: Configuration distribution and versioning
- VM image catalog management
- VM flavor definitions
- Network policies
- GPU allocation policies
- Runtime configuration versioning

**Dependencies**: PostgreSQL, Redis, Kubernetes ConfigMap
**Deployment**: Multi-replica Deployment with caching layer

### 10. Audit Logger Service
**Responsibility**: Compliance and security auditing
- Logs all state-changing operations
- Tracks resource access patterns
- Maintains audit trail for compliance
- Generates compliance reports
- Integrates with SIEM systems

**Dependencies**: PostgreSQL, Kubernetes Audit API
**Deployment**: Multi-replica Deployment with persistent storage

## Data Flow Architecture

### VM Provisioning Flow
```
1. User → API Server (create VM request)
2. API Server → Database (persist intent)
3. API Server → Task Executor (queue provisioning task)
4. Task Executor → Scheduler (find suitable host)
5. Scheduler → GPU Orchestrator (allocate GPU)
6. GPU Orchestrator → Resource Monitor (verify capacity)
7. Task Executor → VM Host Agent (create VM via Libvirt)
8. VM Host Agent → Libvirt/KVM → VM Creation
9. VM Host Agent → Telemetry Collector (report status)
10. API Server → WebSocket Clients (notify completion)
```

### GPU Allocation Flow
```
1. Scheduler requests GPU allocation
2. GPU Orchestrator checks GPU availability
3. GPU Orchestrator applies allocation policy
4. GPU Orchestrator updates GPU assignment table
5. Scheduler proceeds with VM placement
6. VM Host Agent applies GPU device mapping
7. GPU Monitor tracks utilization
8. Resource Monitor aggregates metrics
```

### Monitoring Flow
```
1. VM Host Agent collects metrics (CPU, Memory, GPU, Network)
2. Telemetry Collector aggregates metrics
3. Prometheus scrapes metrics endpoint
4. Alertmanager evaluates alert rules
5. Grafana visualizes dashboards
6. Audit Logger records state changes
7. Event stream (NATS/Redis) publishes updates
```

## State Management

### VM State Machine
```
┌─────────────────────────────────────────────────────┐
│                                                       │
│  Created → Provisioning → Running → Scaling        │
│    │           ↓            ↓        ↓              │
│    └───────────→ Failed ←────┴────────┴──────┐      │
│                   ↓                           │      │
│                Deleting → Deleted             │      │
│                   ↑                           │      │
│                   └───────────────────────────┘      │
│                                                       │
│  States:                                             │
│  - Created: VM definition stored, not yet created   │
│  - Provisioning: Creating VM infrastructure        │
│  - Running: VM operational, workloads can run      │
│  - Scaling: VM resources being modified            │
│  - Failed: Error during creation/operation         │
│  - Deleting: VM being removed                      │
│  - Deleted: VM fully removed                       │
│                                                       │
└─────────────────────────────────────────────────────┘
```

### GPU Allocation States
```
Available → Allocated → Assigned → In-Use → Released → Available
            ↓           ↓           ↓        ↓
            └──→ Faulty ←──────────→ Error
```

## Database Schema (High-Level)

### Core Tables
- `vms`: VM definitions, state, metadata
- `vm_instances`: Running VM instances with lifecycle info
- `gpus`: GPU inventory and availability
- `gpu_allocations`: GPU-to-VM mappings
- `vm_hosts`: Physical hosts/nodes
- `tasks`: Background task queue and history
- `audit_logs`: Compliance and audit trail
- `configurations`: System and resource configurations
- `secrets`: Encrypted secrets and credentials

## API Contract Overview

### REST Endpoints
```
VMs:
  POST   /api/v1/vms              - Create VM
  GET    /api/v1/vms              - List VMs
  GET    /api/v1/vms/{id}         - Get VM details
  PATCH  /api/v1/vms/{id}         - Update VM
  DELETE /api/v1/vms/{id}         - Delete VM
  POST   /api/v1/vms/{id}/start   - Start VM
  POST   /api/v1/vms/{id}/stop    - Stop VM
  POST   /api/v1/vms/{id}/reboot  - Reboot VM

GPUs:
  GET    /api/v1/gpus             - List GPUs
  GET    /api/v1/gpus/{id}        - Get GPU details
  GET    /api/v1/gpus/allocation  - Get allocation status

Hosts:
  GET    /api/v1/hosts            - List hosts
  GET    /api/v1/hosts/{id}       - Get host details
  GET    /api/v1/hosts/{id}/metrics - Get host metrics

Monitoring:
  GET    /api/v1/metrics          - Get system metrics
  GET    /api/v1/health           - Health check

WebSocket:
  WS     /ws/vm/{id}/logs         - VM console output
  WS     /ws/vm/{id}/metrics      - VM metrics stream
  WS     /ws/cluster/events       - Cluster events stream
```

## Observability Strategy

### Metrics (Prometheus)
- VM resource utilization (CPU, Memory, Network, Disk I/O)
- GPU utilization and temperature
- API request latency and throughput
- Task execution duration and success rates
- Database query performance
- Scheduler decision metrics
- GPU allocation efficiency metrics

### Logging
- Structured JSON logging with correlation IDs
- Separate log streams: application, audit, infrastructure
- Log levels: DEBUG, INFO, WARN, ERROR, FATAL
- Centralized log aggregation (ELK/Loki)

### Tracing
- Distributed tracing with OpenTelemetry
- Trace each VM provisioning request end-to-end
- Correlate traces across service boundaries
- Trace GPU allocation decisions

### Alerting
- VM provisioning failure rates > threshold
- GPU allocation failures
- Host resource exhaustion
- API server latency p99 > threshold
- Database replication lag
- Task queue backlog growth

## Security Architecture

### Authentication & Authorization
- Service-to-service: mTLS with certificate rotation
- API Clients: JWT tokens with API keys
- Kubernetes RBAC for service accounts
- Audit logging for all authorization decisions

### Data Security
- Encryption at rest: PostgreSQL encryption, encrypted secrets
- Encryption in transit: TLS for all communication
- Secret management: Kubernetes Secrets, sealed secrets for backup
- SSH keys: Rotated regularly, access logged

### VM Isolation
- Network: CNI plugin with network policies
- Storage: Separate PVCs per VM
- Compute: Dedicated libvirt domains per VM
- GPU: GPU isolation via kernel drivers

### Compliance
- Audit trail for all state changes
- Immutable audit logs
- Compliance reports generation
- Access control policies
- Data retention policies

## Scaling Considerations

### Horizontal Scaling
- API Server: Stateless, scales with load
- GPU Orchestrator: Distributed state with consensus
- Task Executor: Work queue distribution
- VM Host Agents: One per node, scales with cluster

### Performance Optimization
- Scheduler caching of resource state
- GPU allocation state replication
- Task batch processing
- Connection pooling to PostgreSQL
- Redis caching layer

### Resource Limits
- Max VMs per host: 50-100 (configurable)
- Max GPUs per host: 8 (hardware dependent)
- Max concurrent provisioning tasks: Configurable queue depth
- API rate limiting: Token bucket per tenant

## Deployment Architecture

### Kubernetes Namespaces
- `aihypervisor`: Core control plane services
- `aihypervisor-agents`: VM host agents (DaemonSet)
- `infra`: Infrastructure services (DB, monitoring)
- `ai-workloads`: Tenant AI inference workloads

### Image Strategy
- Multi-stage builds for Go services
- Distroless images for minimal surface area
- Image scanning and vulnerability assessment
- Private registry with access control

### Configuration Management
- Helm charts for templating
- ConfigMaps for application config
- Secrets for credentials
- Environment-specific values overrides

## Disaster Recovery

### Backup Strategy
- Regular PostgreSQL snapshots
- VM state exports to external storage
- Configuration backups to version control
- Encryption key backup (separate secure location)

### Recovery Procedures
- RTO (Recovery Time Objective): 1 hour
- RPO (Recovery Point Objective): 15 minutes
- Failover to backup cluster: Automated
- Data validation before restoration

## Roadmap Extensions

1. **Multi-tenancy v2**: Tenant-level resource quotas and isolation
2. **GPU Virtualization**: MIG profiles for NVIDIA GPUs
3. **Live Migration**: VM live migration between hosts
4. **Advanced Scheduling**: Machine learning-based placement
5. **Cost Optimization**: Spot instance integration
6. **Multi-cluster Federation**: Cross-cluster VM orchestration
7. **GitOps Integration**: ArgoCD for configuration management
8. **Policy Engine**: OPA/Rego for compliance policies

---

This architecture document serves as the blueprint for all implementation decisions and is maintained as a living document throughout the project lifecycle.
