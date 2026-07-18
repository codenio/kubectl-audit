---
title: Find Terminating and Empty Kubernetes Namespaces — kubectl-audit
description: Use kubectl audit ns to find namespaces stuck in Terminating or with no workloads.
---

# Namespaces

`kubectl audit ns` lists namespaces that need attention — either stuck in **`Terminating`** phase, or **`Active`** but containing **no workloads**.

Workloads checked: Pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, and CronJobs. ConfigMaps, Secrets, and Services alone do **not** count.

Namespaces are cluster-scoped; `-n` and `-A` do not apply. Use `-l` / `--selector` to filter namespaces by label.

!!! info
    There is no demo manifest for namespaces — empty or stuck namespaces are cluster-specific. Run against your own cluster.

## Sample output

```
$ kubectl audit ns
-------------------------------------------------------
Namespace Audit summary: total = 42 benign = 38 attention = 4
-------------------------------------------------------
Namespace that requires attention
NAME              STATUS        AGE
kube-public       Active        400d
old-team-a        Active        90d
stuck-delete      Terminating   2d
unused-staging    Active        14d
```

**`Terminating`** — deletion is stuck; resources may remain and block removal.  
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
