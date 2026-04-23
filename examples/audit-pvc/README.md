# Example: `kubectl audit pvc`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-pvc-demo`**, label **`demo=kubectl-audit-pvc`**.

**`kubectl-audit-demo-pvc-pending`** references StorageClass **`kubectl-audit-nonexistent-sc`** (not created) → stays **Pending** → **not Bound** → should appear in the audit.

## Try it

```bash
kubectl apply -f examples/audit-pvc/demo.yaml
kubectl audit pvc -n kubectl-audit-pvc-demo -l demo=kubectl-audit-pvc
kubectl delete -f examples/audit-pvc/demo.yaml
```
