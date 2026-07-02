---
title: Quick Start — kubectl-audit
description: Run your first Kubernetes health audit in under a minute.
---

# Quick Start

After [installing](install.md), run your first audit:

```bash
kubectl audit pods
```

![kubectl audit pods demo](../assets/audit-pods.gif){ .demo-gif }

## Choose the right subcommand

| I want to… | Command |
| ---------- | ------- |
| See all unhealthy pods | `kubectl audit pods -A` |
| Drill into a specific pod's containers | `kubectl audit containers -p <pod-name>` |
| Check which nodes can't schedule work | `kubectl audit nodes` |
| Find services with no backing pods | `kubectl audit svc -A` |
| Check deployment rollout health | `kubectl audit deploy -A` |
| Catch warning events cluster-wide | `kubectl audit events -A` |
| Check storage (PVC/PV) issues | `kubectl audit pvc -A && kubectl audit pv` |

!!! tip "pods vs containers"
    - `kubectl audit pods` — one row per pod, aggregate status
    - `kubectl audit containers` — one row per container, shows exact image/restart detail
    
    Start with `pods` for a quick count. Use `containers` when you need to identify *which* container inside a multi-container pod is broken.

## Common workflows

### Incident triage — cluster-wide sweep

```bash
kubectl audit pods -A
kubectl audit nodes
kubectl audit events -A
```

### Pre-deploy smoke check

```bash
kubectl audit pods -n my-namespace
kubectl audit deploy -n my-namespace
kubectl audit svc -n my-namespace
```

### Scope to a team namespace with label filter

```bash
kubectl audit pods -n team-a -l app=api
kubectl audit containers -n team-a -l app=api
```

## Try a demo (no broken cluster needed)

Each audit type has a demo manifest that creates intentionally broken resources:

=== "Pods"
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-pods/demo.yaml
    kubectl audit pods -n kubectl-audit-pods-demo
    kubectl delete -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-pods/demo.yaml
    ```

=== "Services"
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-svc/demo.yaml
    kubectl audit service -n kubectl-audit-svc-demo
    kubectl delete -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-svc/demo.yaml
    ```

=== "Events"
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-events/demo.yaml
    kubectl audit events -n kubectl-audit-events-demo
    kubectl delete -f https://raw.githubusercontent.com/codenio/kubectl-audit/main/examples/audit-events/demo.yaml
    ```

See [all demo examples →](../examples/index.md)

## Scripting and CI/CD

The audit summary is written to **stderr**, so stdout is clean for piping:

```bash
# Extract names of failing pods as JSON
kubectl audit pods -o json 2>/dev/null | jq '[.items[].metadata.name]'

# Non-zero exit if any unhealthy pods found
kubectl audit pods -n my-namespace -o name | grep -q . && echo "Issues found"
```

See the full [CI/CD integration guide →](../guides/cicd.md)

---

Next: [RBAC & Permissions →](rbac.md)
