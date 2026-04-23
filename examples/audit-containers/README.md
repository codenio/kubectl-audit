# Example: `kubectl audit containers`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-containers-demo`**, label **`demo=kubectl-audit-containers`**.

- **`kubectl-audit-demo-healthy`** — single `pause` container (benign rows).
- **`kubectl-audit-demo-bad-init`** — **init** container exits `1`, so the app container never becomes ready → the audit lists at least the **init** row (`TYPE` = `init-container`), and often the **app** row while the Pod is not terminal.

This is separate from [audit-pods](../audit-pods/README.md), which uses a **bad image** on the main container to exercise **`kubectl audit pods`**.

## Try it

```bash
kubectl apply -f examples/audit-containers/demo.yaml
kubectl wait --for=condition=available deployment/kubectl-audit-demo-healthy -n kubectl-audit-containers-demo --timeout=120s
# Allow a few seconds for the bad-init Pod to run init and retry
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers -o wide
kubectl delete -f examples/audit-containers/demo.yaml
```
