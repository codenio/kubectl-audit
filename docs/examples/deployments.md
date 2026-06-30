---
title: Find Unhealthy Kubernetes Deployments — kubectl-audit
description: Use kubectl audit deploy to find deployments scaled to zero or with fewer ready replicas than desired.
---

# Deployments

`kubectl audit deploy` finds Deployments that are scaled to zero or have fewer ready replicas than desired.

## Demo

```bash
kubectl apply -f examples/audit-deploy/demo.yaml
kubectl audit deploy -n kubectl-audit-deploy-demo
kubectl delete -f examples/audit-deploy/demo.yaml
```

## Sample output

```
$ kubectl audit deploy -A
-------------------------------------------------------
Deployment Audit summary: total = 30 benign = 27 attention = 3
-------------------------------------------------------
Deployment that requires attention
NAMESPACE   NAME               READY   UP-TO-DATE   AVAILABLE   AGE
app-ns      web-frontend       0/3     3            0           2h
ops-ns      batch-worker       0/0     0            0           5d
staging     api-gateway        1/3     3            1           45m
```

## Detection rules

| Condition | Flagged |
| --------- | ------- |
| `spec.replicas == 0` | ✅ Explicitly scaled to zero |
| `spec.replicas` unset and `readyReplicas < 1` | ✅ API default is 1 |
| `readyReplicas < spec.replicas` when `spec.replicas > 0` | ✅ Rollout not complete or replicas crashing |
| `readyReplicas == spec.replicas > 0` | ❌ Healthy |

## Useful commands

```bash
kubectl audit deploy -A
kubectl audit deploy -n my-namespace -l app=web
kubectl audit deploy -o json | jq '.items[] | {name: .metadata.name, ready: .status.readyReplicas}'
```

## Next steps

- Deployment unhealthy? Check the pods → [kubectl audit pods](pods.md)
- See full rules → [Audit Rules: Deployments](../reference/audit-rules.md#deployments)
