---
title: Subcommands Reference
description: Full reference for all kubectl-audit subcommands, aliases, and flags.
---

# Subcommands Reference

## Subcommand overview

| Subcommand | Aliases | Scope |
| ---------- | ------- | ----- |
| `pods` | `pod`, `po` | Namespaced |
| `containers` | `container` | Namespaced |
| `nodes` | `node`, `no` | Cluster-wide |
| `pvc` | `pvcs`, `persistentvolumeclaim`, `persistentvolumeclaims` | Namespaced |
| `pv` | `pvs`, `persistentvolume`, `persistentvolumes` | Cluster-wide |
| `jobs` | `job` | Namespaced |
| `cronjobs` | `cronjob`, `cj` | Namespaced |
| `service` | `services`, `svc` | Namespaced |
| `deploy` | `deployment`, `deployments` | Namespaced |
| `events` | `event`, `ev` | Namespaced |

For exact detection rules per subcommand, see [Audit Rules](audit-rules.md).

---

## pods

Lists pods that are not in a healthy steady state.

```bash
kubectl audit pods [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** phase ≠ `Running`, any container not `Ready`, or restart count ≥ 5. `Succeeded`/`Completed` pods are included.

**Default columns:** `NAME`, `READY`, `STATUS`, `RESTARTS`, `AGE`  
**With `-A`:** prepends `NAMESPACE` column (same as `kubectl get pods -A`)

---

## containers

Per-container rows — finer detail than `pods`.

```bash
kubectl audit containers [-n NAMESPACE] [-A] [-l SELECTOR] [-p POD] [-o FORMAT]
```

**Flagged when:** image pull errors, crash-loop waiting states, not ready, restart count ≥ 5, failed init containers.

**Default columns:** `POD`, `NAME`, `READY`, `STATUS`, `RESTARTS`, `AGE`, `TYPE`  
**With `-o wide`:** adds `PORTS`, `IMAGE`, `PULLPOLICY`

!!! tip "pods vs containers"
    Use `pods` for a quick count. Use `containers` when you need to pinpoint which container inside a multi-container pod is broken.

---

## nodes

Lists nodes that cannot be scheduled to.

```bash
kubectl audit nodes [-o FORMAT]
```

**Flagged when:** `NotReady` or `SchedulingDisabled` (cordoned via `kubectl cordon`).  
Nodes are cluster-scoped — `-n` and `-A` do not apply.

---

## pvc

Lists PersistentVolumeClaims not in `Bound` phase.

```bash
kubectl audit pvc [-n NAMESPACE] [-A] [-o FORMAT]
```

**Flagged phases:** `Pending`, `Lost`

---

## pv

Lists PersistentVolumes not in `Bound` phase.

```bash
kubectl audit pv [-o FORMAT]
```

**Flagged phases:** `Available`, `Released`, `Failed`  
PVs are cluster-scoped — `-n` and `-A` do not apply.

---

## jobs

Lists failed Jobs.

```bash
kubectl audit jobs [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** `status.failed > 0` (includes backoff limit reached, active deadline exceeded).

---

## cronjobs

Lists suspended CronJobs.

```bash
kubectl audit cronjobs [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** `spec.suspend = true`

---

## service

Lists Services whose selector matches no Pods.

```bash
kubectl audit service [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** non-empty `spec.selector` with no matching Pods in the namespace.  
**Skipped:** `type: ExternalName`, empty selectors.

!!! note
    `-l`/`--selector` filters the **Service list**. Pod matching always uses all Pods in the namespace.

---

## deploy

Lists Deployments that are scaled down or have unhealthy replicas.

```bash
kubectl audit deploy [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** `spec.replicas == 0`, or `readyReplicas < desired` when `desired > 0`.  
`desired` defaults to **1** when `spec.replicas` is unset (matches Kubernetes API default).

---

## events

Lists Warning events.

```bash
kubectl audit events [-n NAMESPACE] [-A] [-l SELECTOR] [-o FORMAT]
```

**Flagged when:** `type = Warning`. Normal events are excluded.

---

## Common flags

| Flag | Short | Applies to | Description |
| ---- | ----- | ---------- | ----------- |
| `--all-namespaces` | `-A` | Namespaced targets | List across all namespaces |
| `--namespace` | `-n` | Namespaced targets | Target namespace |
| `--selector` | `-l` | All | Label selector (e.g., `app=web`) |
| `--output` | `-o` | All | Output format (see [Output Formats](output-formats.md)) |
| `--context` | | All | Kubeconfig context override |
| `--kubeconfig` | | All | Path to kubeconfig file |
| `--pod` | `-p` | `containers` only | Filter to a single pod by exact name |

---

## Getting help

```bash
kubectl audit --help
kubectl audit <subcommand> --help
```
