---
title: Examples — kubectl-audit
description: Demo manifests and sample output for every kubectl-audit subcommand.
---

# Examples

Each subcommand has a demo manifest under [`examples/`](https://github.com/codenio/kubectl-audit/tree/main/examples) that creates intentionally broken resources — so you can see the plugin in action without needing a broken cluster.

## All demos

| Example | Subcommand | Demo folder |
| ------- | ---------- | ----------- |
| [Pods](pods.md) | `kubectl audit pods` | `examples/audit-pods/` |
| [Containers](containers.md) | `kubectl audit containers` | `examples/audit-containers/` |
| [Jobs & CronJobs](jobs.md) | `kubectl audit jobs` / `cronjobs` | `examples/audit-job/`, `audit-cronjob/` |
| [Services](services.md) | `kubectl audit service` | `examples/audit-svc/` |
| [Deployments](deployments.md) | `kubectl audit deploy` | `examples/audit-deploy/` |
| [PV & PVC](storage.md) | `kubectl audit pv` / `pvc` | `examples/audit-pv/`, `audit-pvc/` |
| [Events](events.md) | `kubectl audit events` | `examples/audit-events/` |
| [Nodes](nodes.md) | `kubectl audit nodes` | *(cluster-specific, no demo)* |
| [Namespaces](namespaces.md) | `kubectl audit ns` | `examples/audit-namespace/` |

## Demo pattern

Every demo follows this three-step pattern:

```bash
# 1. Create the broken resources
kubectl apply -f examples/audit-<target>/demo.yaml

# 2. Run the audit
kubectl audit <target> -n kubectl-audit-<target>-demo

# 3. Clean up
kubectl delete -f examples/audit-<target>/demo.yaml
```

## Reading the output

Every audit prints a summary to **stderr**:
```
-------------------------------------------------------
Pod Audit summary: total = 9 benign = 3 attention = 6
-------------------------------------------------------
```

- **total** — all resources of that type found in scope
- **benign** — resources passing all checks
- **attention** — resources with at least one failing check

The table that follows lists only the **attention** resources. If `attention = 0`, there is nothing to report and no table is printed.
