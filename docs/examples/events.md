---
title: Find Kubernetes Warning Events — kubectl-audit
description: Use kubectl audit events to surface Warning events across your cluster for quick triage.
---

# Events

`kubectl audit events` lists all events with `type: Warning`. Normal events are excluded — only Warning events signal potential problems.

## Demo

```bash
kubectl apply -f examples/audit-events/demo.yaml
kubectl audit events -n kubectl-audit-events-demo
kubectl delete -f examples/audit-events/demo.yaml
```

## Sample output

```
$ kubectl audit events -A
-------------------------------------------------------
Event Audit summary: total = 156 benign = 143 attention = 13
-------------------------------------------------------
Event that requires attention
NAMESPACE   LAST SEEN   TYPE      REASON              OBJECT                          MESSAGE
app-ns      5m          Warning   BackOff             Pod/web-dep-abc-xyz             Back-off restarting failed container
ops-ns      2m          Warning   FailedScheduling    Pod/batch-job-pod               0/3 nodes available: Insufficient cpu
staging     12m         Warning   Unhealthy           Pod/api-dep-rst-uvw             Liveness probe failed: ...
```

## Common warning reasons

| Reason | What it means |
| ------ | ------------- |
| `BackOff` | Container is crash-looping |
| `FailedScheduling` | No node can satisfy the Pod's resource/affinity requirements |
| `Unhealthy` | Liveness or readiness probe failing |
| `OOMKilled` | Container exceeded memory limit |
| `FailedMount` | Volume couldn't be mounted |
| `ImagePullBackOff` | Image pull failed after retries |

## Useful commands

```bash
kubectl audit events -A               # cluster-wide warning events
kubectl audit ev -n my-namespace      # namespace-scoped
kubectl audit events -o wide          # includes source, count columns
kubectl audit events -o json | jq -r '.items[] | "\(.reason): \(.message)"'
```

## Next steps

- `BackOff` events → investigate with [kubectl audit pods](pods.md) and [kubectl audit containers](containers.md)
- `FailedScheduling` → check [kubectl audit nodes](nodes.md) for unschedulable nodes
- See detection rules → [Audit Rules: Events](../reference/audit-rules.md#events)
