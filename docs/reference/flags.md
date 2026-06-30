# Flags Reference

`kubectl-audit` uses standard `kubectl` configuration for context, namespace, and kubeconfig. All global flags accepted by `kubectl` (e.g., `--kubeconfig`, `--context`, `-n`) apply.

## Global flags (all subcommands)

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `--context` | current context | Kubeconfig context to use |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |
| `-n`, `--namespace` | current namespace | Target namespace for namespaced resources |
| `-o`, `--output` | table | Output format (see [Output Formats](output-formats.md)) |
| `-v`, `--v` | 0 | Log verbosity level |

## Namespace flags

| Flag | Short | Description |
| ---- | ----- | ----------- |
| `--all-namespaces` | `-A` | List across all namespaces (namespaced subcommands only) |
| `--namespace` | `-n` | Target a specific namespace |

## Selector flag

| Flag | Short | Description |
| ---- | ----- | ----------- |
| `--selector` | `-l` | Label selector (e.g., `app=web`, `env in (prod,staging)`) |

For `containers`: applies to the Pod list. For `service`: applies to the Service list. For `deploy`: applies to the Deployment list.

## containers-specific flag

| Flag | Short | Description |
| ---- | ----- | ----------- |
| `--pod` | `-p` | Filter to an exact pod `metadata.name` |

## Getting help

```bash
kubectl audit --help
kubectl audit <subcommand> --help
```
