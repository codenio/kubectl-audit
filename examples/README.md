# Example manifests

Each folder is self-contained: apply **`demo.yaml`**, run the matching **`kubectl audit …`** command, then delete the same file.

## Local kind cluster (podman)

Several demos use `registry.k8s.io/pause:3.9` and `busybox:1.36`. On **kind** with **podman**, pre-load them once (TLS verify is often required off for local pulls):

```bash
podman pull --tls-verify=false registry.k8s.io/pause:3.9
podman pull --tls-verify=false docker.io/library/busybox:1.36
podman save registry.k8s.io/pause:3.9 -o /tmp/pause-3.9.tar
podman save docker.io/library/busybox:1.36 -o /tmp/busybox-1.36.tar
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive /tmp/pause-3.9.tar --name kind
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive /tmp/busybox-1.36.tar --name kind
rm -f /tmp/pause-3.9.tar /tmp/busybox-1.36.tar
```

| Folder | Command |
| ------ | ------- |
| [audit-pods](audit-pods/README.md) | `kubectl audit pods` |
| [audit-containers](audit-containers/README.md) | `kubectl audit containers` |
| [audit-job](audit-job/README.md) | `kubectl audit jobs` |
| [audit-cronjob](audit-cronjob/README.md) | `kubectl audit cronjobs` |
| [audit-svc](audit-svc/README.md) | `kubectl audit service` |
| [audit-deploy](audit-deploy/README.md) | `kubectl audit deploy` |
| [audit-events](audit-events/README.md) | `kubectl audit events` |
| [audit-pvc](audit-pvc/README.md) | `kubectl audit pvc` |
| [audit-pv](audit-pv/README.md) | `kubectl audit pv` |

There is **no** `audit-nodes` or `audit-namespace` sample (not portable across clusters). See the main [README](../README.md) for `kubectl audit nodes` and `kubectl audit ns`.
