
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

## How it works
The plugin queries Kubernetes resources using client-go and prints kubectl-style
tables for resources that require attention (for example pods that are not running/ready or have high container restart counts,
NotReady nodes, non-Bound PVC/PV, failed jobs, and suspended cronjobs). See the README for pod criteria and the restart threshold.