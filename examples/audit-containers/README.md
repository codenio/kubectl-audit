# Example: `kubectl audit containers`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-containers-demo`**, label **`demo=kubectl-audit-containers`**.

| Workload | What it exercises |
| -------- | ----------------- |
| **`kubectl-audit-demo-healthy`** | Single `pause` container (benign row). |
| **`kubectl-audit-demo-bad-init`** | Init container exits `1` → **init-container** row in `Error`; app row stays `PodInitializing`. |
| **`kubectl-audit-demo-bad-sidecar`** | Main container healthy; sidecar has bad image → only the **sidecar** row flagged. |
| **`kubectl-audit-demo-crashloop`** | App container crash loop → **container** row in `CrashLoopBackOff`. |

Pod-level checks (phase, aggregate readiness) are in [audit-pods](../audit-pods/README.md).

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
kubectl apply -f examples/audit-containers/demo.yaml
kubectl wait --for=condition=available deployment/kubectl-audit-demo-healthy -n kubectl-audit-containers-demo --timeout=120s
# Allow a few seconds for init, sidecar, and crash-loop failures
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers -o wide
kubectl delete -f examples/audit-containers/demo.yaml
```
