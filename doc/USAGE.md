
## Usage
The following assumes you have the plugin installed via

```shell
kubectl krew install audit
```

### List pods that need attention in your current namespace

```shell
kubectl audit pods
```

### List pods that need attention in another context

```shell
kubectl audit pods --context=context-name
```

### List Services whose selector matches no Pods

```shell
kubectl audit service
kubectl audit svc -A
```

Try the demo manifests in the repository (per audit under `examples/audit-*`; index at `examples/README.md`). For services: `examples/audit-svc/demo.yaml`.

### Try a demo (example: pods)

```shell
kubectl apply -f examples/audit-pods/demo.yaml
kubectl audit pods -n kubectl-audit-pods-demo -l demo=kubectl-audit-pods
kubectl delete -f examples/audit-pods/demo.yaml
```

### Try containers (per-container rows)

```shell
kubectl apply -f examples/audit-containers/demo.yaml
kubectl audit containers -n kubectl-audit-containers-demo -l demo=kubectl-audit-containers
kubectl delete -f examples/audit-containers/demo.yaml
```

See `examples/README.md` for **jobs**, **cronjobs**, **PVC**, **PV**, and **service** folders.

## How it works
The plugin queries Kubernetes resources using client-go and prints kubectl-style
tables for resources that require attention (for example pods that are not running/ready or have high container restart counts,
NotReady nodes, non-Bound PVC/PV, failed jobs, suspended cronjobs, and Services with a pod selector but no matching Pods). See the README for pod criteria and the restart threshold.