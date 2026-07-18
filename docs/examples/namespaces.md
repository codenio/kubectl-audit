---
title: Find Terminating and Empty Kubernetes Namespaces — kubectl-audit
description: Use kubectl audit ns to find namespaces stuck in Terminating or with no workloads.
---

# Namespaces

`kubectl audit ns` lists namespaces that need attention — either stuck in **`Terminating`** phase, or **`Active`** but containing **no workloads**.

Workloads checked: Pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, and CronJobs. ConfigMaps, Secrets, and Services alone do **not** count.

Namespaces are cluster-scoped; `-n` and `-A` do not apply. Use `-l` / `--selector` to filter namespaces by label.

## Demo

```bash
kubectl apply -f examples/audit-namespace/demo.yaml
kubectl wait --for=condition=available deployment/healthy -n kubectl-audit-namespace-workload-demo --timeout=120s
kubectl delete namespace kubectl-audit-namespace-terminating-demo --wait=false
kubectl audit ns -l demo=kubectl-audit-namespace
kubectl patch namespace kubectl-audit-namespace-terminating-demo -p '{"metadata":{"finalizers":[]}}' --type=merge 2>/dev/null || true
kubectl delete -f examples/audit-namespace/demo.yaml --ignore-not-found
```

See [`examples/audit-namespace/README.md`](https://github.com/codenio/kubectl-audit/tree/main/examples/audit-namespace/README.md) for the full walkthrough (empty **Active** namespace, benign workload namespace, and stuck **Terminating** namespace).

## Sample output

```
$ kubectl audit ns -l demo=kubectl-audit-namespace
-------------------------------------------------------
Namespace Audit summary: total = 3 benign = 1 attention = 2
-------------------------------------------------------
Namespace that requires attention
NAME                                      STATUS        AGE
kubectl-audit-namespace-empty-demo        Active        2m
kubectl-audit-namespace-terminating-demo  Terminating   1m
```

**`Terminating`** — deletion is stuck; resources or finalizers may block removal.  
**`Active` (empty)** — no Pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, or CronJobs; often safe to delete after review.

## Useful commands

```bash
kubectl audit ns
kubectl audit namespace -l team=platform
kubectl audit ns -o json | jq '.items[] | {name: .metadata.name, phase: .status.phase}'
kubectl audit ns -o wide
```

## Next steps

- See detection rules → [Audit Rules: Namespaces](../reference/audit-rules.md#namespaces)
- Delete an unused namespace: `kubectl delete namespace <name>`
- Debug a stuck namespace: `kubectl get namespace <name> -o yaml`
