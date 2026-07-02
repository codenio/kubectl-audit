# Example: `kubectl audit pods`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-pods-demo`**, label **`demo=kubectl-audit-pods`**.

| Workload | What it exercises |
| -------- | ----------------- |
| **`kubectl-audit-demo-healthy`** | Running, all containers ready (benign). |
| **`kubectl-audit-demo-badpod`** | Invalid image → not ready (`ImagePullBackOff` / `ErrImagePull`). |
| **`kubectl-audit-demo-crashloop`** | Container exits immediately → `CrashLoopBackOff`. |
| **`kubectl-audit-demo-completed`** | Finished Job pod in `Succeeded` phase (flagged because phase ≠ `Running`). |

## Prerequisites (kind + podman)

Healthy and crash-loop workloads use `registry.k8s.io/pause:3.9` and `busybox:1.36`. On a local **kind** cluster with **podman**, pre-load them once:

```bash
podman pull --tls-verify=false registry.k8s.io/pause:3.9
podman pull --tls-verify=false docker.io/library/busybox:1.36
podman save registry.k8s.io/pause:3.9 -o /tmp/pause-3.9.tar
podman save docker.io/library/busybox:1.36 -o /tmp/busybox-1.36.tar
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive /tmp/pause-3.9.tar --name kind
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive /tmp/busybox-1.36.tar --name kind
rm -f /tmp/pause-3.9.tar /tmp/busybox-1.36.tar
```

## Try it

```bash
kubectl apply -f examples/audit-pods/demo.yaml
kubectl wait --for=condition=complete job/kubectl-audit-demo-completed -n kubectl-audit-pods-demo --timeout=120s
kubectl wait --for=condition=available deployment/kubectl-audit-demo-healthy -n kubectl-audit-pods-demo --timeout=120s
# Allow ~30s for image-pull and crash-loop failures to surface
kubectl audit pods -n kubectl-audit-pods-demo -l demo=kubectl-audit-pods
kubectl audit pods -n kubectl-audit-pods-demo -l demo=kubectl-audit-pods -o wide
kubectl delete -f examples/audit-pods/demo.yaml
```

For per-container detail (init failures, sidecars), use [`../audit-containers/`](../audit-containers/README.md). For Warning events, use [`../audit-events/`](../audit-events/README.md).
