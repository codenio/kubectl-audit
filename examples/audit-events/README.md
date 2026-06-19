# Example: `kubectl audit events`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-events-demo`**, label **`demo=kubectl-audit-events`**.

**`kubectl-audit-demo-badpod`** uses an invalid image → the kubelet emits **Warning** events (for example `Failed`, `FailedToPullImage`). **`kubectl-audit-demo-healthy`** may emit **Normal** events only, which the audit treats as benign.

## Try it

```bash
kubectl apply -f examples/audit-events/demo.yaml
kubectl wait --for=condition=available deployment/kubectl-audit-demo-badpod -n kubectl-audit-events-demo --timeout=120s || true
sleep 10
kubectl audit events -n kubectl-audit-events-demo
kubectl delete -f examples/audit-events/demo.yaml
```
