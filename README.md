# kubectl-audit

[`kubectl-audit`](https://github.com/codenio/kubectl-audit) is a [`kubectl` plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that lists Kubernetes resources failing common health checks: unhealthy or crash-prone pods (including high container restarts), unhealthy nodes, unbound volumes, failed jobs, and suspended cron jobs. Output uses the same printers as `kubectl get` (default table, `-o wide`, JSON, YAML, custom columns, Go templates, and more).

## Contents

- [Install](#install)
- [Usage](#usage)
- [Resources and filters](#resources-and-filters)
- [Output formats](#output-formats)
- [Examples](#examples)
  - [Sample output](#sample-output)
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
| `pods`      | `pod`, `po` | Pods that need attention: phase is not `Running`, any regular container is not `Ready`, **or** any regular or init container has `RestartCount` ≥ **5** (threshold is fixed in code). `Succeeded` / `Completed` job pods are included because they are not in a running steady state. |
| `nodes`     | `node`, `no` | Nodes that are `NotReady` or have `SchedulingDisabled`. |
| `pvc`       | `pvcs`, `persistentvolumeclaim`, `persistentvolumeclaims` | PVCs not in `Bound` phase. |
| `pv`        | `pvs`, `persistentvolume`, `persistentvolumes` | PVs not in `Bound` phase. |
| `jobs`      | `job` | Failed jobs (including backoff / deadline failures). |
| `cronjobs`  | `cronjob`, `cj` | Suspended cron jobs. |

**Common flags**

- **All namespaces:** `-A` or `--all-namespaces` for namespaced resources (`pods`, `pvc`, `jobs`, `cronjobs`).
- **Labels:** `-l` / `--selector` (same semantics as `kubectl get`).

There are no separate pod subcommands or `--pending` / `--failed` style switches: one `kubectl audit pods` run applies all pod rules above.

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

```

### Sample output

For default and wide output, a **summary** is printed first, then the same table layout as `kubectl get` (including the **RESTARTS** column for pods).

- **total** — resources in scope (full list before the audit filter).
- **benign** — resources that pass the audit’s “OK” bar for that kind (for **pods**: `Running`, every regular container `Ready`, and no regular or init container has `RestartCount` of 5 or more).
- **attention** — rows in the filtered result (same as the table below the summary).

If the table is empty, **stderr** states either that nothing of that kind exists in the namespace (or cluster scope), or that **no resources require attention** (resources exist but none matched the audit). **Stdout** still shows the summary line in that case.

Names and namespaces below are **illustrative and masked**.

**Pods** (`kubectl audit pods -A`)

```text
Pod Audit summary: total = 2000 benign = 1996 attention = 4

NAMESPACE     NAME                                   READY   STATUS             RESTARTS   AGE
app-demo      web-frontend-7fb2c9-0                   0/1     Pending            0          2m
app-demo      api-rolling-xyz                         1/1     Running            12         1d
batch-demo    nightly-sync-abc12                      0/1     Completed          0          88m
data-demo     indexer-sidecar-xyz                     0/1     Completed          0          17m
```

*(The Running row is listed because restarts reached the attention threshold; Pending and Completed rows are non-steady-state phases.)*

**Nodes** (`kubectl audit nodes`)

```text
Node Audit summary: total = 53 benign = 51 attention = 2

NAME               STATUS     ROLES    AGE
node-worker-b02    NotReady   worker   40d
node-worker-d04    NotReady   worker   30d
```

*(Listed nodes are `NotReady` or cordoned `SchedulingDisabled`; benign counts Ready, schedulable nodes.)*

**Persistent volumes** (`kubectl audit pv`)

```text
PersistentVolume Audit summary: total = 120 benign = 119 attention = 1

NAME              CAPACITY   ACCESS MODES   STATUS      CLAIM
pv-archive-001    500Gi      RWO            Released    demo-ns/pvc-old-claim
```

**Persistent volume claims** (`kubectl audit pvc -A`)

```text
PersistentVolumeClaim Audit summary: total = 45 benign = 43 attention = 2

NAMESPACE    NAME              STATUS    VOLUME
app-demo     logs-claim-01     Pending
data-demo    backup-claim-02   Lost
```

**Jobs** (`kubectl audit jobs -A`)

```text
Job Audit summary: total = 28 benign = 26 attention = 2

NAMESPACE    NAME               COMPLETIONS   DURATION   AGE
batch-demo   daily-import       0/1           5m         5m
batch-demo   retry-migrate      0/1           1h         1h
```

**CronJobs** (`kubectl audit cronjobs -A`)

```text
CronJob Audit summary: total = 15 benign = 13 attention = 2

NAMESPACE   NAME               SCHEDULE      SUSPEND   ACTIVE
ops-demo    pause-backup       0 2 * * *     True      0
ops-demo    hold-reports       15 * * * *    True      0
```

For `-o json`, `-o yaml`, and other machine-oriented formats, the summary line is written to **stderr** so you can pipe **stdout** to `jq` or other tools unchanged.

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
