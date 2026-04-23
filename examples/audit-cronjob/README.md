# Example: `kubectl audit cronjobs`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-cronjob-demo`**, label **`demo=kubectl-audit-cronjob`**.

**`kubectl-audit-demo-suspended`** has **`spec.suspend: true`** → should appear in the audit immediately after apply.

## Try it

```bash
kubectl apply -f examples/audit-cronjob/demo.yaml
kubectl audit cronjobs -n kubectl-audit-cronjob-demo -l demo=kubectl-audit-cronjob
kubectl delete -f examples/audit-cronjob/demo.yaml
```
