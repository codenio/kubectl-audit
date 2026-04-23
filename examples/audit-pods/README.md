# Example: `kubectl audit pods`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-pods-demo`**, label **`demo=kubectl-audit-pods`**.

- **`kubectl-audit-demo-healthy`** — running `pause` (benign).
- **`kubectl-audit-demo-badpod`** — invalid image → attention (e.g. `ImagePullBackOff` / not ready).

## Try it

```bash
kubectl apply -f examples/audit-pods/demo.yaml
kubectl wait --for=condition=available deployment/kubectl-audit-demo-healthy -n kubectl-audit-pods-demo --timeout=120s
# Allow time for the bad Deployment to fail image pull (~30–60s)
kubectl audit pods -n kubectl-audit-pods-demo -l demo=kubectl-audit-pods
kubectl delete -f examples/audit-pods/demo.yaml
```

For **`kubectl audit containers`**, use [`../audit-containers/`](../audit-containers/README.md) (failing init + healthy Pod) or the same apply/delete here for image-pull issues on the **main** container.
