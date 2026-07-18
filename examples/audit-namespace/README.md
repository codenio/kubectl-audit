# Example: `kubectl audit ns`

Manifests: [`demo.yaml`](demo.yaml) — three namespaces labeled **`demo=kubectl-audit-namespace`**.

| Namespace | What it exercises |
| --------- | ----------------- |
| **`kubectl-audit-namespace-empty-demo`** | **Active** with ConfigMap, Secret, and Service only → **attention** (no workloads). |
| **`kubectl-audit-namespace-workload-demo`** | **Active** with a running Deployment → **benign** (not listed). |
| **`kubectl-audit-namespace-terminating-demo`** | Custom finalizer blocks delete → **Terminating** after you delete the namespace → **attention**. |

Workloads checked: Pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, and CronJobs. ConfigMaps, Secrets, and Services alone do **not** count.

## Prerequisites (kind + podman)

The workload namespace uses `registry.k8s.io/pause:3.9`. On a local **kind** cluster with **podman**, pre-load it once:

```bash
podman pull --tls-verify=false registry.k8s.io/pause:3.9
podman save registry.k8s.io/pause:3.9 -o /tmp/pause-3.9.tar
KIND_EXPERIMENTAL_PROVIDER=podman kind load image-archive /tmp/pause-3.9.tar --name kind
rm -f /tmp/pause-3.9.tar
```

## Try it

**1. Apply**

```bash
kubectl apply -f examples/audit-namespace/demo.yaml
kubectl wait --for=condition=available deployment/healthy -n kubectl-audit-namespace-workload-demo --timeout=120s
```

**2. Start terminating demo** (custom finalizer keeps the namespace stuck in **Terminating**)

```bash
kubectl delete namespace kubectl-audit-namespace-terminating-demo --wait=false
```

**3. Audit** (label filter scopes to demo namespaces only)

```bash
kubectl audit ns -l demo=kubectl-audit-namespace
kubectl audit namespace -l demo=kubectl-audit-namespace -o wide
```

**4. Optional**

```bash
kubectl audit ns -l demo=kubectl-audit-namespace -o json
kubectl audit ns -l demo=kubectl-audit-namespace -o name
```

**5. Cleanup** (remove the blocking finalizer before deleting the stuck namespace)

```bash
kubectl patch namespace kubectl-audit-namespace-terminating-demo -p '{"metadata":{"finalizers":[]}}' --type=merge 2>/dev/null || true
kubectl delete -f examples/audit-namespace/demo.yaml --ignore-not-found
```

You should see **`kubectl-audit-namespace-empty-demo`** (**Active**) and **`kubectl-audit-namespace-terminating-demo`** (**Terminating**) in the attention table. **`kubectl-audit-namespace-workload-demo`** should not appear.
