---
title: Find Failed Kubernetes Jobs and Suspended CronJobs — kubectl-audit
description: Use kubectl audit jobs and kubectl audit cronjobs to find failed jobs and suspended cron jobs.
---

# Jobs & CronJobs

## Jobs

`kubectl audit jobs` lists Jobs with `status.failed > 0` — including those that hit their backoff limit or active deadline.

### Demo

```bash
kubectl apply -f examples/audit-job/demo.yaml
kubectl audit jobs -n kubectl-audit-job-demo
kubectl delete -f examples/audit-job/demo.yaml
```

### Sample output

```
$ kubectl audit jobs -A
-------------------------------------------------------
Job Audit summary: total = 28 benign = 26 attention = 2
-------------------------------------------------------
Job that requires attention
NAMESPACE    NAME               COMPLETIONS   DURATION   AGE
batch-demo   daily-import       0/1           5m         5m
batch-demo   retry-migrate      0/1           1h         1h
```

### Useful commands

```bash
kubectl audit jobs -A
kubectl audit job -n batch-ns
kubectl audit jobs -o json | jq '.items[].metadata.name'
```

## CronJobs

`kubectl audit cronjobs` lists CronJobs where `spec.suspend = true`. Suspended CronJobs will not create new Jobs until unsuspended.

### Demo

```bash
kubectl apply -f examples/audit-cronjob/demo.yaml
kubectl audit cronjobs -n kubectl-audit-cronjob-demo
kubectl delete -f examples/audit-cronjob/demo.yaml
```

### Sample output

```
$ kubectl audit cronjobs -A
-------------------------------------------------------
CronJob Audit summary: total = 15 benign = 13 attention = 2
-------------------------------------------------------
CronJob that requires attention
NAMESPACE   NAME               SCHEDULE      SUSPEND   ACTIVE
ops-demo    pause-backup       0 2 * * *     True      0
ops-demo    hold-reports       15 * * * *    True      0
```

### Useful commands

```bash
kubectl audit cj -A
kubectl audit cronjobs -n ops-ns
```

## Next steps

- See detection rules → [Audit Rules: Jobs](../reference/audit-rules.md#jobs) / [CronJobs](../reference/audit-rules.md#cronjobs)
- Failed Job? Check the pods it created:
  ```bash
  kubectl audit pods -n <namespace> -l job-name=<job-name>
  ```
