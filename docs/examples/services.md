---
title: Find Kubernetes Services with No Backing Pods — kubectl-audit
description: Use kubectl audit service to detect services whose pod selector matches no pods, causing silent traffic drops.
---

# Services

`kubectl audit service` finds Services whose pod selector matches no Pods in the namespace — these Services will silently drop all traffic.

**Skipped:** `type: ExternalName` (DNS-based routing, no selector) and Services with an empty selector (headless / manually managed endpoints).

## Demo

```bash
kubectl apply -f examples/audit-svc/demo.yaml
kubectl audit service -n kubectl-audit-svc-demo
kubectl delete -f examples/audit-svc/demo.yaml
```

## Sample output

```
$ kubectl audit svc -A
-------------------------------------------------------
Service Audit summary: total = 24 benign = 22 attention = 2
-------------------------------------------------------
Service that requires attention
NAMESPACE   NAME              TYPE        CLUSTER-IP     PORT(S)    AGE
app-ns      broken-svc        ClusterIP   10.96.0.100    80/TCP     2h
web-ns      orphan-frontend   ClusterIP   10.96.0.200    443/TCP    45m
```

## Common causes

- Selector label mismatch — Service selects `app: frontend` but Pods are labeled `app: web-frontend`
- All matching Pods were deleted or evicted
- Namespace mismatch — Pods exist in a different namespace
- Wrong namespace on the Service — Service is in `staging`, Pods are in `production`

## Useful commands

```bash
kubectl audit svc -A
kubectl audit service -n my-namespace -l app=myapp
kubectl audit service -o json | jq -r '.items[] | "\(.metadata.namespace)/\(.metadata.name)"'
```

!!! note
    The `-l`/`--selector` flag filters **Services**, not Pods. Pod matching always uses all Pods in the namespace.

## Next steps

- Found a broken service? Check the pods it should be selecting:
  ```bash
  kubectl audit pods -n <namespace> -l <service-selector-labels>
  ```
- See full rules → [Audit Rules: Services](../reference/audit-rules.md#services)
