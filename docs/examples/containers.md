---
title: Per-Container Kubernetes Health Check — kubectl-audit
description: Use kubectl audit containers to find individual containers with image pull errors, crash loops, or high restarts.
---

# Containers

`kubectl audit containers` gives per-container visibility — one row per container (init or app), showing exactly which container inside a pod is broken.

Use this after `kubectl audit pods` shows a pod in a bad state but you need to know which specific container is the cause.

## Demo

```bash
kubectl apply -f examples/audit-containers/demo.yaml
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers
kubectl delete -f examples/audit-containers/demo.yaml
```

## Sample output

=== "Default"
    ```
    $ kubectl audit containers -n demo
    -------------------------------------------------------
    Container Audit summary: total = 12 benign = 9 attention = 3
    -------------------------------------------------------
    Container that requires attention
    POD                     NAME       READY   STATUS             RESTARTS   AGE   TYPE
    workload-dep-abc-xyz    sidecar    0       ImagePullBackOff   0          -     container
    workload-dep-abc-xyz    init-cfg   -       Failed             2          -     init-container
    ```

=== "Wide (-o wide)"
    ```
    $ kubectl audit containers -n demo -o wide
    POD                     NAME       READY   STATUS             RESTARTS   AGE   TYPE            PORTS   IMAGE                     PULLPOLICY
    workload-dep-abc-xyz    sidecar    0       ImagePullBackOff   0          -     container               myrepo/sidecar:bad-tag    IfNotPresent
    ```

## Scope to a single pod

```bash
kubectl audit containers -p my-pod-0
kubectl audit containers --pod my-pod-0 -n my-namespace
```

## Output formats

| Format | Description |
| ------ | ----------- |
| *(default)* | Fixed table: `POD`, `NAME`, `READY`, `STATUS`, `RESTARTS`, `AGE`, `TYPE` |
| `-o wide` | Adds `PORTS`, `IMAGE`, `PULLPOLICY` |
| `-o json` | Full JSON |
| `-o yaml` | Full YAML |
| `-o name` | `namespace/pod-name` lines |

!!! warning
    `custom-columns`, `jsonpath`, and Go templates are not supported for `containers`. Use `-o json | jq` instead.

## Next steps

- See detection rules → [Audit Rules: Containers](../reference/audit-rules.md#containers)
- Aggregate pod view → [kubectl audit pods](pods.md)
