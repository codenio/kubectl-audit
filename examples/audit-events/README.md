# Example: `kubectl audit events`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-events-demo`**, label **`demo=kubectl-audit-events`**.

The audit lists every event with **`type: Warning`**. Normal events from the healthy Deployment are counted as benign but not shown.

| Workload | Typical Warning reasons |
| -------- | ----------------------- |
| **`kubectl-audit-demo-healthy`** | Normal events only (e.g. `Scheduled`, `Pulled`, `Started`). |
| **`kubectl-audit-demo-badpod`** | `Failed`, `ErrImagePull`, `ImagePullBackOff`. |
| **`kubectl-audit-demo-crashloop`** | `BackOff`, `Failed` (container restart). |
| **`kubectl-audit-demo-bad-init`** | `Failed` from the init container retry loop. |

## Prerequisites (kind + podman)

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
kubectl apply -f examples/audit-events/demo.yaml
kubectl wait --for=condition=available deployment/kubectl-audit-demo-healthy -n kubectl-audit-events-demo --timeout=120s
sleep 15
kubectl audit events -n kubectl-audit-events-demo
kubectl audit events -n kubectl-audit-events-demo -o wide
kubectl delete -f examples/audit-events/demo.yaml
```
