---
title: Find Unhealthy Kubernetes Pods — kubectl-audit
description: Use kubectl audit pods to find pods that are not running, not ready, or have high restart counts.
---

# Pods

`kubectl audit pods` lists pods flagged by health checks — see [Audit Rules: Pods](../reference/audit-rules.md#pods) for the exact criteria.

![kubectl audit pods — summary and attention table](../assets/audit-pods.gif){ .demo-gif }
<p class="demo-gif-caption"><code>kubectl audit pods</code> on a demo namespace with healthy, failing, and completed workloads.</p>

## Demo

```bash
kubectl apply -f examples/audit-pods/demo.yaml
kubectl audit pods -n kubectl-audit-pods-demo -l demo=kubectl-audit-pods
kubectl delete -f examples/audit-pods/demo.yaml
```

## Sample output

=== "Single namespace"
    ```
    $ kubectl audit pods -n ns-prod
    -------------------------------------------------------
    Pod Audit summary: total = 9 benign = 3 attention = 6
    -------------------------------------------------------
    Pod that requires attention
    NAME                                            READY   STATUS                  RESTARTS   AGE
    workload-a-dep-54b6948c9c-pqx12                 0/1     ImagePullBackOff        0          8h
    workload-a-dep-785b496f5d-rst34                 0/1     ContainerCreating       0          11h
    svc-b-dep-6c87c74674-uvw56                      2/3     ErrImagePull            0          3h5m
    svc-b-dep-9f6d7b5cc-xyz99                       0/3     ContainerStatusUnknown  0          12h
    sidecar-c-dep-9c5c948dd-abccd                   0/1     CrashLoopBackOff        0          18h
    sidecar-c-dep-c45d7c6c5-efghi                   0/1     InvalidImageName        0          8h
    ```

=== "All namespaces"
    ```
    $ kubectl audit pods -A
    -------------------------------------------------------
    Pod Audit summary: total = 120 benign = 114 attention = 6
    -------------------------------------------------------
    Pod that requires attention
    NAMESPACE      NAME                                            READY   STATUS             RESTARTS   AGE
    ns-team-a      workload-dep-5f7677d8c-pqx12                    0/1     ImagePullBackOff   0          8h
    ns-team-a      workload-dep-6a8b9c0d-rst34                     0/1     ImagePullBackOff   0          11h
    ns-team-b      indexer-dep-7c8d9e0f-uvw56                      0/1     ImagePullBackOff   0          3h55m
    ns-shared      sidecar-dep-8b9c0d1e-xyz78                      0/1     ErrImagePull       0          45m
    ns-shared      batch-harness-dep-9c0d1e2f-ab901                1/1     Running            6          4h
    ns-monitoring  obs-collector-dep-0d1e2f3g-cd234                1/1     Running            18         2d
    ```
    
    !!! note
        `Running` pods appear when restart count reaches the threshold (≥ 5). The phase is healthy but the pod is worth investigating.

## Useful commands

```bash
kubectl audit pods -A                      # all namespaces
kubectl audit pods -n ns-prod -l app=api   # label filter
kubectl audit pods -o wide                 # extended columns
kubectl audit pods -o json | jq ...        # JSON for scripting
```

## Next steps

- Need per-container detail? → [kubectl audit containers](containers.md)
- Seeing high restarts? → Check [Events](events.md) for `BackOff` warnings
- See the full detection criteria → [Audit Rules: Pods](../reference/audit-rules.md#pods)
