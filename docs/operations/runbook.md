# AI Hypervisor Platform Runbook

This runbook covers the minimum operational steps for deploying, observing, and recovering the platform.

## Deploy

1. Validate the config and manifests.

```bash
go test ./...
go vet ./...
$(go env GOPATH)/bin/staticcheck ./...
$(go env GOPATH)/bin/kubeconform -strict -summary deploy/kubernetes/manifests.yaml
```

2. Apply the Kubernetes manifests.

```bash
kubectl apply -f deploy/kubernetes/manifests.yaml
kubectl rollout status deploy/api-server -n aihypervisor
```

3. Confirm the API is reachable.

```bash
kubectl -n aihypervisor port-forward svc/api-server 8080:80
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/ready
```

## Observe

- Scrape `/metrics` on the API metrics port exposed in `deploy/kubernetes/manifests.yaml`.
- Check `/health`, `/ready`, and `/live` for control-plane status.
- Use the Grafana dashboards in `deploy/grafana/dashboards/` to review cluster, GPU, and VM lifecycle trends.
- Review structured logs for `request_id` and `trace_id` fields when tracing user-visible failures.

Example log and metrics checks:

```bash
kubectl logs -n aihypervisor deploy/api-server -f
kubectl top pods -n aihypervisor
```

## Recover

1. If the API server is unhealthy, inspect its pod events and logs.

```bash
kubectl describe pod -n aihypervisor -l app=api-server
kubectl logs -n aihypervisor -l app=api-server --tail=200
```

2. Restart the deployment after fixing the underlying dependency.

```bash
kubectl rollout restart deployment/api-server -n aihypervisor
kubectl rollout status deployment/api-server -n aihypervisor
```

3. If the issue is external storage or messaging, verify Postgres, Redis, and NATS endpoints before recycling pods.

4. If only the scheduler or GPU orchestration path is impacted, isolate the failing component and check the corresponding service logs before rebalancing workloads.

5. Escalate to a full redeploy only after health checks stabilize and the metrics path is reporting again.

## References

- [Observability guide](observability.md)
- [Deployment guide](../deployment/DEPLOYMENT_GUIDE.md)
