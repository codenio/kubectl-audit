---
title: Output Formats — kubectl-audit
description: All supported -o output formats for kubectl-audit subcommands.
---

# Output Formats

## Standard subcommands

For `pods`, `nodes`, `pv`, `pvc`, `jobs`, `cronjobs`, `service`, `deploy`, and `events`, the `-o` flag works identically to `kubectl get`:

```bash
kubectl audit pods -o wide
kubectl audit nodes -o json
kubectl audit pvc -o yaml
kubectl audit jobs -o custom-columns=NAME:.metadata.name
kubectl audit service -o json
kubectl audit events -o wide
kubectl audit deploy -o jsonpath='{.items[*].metadata.name}'
```

## containers subcommand

The `containers` subcommand uses a dedicated printer with limited format support:

| Format | Output |
| ------ | ------ |
| *(default)* | Table: `POD`, `NAME`, `READY`, `STATUS`, `RESTARTS`, `AGE`, `TYPE` |
| `-o wide` | Adds `PORTS`, `IMAGE`, `PULLPOLICY` |
| `-o json` | Full JSON |
| `-o yaml` | Full YAML |
| `-o name` | `namespace/pod-name` lines |

!!! warning "Unsupported formats for `containers`"
    `custom-columns`, `jsonpath`, and Go templates are **not supported** and will return a clear error. Use `-o json | jq` instead:
    ```bash
    kubectl audit containers -o json | jq '.items[] | {pod: .pod, name: .name, status: .status}'
    ```

## Audit summary line

Every audit writes a summary to **stderr**, keeping stdout clean for piping:

```
Pod Audit summary: total = 120 benign = 114 attention = 6
```

Redirect as needed:

```bash
kubectl audit pods -o json 2>/dev/null | jq ...    # suppress summary
kubectl audit pods 2>&1 | tee audit.log             # capture both
```

## Custom columns example

```bash
kubectl audit deploy -o custom-columns=\
NAME:.metadata.name,\
NAMESPACE:.metadata.namespace,\
DESIRED:.spec.replicas,\
READY:.status.readyReplicas
```

## Machine-readable scripting

```bash
# Count unhealthy pods
kubectl audit pods -A -o json 2>/dev/null | jq '.items | length'

# Names of flagged services
kubectl audit service -A -o json 2>/dev/null \
  | jq -r '.items[] | "\(.metadata.namespace)/\(.metadata.name)"'
```

See the [CI/CD integration guide](../guides/cicd.md) for complete scripting patterns.
