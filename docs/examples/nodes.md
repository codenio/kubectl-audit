---
title: Find NotReady and Cordoned Kubernetes Nodes — kubectl-audit
description: Use kubectl audit nodes to find nodes that are NotReady or SchedulingDisabled (cordoned).
---

# Nodes

`kubectl audit nodes` lists nodes that cannot be scheduled to — either because they are `NotReady` or because they have been cordoned (`SchedulingDisabled`).

Nodes are cluster-scoped; `-n` and `-A` do not apply.

!!! info
    There is no demo manifest for nodes — node states are cluster-specific. Run against your own cluster.

## Sample output

```
$ kubectl audit nodes
-------------------------------------------------------
Node Audit summary: total = 66 benign = 46 attention = 20
-------------------------------------------------------
Node that requires attention
NAME            STATUS                        ROLES    AGE     VERSION
default-0       Ready,SchedulingDisabled      <none>   459d    v1.30.3
default-1       Ready,SchedulingDisabled      <none>   35d     v1.30.3
worker-0        NotReady,SchedulingDisabled   <none>   2d21h   v1.33.3
worker-1        NotReady,SchedulingDisabled   <none>   23h     v1.33.3
worker-3        NotReady                      <none>   24h     v1.33.3
```

**`Ready,SchedulingDisabled`** — node is healthy but cordoned (e.g., draining for maintenance). No new Pods will be scheduled here.  
**`NotReady`** — node is not reporting healthy to the control plane. Pods on this node may be evicted.  
**`NotReady,SchedulingDisabled`** — both: node is down and cordoned.

## Useful commands

```bash
kubectl audit nodes
kubectl audit nodes -o json | jq '.items[].metadata.name'
kubectl audit nodes -o wide
```

## Next steps

- See detection rules → [Audit Rules: Nodes](../reference/audit-rules.md#nodes)
- Cordoning a node: `kubectl cordon <node-name>`
- Draining a node: `kubectl drain <node-name> --ignore-daemonsets`
