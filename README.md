# kubectl-audit

`kubectl-audit` is a `kubectl` plugin for auditing Kubernetes resources, helping you quickly identify unhealthy or non-compliant resources in your cluster.

## Install

### Local development install

```bash
make install
```

This builds `bin/audit` and installs it as `~/.krew/bin/kubectl-audit`.

### Krew install

```bash
kubectl krew install audit
```

## Usage

```bash
kubectl audit
```

The root command shows help and available subcommands.

## Available audits

- `kubectl audit pods` (aliases: `pod`, `po`)  
  Show non-running pods.
- `kubectl audit nodes` (aliases: `node`, `no`)  
  Show `NotReady` or `SchedulingDisabled` nodes.
- `kubectl audit pvc` (aliases: `pvcs`, `persistentvolumeclaim`, `persistentvolumeclaims`)  
  Show non-`Bound` PVCs.
- `kubectl audit pv` (aliases: `pvs`, `persistentvolume`, `persistentvolumes`)  
  Show non-`Bound` PVs.
- `kubectl audit jobs` (alias: `job`)  
  Show failed jobs, including backoff/deadline failures.
- `kubectl audit cronjobs` (aliases: `cronjob`, `cj`)  
  Show suspended cronjobs.

## Common examples

```bash
# Show unhealthy pods in current namespace
kubectl audit pods

# Show unhealthy pods cluster-wide
kubectl audit po -A

# Namespace-specific audit
kubectl audit pods -n tester

# PVC and PV audits
kubectl audit pvc -A
kubectl audit pv

# Job and CronJob audits
kubectl audit job -A
kubectl audit cj -A
```

