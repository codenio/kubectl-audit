---
title: Find Unbound Kubernetes PV and PVC — kubectl-audit
description: Use kubectl audit pvc and kubectl audit pv to find PersistentVolumeClaims and PersistentVolumes not in Bound phase.
---

# PV & PVC

## PersistentVolumeClaims (PVC)

`kubectl audit pvc` lists PVCs not in `Bound` phase. Pods that reference an unbound PVC will remain `Pending`.

**Flagged phases:** `Pending` (not yet provisioned), `Lost` (backing PV lost or deleted)

### Demo

```bash
kubectl apply -f examples/audit-pvc/demo.yaml
kubectl audit pvc -n kubectl-audit-pvc-demo
kubectl delete -f examples/audit-pvc/demo.yaml
```

### Sample output

```
$ kubectl audit pvc -A
-------------------------------------------------------
PersistentVolumeClaim Audit summary: total = 45 benign = 43 attention = 2
-------------------------------------------------------
PersistentVolumeClaim that requires attention
NAMESPACE    NAME              STATUS    VOLUME
app-demo     logs-claim-01     Pending
data-demo    backup-claim-02   Lost
```

## PersistentVolumes (PV)

`kubectl audit pv` lists PVs not in `Bound` phase. PVs are cluster-scoped; `-n` and `-A` do not apply.

**Flagged phases:** `Available` (no claim bound), `Released` (claim deleted, awaiting reclaim), `Failed` (reclaim failure)

### Demo

```bash
kubectl apply -f examples/audit-pv/demo.yaml
kubectl audit pv
kubectl delete -f examples/audit-pv/demo.yaml
```

### Sample output

```
$ kubectl audit pv
-------------------------------------------------------
PersistentVolume Audit summary: total = 12 benign = 11 attention = 1
-------------------------------------------------------
PersistentVolume that requires attention
NAME              CAPACITY   ACCESS MODES   STATUS      CLAIM
pv-archive-001    500Gi      RWO            Released    demo-ns/pvc-old-claim
```

## Useful commands

```bash
kubectl audit pvc -A                 # PVCs across all namespaces
kubectl audit pv                     # PVs (cluster-scoped)
kubectl audit pvc -o json | jq ...   # JSON for scripting
```

## Next steps

- See detection rules → [Audit Rules: PVC](../reference/audit-rules.md#persistentvolumeclaims-pvc) / [PV](../reference/audit-rules.md#persistentvolumes-pv)
- A `Pending` PVC means Pods using it will be stuck in `Pending` too → check with `kubectl audit pods -A`
