# kubectl-audit

[`kubectl-audit`](https://github.com/codenio/kubectl-audit) is a [`kubectl` plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that lists Kubernetes resources failing common health checks: stuck pods, unhealthy nodes, unbound volumes, failed jobs, and suspended cron jobs. Output uses the same printers as `kubectl get` (default table, `-o wide`, JSON, YAML, custom columns, Go templates, and more).

## Contents

- [Install](#install)
- [Usage](#usage)
- [Resources and filters](#resources-and-filters)
- [Output formats](#output-formats)
- [Examples](#examples)
- [Development](#development)
- [Contributing](#contributing)
- [Acknowledgments](#acknowledgments)

## Install

### Krew (codenio custom index)

Published through [codenio-krew-index](https://github.com/codenio/codenio-krew-index):

```bash
kubectl krew index add codenio https://github.com/codenio/codenio-krew-index.git
kubectl krew install codenio/audit
```

Upgrade when the index is updated:

```bash
kubectl krew upgrade codenio/audit
```

If the local index name `codenio` is already taken, pick another name for `kubectl krew index add` and use that same prefix in `install` / `upgrade` (for example `codenio-krew/audit`). Maintainer notes for bumping versions live in the index repository.

### Krew (default index)

If the plugin is listed in Krew’s default index:

```bash
kubectl krew install audit
```

### From source

Requires [Go](https://go.dev/dl/) 1.21+ and `make`:

```bash
git clone https://github.com/codenio/kubectl-audit.git
cd kubectl-audit
make install
```

This builds `bin/audit` and copies it to `~/.krew/bin/kubectl-audit`. Put `~/.krew/bin` on your `PATH` (or adjust the install path) so `kubectl audit` resolves.

## Usage

Every invocation needs a **resource** argument (see [Resources and filters](#resources-and-filters)):

```bash
kubectl audit pods
kubectl audit --help
```

Standard `kubectl` config applies: current context, `KUBECONFIG`, `-n` / `--namespace`, `--context`, and so on.

## Resources and filters

| Resource    | Aliases | What is listed |
| ----------- | ------- | -------------- |
| `pods`      | `pod`, `po` | By default, pods whose phase is not `Running` (for example `Pending`, `Failed`, `Succeeded`). Use `--pending`, `--failed`, and `--not-ready` to narrow or extend what is shown. |
| `nodes`     | `node`, `no` | Nodes that are `NotReady` or have `SchedulingDisabled`. |
| `pvc`       | `pvcs`, `persistentvolumeclaim`, `persistentvolumeclaims` | PVCs not in `Bound` phase. |
| `pv`        | `pvs`, `persistentvolume`, `persistentvolumes` | PVs not in `Bound` phase. |
| `jobs`      | `job` | Failed jobs (including backoff / deadline failures). |
| `cronjobs`  | `cronjob`, `cj` | Suspended cron jobs. |

**Common flags**

- **All namespaces:** `-A` or `--all-namespaces` for namespaced resources (`pods`, `pvc`, `jobs`, `cronjobs`).
- **Labels:** `-l` / `--selector` (same semantics as `kubectl get`).
- **Pods only:** combine `--pending`, `--failed`, and `--not-ready` to widen which pods are included (`--help` describes each).

Further notes live in [doc/USAGE.md](doc/USAGE.md).

## Output formats

[`kubectl get`](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_get/)-style `-o` flags are supported, for example:

```bash
kubectl audit pods -o wide
kubectl audit nodes -o json
kubectl audit pvc -o yaml
kubectl audit jobs -o custom-columns=NAME:.metadata.name
```

## Examples

```bash
# Pods in the current namespace (default table)
kubectl audit pods

# Pods cluster-wide (namespace column when listing all namespaces)
kubectl audit po -A

# Specific namespace
kubectl audit pods -n tester

# Label filter
kubectl audit pods -A -l app=nginx

# PVC / PV
kubectl audit pvc -A
kubectl audit pv

# Jobs and CronJobs
kubectl audit job -A
kubectl audit cj -A

# Pod filters (combine as needed)
kubectl audit pods --pending --failed --not-ready
```

## Development

**Prerequisites:** Go 1.21+, `make`, and a working cluster context if you want to run the plugin end to end.

```bash
git clone https://github.com/codenio/kubectl-audit.git
cd kubectl-audit
make bin          # writes bin/audit
make test         # tests + coverage profile
make fmt && make vet
```

Run without installing (same flags as under `kubectl audit`):

```bash
go run ./cmd/plugin --help
go run ./cmd/plugin pods
```

**Repository layout**

- `cmd/plugin/` — entrypoint and CLI (`cobra`, config flags, printing).
- `pkg/plugin/` — audit logic and server-side table handling.
- `deploy/krew/plugin.yaml` — Krew manifest template for releases.

To bump pinned Kubernetes dependencies, use the `kubernetes-deps` target in the `Makefile`.

## Contributing

Issues and pull requests are welcome: [github.com/codenio/kubectl-audit/issues](https://github.com/codenio/kubectl-audit/issues).

Before you open a PR:

- `make test` passes.
- `make bin` succeeds and `make fmt` / `make vet` are clean (or run `make bin`, which runs `fmt` and `vet` first).

Validate changes against a cluster with `make install` and `kubectl audit …` as needed.

## Acknowledgments

This plugin was created using [replicatedhq/krew-plugin-template](https://github.com/replicatedhq/krew-plugin-template).
