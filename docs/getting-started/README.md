# Getting Started

This guide will walk you through deploying the AI Hypervisor Platform for the first time.

## 1. Prerequisites

You need a Kubernetes cluster (e.g., Minikube, kind, or a managed service) and `kubectl` configured.

## 2. Deploy Infrastructure

The platform requires PostgreSQL, Redis, and NATS. You can deploy these using Helm or the provided manifests.

```bash
# Example using provided manifests (for testing only)
kubectl apply -f deploy/kubernetes/manifests.yaml
```

## 3. Configuration

Review and copy the `config/sample-config.yaml` to `config/config.yaml`. Update the database connection strings and NATS URLs to point to your deployed infrastructure.

## 4. Deploy Platform Services

Apply the AI Hypervisor manifests.

```bash
kubectl apply -f deploy/kubernetes/platform.yaml
```

## 5. Verify

Check that all pods are running:

```bash
kubectl get pods -n ai-hypervisor
```