
## Usage
The following assumes you have the plugin installed via

```shell
kubectl krew install audit
```

### List non-running pods in your current namespace

```shell
kubectl audit pods
```

### List non-running pods in another context

```shell
kubectl audit pods --context=context-name
```

## How it works
The plugin queries Kubernetes resources using client-go and prints kubectl-style
tables for resources that require attention (for example non-running pods,
NotReady nodes, non-Bound PVC/PV, failed jobs, and suspended cronjobs).