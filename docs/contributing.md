---
title: Contributing to kubectl-audit
description: How to contribute to kubectl-audit — development setup, testing, and PR guidelines.
---

# Contributing

Issues and pull requests are welcome at [github.com/codenio/kubectl-audit/issues](https://github.com/codenio/kubectl-audit/issues).

## Development setup

**Requirements:** Go 1.21+, `make`, a working Kubernetes cluster context.

```bash
git clone https://github.com/codenio/kubectl-audit.git
cd kubectl-audit

make bin          # build bin/audit
make test         # run tests + generate coverage profile
make fmt          # format code
make vet          # vet code
make install      # install to ~/.krew/bin/kubectl-audit
```

Run without installing:

```bash
go run ./cmd/plugin --help
go run ./cmd/plugin pods
go run ./cmd/plugin containers --help
```

## Before opening a PR

- [ ] `make test` passes
- [ ] `make fmt` and `make vet` are clean (`make bin` runs both)
- [ ] Validated against a real cluster with `make install && kubectl audit …`

## Repository layout

| Path | Purpose |
| ---- | ------- |
| `cmd/plugin/` | Entrypoint and CLI (Cobra, config flags, printing) |
| `pkg/plugin/` | Audit logic, container list/table data, server-side table handling |
| `deploy/krew/plugin.yaml` | Krew manifest template for releases |
| `examples/audit-*/` | Demo manifests per audit type |
| `docs/` | MkDocs documentation source |

## Releasing

Releases are automated via GitHub Actions on `v*.*.*` tags. See `.github/workflows/release.yml`.

## Bumping Kubernetes dependencies

```bash
make kubernetes-deps
```

## Reporting bugs

Please include:

- Plugin version (binary or release tag)
- `kubectl version` output
- Exact command and full error message
- Kubernetes distribution and version (EKS, GKE, k3s, etc.)

Open an issue: [github.com/codenio/kubectl-audit/issues/new](https://github.com/codenio/kubectl-audit/issues/new)
