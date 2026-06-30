---
title: RBAC & Permissions
description: Minimum Kubernetes RBAC permissions required to run kubectl-audit.
---

# RBAC & Permissions

`kubectl-audit` is a **read-only** plugin — it never creates, patches, or deletes resources. It uses `list` and `get` on the following API groups.

## Minimum ClusterRole

Apply this to grant the minimum permissions needed to run all audit subcommands:

```yaml title="kubectl-audit-reader.yaml"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubectl-audit-reader
rules:
  - apiGroups: [""]
    resources:
      - pods
      - nodes
      - persistentvolumes
      - persistentvolumeclaims
      - services
      - events
    verbs: ["get", "list"]
  - apiGroups: ["apps"]
    resources:
      - deployments
    verbs: ["get", "list"]
  - apiGroups: ["batch"]
    resources:
      - jobs
      - cronjobs
    verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubectl-audit-reader-binding
subjects:
  - kind: User
    name: <your-username>     # replace with your user or ServiceAccount
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: kubectl-audit-reader
  apiGroup: rbac.authorization.k8s.io
```

```bash
kubectl apply -f kubectl-audit-reader.yaml
```

## Per-subcommand permissions

If you want to grant access only to specific audits:

| Subcommand | API group | Resource |
| ---------- | --------- | -------- |
| `pods` | `""` (core) | `pods` |
| `containers` | `""` (core) | `pods` |
| `nodes` | `""` (core) | `nodes` |
| `pvc` | `""` (core) | `persistentvolumeclaims` |
| `pv` | `""` (core) | `persistentvolumes` |
| `service` | `""` (core) | `services`, `pods` |
| `events` | `""` (core) | `events` |
| `deploy` | `apps` | `deployments` |
| `jobs` | `batch` | `jobs` |
| `cronjobs` | `batch` | `cronjobs` |

!!! note "`service` requires pods"
    The `service` audit needs to list Pods to check if any match each Service's selector. Granting only `services` list will cause the audit to fail.

## Checking your current permissions

```bash
kubectl auth can-i list pods --all-namespaces
kubectl auth can-i list nodes
kubectl auth can-i list deployments --all-namespaces
```

If `kubectl audit` returns a permission error, check access for the specific resource the subcommand queries.
