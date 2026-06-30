---
title: Installation — kubectl-audit
description: How to install the kubectl-audit plugin via Krew or from source.
---

# Installation

## Requirements

- **kubectl** 1.20+
- **Kubernetes cluster** 1.21+ (uses `batch/v1`, `apps/v1`, `core/v1`, `events.k8s.io/v1`)
- **[Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)** (recommended install method)

## Krew — recommended

`kubectl-audit` is published through the [codenio Krew index](https://github.com/codenio/codenio-krew-index).

```bash
kubectl krew index add codenio https://github.com/codenio/codenio-krew-index.git
kubectl krew install codenio/audit
```

!!! tip "Index name conflict"
    If the name `codenio` is already taken in your Krew config, choose any other name for the `index add` step:
    ```bash
    kubectl krew index add my-codenio https://github.com/codenio/codenio-krew-index.git
    kubectl krew install my-codenio/audit
    ```
    Upgrade stays the same regardless of index name:
    ```bash
    kubectl krew upgrade audit
    ```

## Krew default index

If the plugin is listed in Krew's default index:

```bash
kubectl krew install audit
```

## From source

Requires [Go](https://go.dev/dl/) 1.21+ and `make`.

```bash
git clone https://github.com/codenio/kubectl-audit.git
cd kubectl-audit
make install
```

Builds `bin/audit` and copies it to `~/.krew/bin/kubectl-audit`. Ensure `~/.krew/bin` is on your `PATH`.

## Verify

```bash
kubectl audit --help
```

You should see the plugin's help output listing all subcommands.

---

Next: [Quick Start →](quickstart.md)
