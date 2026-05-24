# AI Hypervisor Platform - Deployment Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Infrastructure Setup](#infrastructure-setup)
3. [Kubernetes Cluster Preparation](#kubernetes-cluster-preparation)
4. [Deploying Infrastructure Services](#deploying-infrastructure-services)
5. [Deploying AI Hypervisor Platform](#deploying-ai-hypervisor-platform)
6. [Verification](#verification)
7. [Post-Deployment Configuration](#post-deployment-configuration)
8. [Troubleshooting](#troubleshooting)

## Prerequisites

### Hardware Requirements

**Control Plane Nodes** (minimum 3 for HA):
- CPU: 4 cores
- RAM: 8 GB
- Disk: 50 GB (fast SSD)

**GPU Compute Nodes** (per your workload):
- CPU: 8+ cores
- RAM: 16+ GB
- Disk: 100+ GB (fast SSD)
- GPU: 1-8 NVIDIA/AMD GPUs (with drivers installed)

### Software Requirements

- Kubernetes 1.24+
- Docker 20.10+
- kubectl 1.24+
- Helm 3.10+
- KVM/QEMU support on compute nodes
- Libvirt 8.0+

### Network Requirements

- Persistent container networking (CNI plugin)
- Network policies support (Calico, Cilium, etc.)
- Load balancer for API Server access
- DNS for service discovery

## Infrastructure Setup

### 1. Setup Kubernetes Cluster

Using kubeadm (example):

```bash
# On each node, install container runtime (Docker)
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# On control plane
kubeadm init --pod-network-cidr=10.244.0.0/16

# Copy config
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Install CNI (Flannel example)
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

# On worker nodes
kubeadm join <control-plane-ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

Or use managed Kubernetes (EKS, GKE, AKS).

### 2. Install NVIDIA GPU Device Plugin

```bash
# For NVIDIA GPUs
kubectl apply -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v0.13.0/nvidia-device-plugin.yml

# Verify GPU detection
kubectl get nodes -o json | grep -A5 "nvidia.com/gpu"
```

### 3. Prepare GPU Nodes

```bash
# Label GPU nodes
kubectl label nodes <gpu-node-name> aihypervisor/gpu-node=true
kubectl label nodes <gpu-node-name> aihypervisor/compute=true

# On each GPU node, install KVM/QEMU
sudo apt-get update
sudo apt-get install -y qemu-kvm libvirt-daemon libvirt-clients \
  libvirt-daemon-system virt-manager bridge-utils

# Enable and start libvirt
sudo systemctl enable libvirtd
sudo systemctl start libvirtd

# Verify libvirt is running
virsh list

# Allow qemu URI access
sudo usermod -G libvirt $(whoami)
```

### 4. Create Required Namespaces

```bash
kubectl create namespace aihypervisor
kubectl create namespace aihypervisor-agents
kubectl create namespace ai-workloads
kubectl create namespace infra
kubectl create namespace monitoring
```

## Kubernetes Cluster Preparation

### Network Policies

Apply network policies to restrict traffic:

```bash
# Create network policy manifests
kubectl apply -f deploy/kubernetes/network-policies.yaml
```

### Storage Classes

```bash
# Create storage classes for persistent volumes
kubectl apply -f deploy/kubernetes/storage-classes.yaml
```

### RBAC Setup

Create ClusterRoles and ClusterRoleBindings:

```bash
# Applied as part of main deployment
kubectl apply -f deploy/kubernetes/manifests.yaml
```

## Deploying Infrastructure Services

### 1. PostgreSQL

```bash
# Add Bitnami Helm repo
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create secrets
kubectl create secret generic postgres-secret \
  --from-literal=password=your-secure-password \
  -n infra

# Install PostgreSQL
helm install postgres bitnami/postgresql \
  --namespace infra \
  --set auth.password=your-secure-password \
  --set auth.username=aihypervisor \
  --set auth.database=aihypervisor \
  --set primary.persistence.size=100Gi \
  --set primary.persistence.storageClassName=fast-ssd
```

Initialize database schema:

```bash
# Connect to PostgreSQL
kubectl -n infra port-forward svc/postgres-postgresql 5432:5432

# In another terminal, run migrations
psql postgresql://aihypervisor:password@localhost:5432/aihypervisor < deploy/scripts/init-db.sql
```

### 2. Redis

```bash
# Create secret
kubectl create secret generic redis-secret \
  --from-literal=password=your-secure-password \
  -n infra

# Install Redis
helm install redis bitnami/redis \
  --namespace infra \
  --set auth.password=your-secure-password \
  --set replica.replicaCount=2 \
  --set persistence.size=50Gi
```

### 3. NATS

```bash
# Add NATS Helm repo
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update

# Install NATS
helm install nats nats/nats \
  --namespace infra \
  --set nats.resources.limits.memory=1Gi \
  --set nats.resources.requests.memory=512Mi
```

### 4. Prometheus & Grafana

```bash
# Add Prometheus community Helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install kube-prometheus-stack
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.retention=15d \
  --set grafana.adminPassword=your-secure-password
```

### 5. Loki (Optional)

```bash
helm install loki grafana/loki-stack \
  --namespace monitoring \
  --set promtail.enabled=true
```

## Deploying AI Hypervisor Platform

### 1. Create ConfigMap

```bash
# Create ConfigMap from sample config
kubectl create configmap aihypervisor-config \
  --from-file=config/sample-config.yaml \
  -n aihypervisor
```

### 2. Apply Main Manifests

```bash
# Deploy all AI Hypervisor Platform services
kubectl apply -f deploy/kubernetes/manifests.yaml

# Verify deployments
kubectl -n aihypervisor get deployments
kubectl -n aihypervisor get statefulsets
kubectl -n aihypervisor-agents get daemonsets
```

### 3. Configure TLS (Optional but Recommended)

```bash
# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes

# Create secret
kubectl create secret tls aihypervisor-tls \
  --cert=tls.crt \
  --key=tls.key \
  -n aihypervisor

# Update manifests to reference TLS secret
kubectl patch deployment api-server -n aihypervisor \
  -p '{"spec":{"template":{"spec":{"volumes":[{"name":"tls","secret":{"secretName":"aihypervisor-tls"}}]}}}}'
```

## Verification

### 1. Check Pod Status

```bash
# All pods should be in Running state
kubectl -n aihypervisor get pods
kubectl -n aihypervisor-agents get pods
kubectl -n infra get pods
kubectl -n monitoring get pods
```

### 2. Check Service Health

```bash
# Check API Server health
kubectl -n aihypervisor port-forward svc/api-server 8080:80
curl http://localhost:8080/health

# Should return: {"status":"healthy","version":"1.0.0",...}
```

### 3. Verify Database Connectivity

```bash
# Connect to API Server pod
kubectl -n aihypervisor exec -it deployment/api-server -- /bin/bash

# Inside pod, test database connection
psql postgresql://aihypervisor:password@postgres.infra.svc.cluster.local:5432/aihypervisor -c "SELECT 1"
```

### 4. Check GPU Detection

```bash
# On GPU node
nvidia-smi

# Check GPU allocation from host agent
kubectl -n aihypervisor-agents logs -f daemonset/host-agent

# Look for "GPU detection" or "GPU available" messages
```

## Post-Deployment Configuration

### 1. Configure Ingress

```yaml
# Create ingress for API Server
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: aihypervisor-api
  namespace: aihypervisor
spec:
  rules:
    - host: api.aihypervisor.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-server
                port:
                  number: 80
```

### 2. Configure Monitoring

```bash
# Access Grafana
kubectl -n monitoring port-forward svc/prometheus-grafana 3000:80

# Default credentials: admin / your-password
# Add Prometheus data source pointing to: http://prometheus:9090
```

### 3. Create VM Flavors

```bash
# Create ConfigMap with VM flavors
kubectl create configmap vm-flavors \
  --from-file config/flavors.yaml \
  -n aihypervisor
```

### 4. Create VM Images

```bash
# Upload VM images to image registry
# Update image registry configuration in config.yaml
```

## Troubleshooting

### Pod Fails to Start

```bash
# Check pod events
kubectl describe pod <pod-name> -n aihypervisor

# Check logs
kubectl logs <pod-name> -n aihypervisor
kubectl logs <pod-name> -n aihypervisor --previous  # For crashed pods
```

### Database Connection Issues

```bash
# Verify database is running
kubectl -n infra get pods | grep postgres

# Check PostgreSQL logs
kubectl logs deployment/postgres -n infra

# Test connectivity from pod
kubectl -n aihypervisor run -it --image=postgres debug -- \
  psql postgresql://aihypervisor:password@postgres:5432/aihypervisor
```

### GPU Not Detected

```bash
# Verify NVIDIA device plugin is running
kubectl get pods -n kube-system | grep nvidia

# Check device plugin logs
kubectl logs -n kube-system -l k8s-app=nvidia-gpu-device-plugin

# Verify GPU availability
kubectl describe node <gpu-node>

# Look for "nvidia.com/gpu" in capacity
```

### API Server Connection Refused

```bash
# Check if API Server pod is running
kubectl -n aihypervisor get pods | grep api-server

# Check for resource constraints
kubectl describe pod <api-server-pod> -n aihypervisor

# Check logs
kubectl logs deployment/api-server -n aihypervisor
```

### High Memory/CPU Usage

```bash
# Check current resource usage
kubectl top nodes
kubectl top pods -n aihypervisor

# Increase resource limits if needed
kubectl set resources deployment api-server \
  -n aihypervisor \
  --limits=cpu=4,memory=2Gi \
  --requests=cpu=2,memory=1Gi
```

## Scaling

### Horizontal Scaling (Add More Nodes)

```bash
# Add worker nodes to cluster
# GPU nodes will auto-join daemon sets
# Control plane nodes increase availability
```

### Vertical Scaling (Increase Pod Resources)

```bash
# Update resource requests/limits in manifests.yaml
# Reapply deployment
kubectl apply -f deploy/kubernetes/manifests.yaml
```

### Storage Expansion

```bash
# Resize PersistentVolumeClaim
kubectl patch pvc postgres-postgresql \
  -p '{"spec":{"resources":{"requests":{"storage":"200Gi"}}}}' \
  -n infra
```

## Backup & Recovery

### Backup Database

```bash
# Create backup
kubectl -n infra exec postgres-0 -- \
  pg_dump -U aihypervisor aihypervisor > backup-$(date +%Y%m%d).sql

# Store backup securely
gsutil cp backup-*.sql gs://your-backup-bucket/
```

### Backup etcd

```bash
# On control plane node
sudo etcdctl --endpoints=127.0.0.1:2379 snapshot save etcd-backup.db

# Store backup
sudo gsutil cp etcd-backup.db gs://your-backup-bucket/
```

### Recovery Procedure

```bash
# Restore from database backup
kubectl -n infra exec -i postgres-0 -- \
  psql -U aihypervisor aihypervisor < backup-*.sql

# Verify data integrity
kubectl -n aihypervisor exec deployment/api-server -- \
  curl localhost:8080/api/v1/metrics
```

---

For additional support, see:
- [ARCHITECTURE.md](../ARCHITECTURE.md)
- [Operational Runbooks](operations/)
- [API Documentation](api/)
