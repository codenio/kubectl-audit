# Example: `kubectl audit jobs`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-job-demo`**, label **`demo=kubectl-audit-job`**.

**`kubectl-audit-demo-failjob`** exits with status 1 and `backoffLimit: 0` → should appear in the audit after the Job fails.

## Try it

```bash
kubectl apply -f examples/audit-job/demo.yaml
kubectl wait --for=condition=failed job/kubectl-audit-demo-failjob -n kubectl-audit-job-demo --timeout=120s
kubectl audit jobs -n kubectl-audit-job-demo -l demo=kubectl-audit-job
kubectl delete -f examples/audit-job/demo.yaml
```
