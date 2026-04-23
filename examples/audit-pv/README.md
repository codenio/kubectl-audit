# Example: `kubectl audit pv`

Manifests: [`demo.yaml`](demo.yaml) — cluster-scoped **PersistentVolume** only, label **`demo=kubectl-audit-pv`**.

**`kubectl-audit-pv-demo-pv`** is an unbound **hostPath** PV → phase **Available** (not **Bound**) → the plugin lists it as attention (any non-Bound phase).

## Try it

```bash
kubectl apply -f examples/audit-pv/demo.yaml
kubectl audit pv -l demo=kubectl-audit-pv
kubectl delete -f examples/audit-pv/demo.yaml
```

If apply fails with **Forbidden** / admission errors, your cluster may disallow **hostPath**; remove this demo or use a local cluster (kind, minikube).
