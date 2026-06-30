---
title: Audit Rules & Detection Logic
description: Every health check, threshold, and detection rule used by kubectl-audit.
---

# Audit Rules & Detection Logic

This page is the single authoritative reference for every rule `kubectl-audit` applies. Rules are fixed in code — there are no flags to adjust thresholds.

## Pods

A pod is flagged if **any** of the following are true:

| Rule | Detail |
| ---- | ------ |
| Phase ≠ `Running` | Includes `Pending`, `Failed`, `Unknown`. `Succeeded`/`Completed` job pods are also flagged — they are not in a running steady state. |
| Any regular container `Ready = false` | Checked via `containerStatuses[*].ready`. Init containers are not checked for readiness here (see Containers). |
| Any container `restartCount ≥ 5` | Applies to both regular and init containers. The threshold is **5** and is not configurable. |

!!! info "Why are Succeeded pods flagged?"
    `Succeeded` pods from completed Jobs are intentionally included — they consume namespace resources and are often cleaned up by a TTL controller. Their presence signals they haven't been collected.

## Containers

Each container row (init or app) is flagged if:

| Rule | Waiting reasons that trigger |
| ---- | ---------------------------- |
| Container is in `Waiting` state with a known error reason | `CrashLoopBackOff`, `ImagePullBackOff`, `ErrImagePull`, `InvalidImageName`, `CreateContainerConfigError`, `RunContainerError` |
| Container `ready = false` and pod is not in a terminal phase | Pod is still expected to be running |
| `restartCount ≥ 5` | Same threshold as pod-level |
| Init container in `Failed` state | `lastState.terminated.exitCode ≠ 0` |

## Nodes

A node is flagged if **any** condition is true:

| Rule | Detail |
| ---- | ------ |
| `Ready` condition = `False` or `Unknown` | The `NodeReady` condition is not True |
| `spec.unschedulable = true` | Node is cordoned — `kubectl cordon` sets this flag; shown as `SchedulingDisabled` |

## PersistentVolumeClaims (PVC)

| Rule | Flagged phases |
| ---- | -------------- |
| Phase ≠ `Bound` | `Pending` (not yet provisioned), `Lost` (backing PV lost) |

## PersistentVolumes (PV)

| Rule | Flagged phases |
| ---- | -------------- |
| Phase ≠ `Bound` | `Available` (no claim), `Released` (claim deleted, not yet reclaimed), `Failed` (reclaim failure) |

## Jobs

| Rule | Detail |
| ---- | ------ |
| `status.failed > 0` | Any failed attempt, including backoff limit reached and active deadline exceeded |

## CronJobs

| Rule | Detail |
| ---- | ------ |
| `spec.suspend = true` | CronJob is suspended — no new Jobs will be created |

## Services

| Rule | Detail |
| ---- | ------ |
| Selector has entries and no Pods match | The Service's `spec.selector` is non-empty, but no Pods in the same namespace have all those labels |
| **Skipped:** `type: ExternalName` | These route by DNS name, not Pod selector |
| **Skipped:** empty `spec.selector` | Headless services or manually-managed endpoints |

!!! tip
    The `-l`/`--selector` flag for `kubectl audit service` filters the **Service list**, not the Pod matching. Pod matching always uses all Pods in each namespace.

## Deployments

A deployment is flagged if **either** condition is true:

| Rule | Detail |
| ---- | ------ |
| `spec.replicas == 0` | Explicitly scaled to zero |
| `status.readyReplicas < desired` when `desired > 0` | `desired` = `spec.replicas` if set, otherwise **1** (Kubernetes API default when `spec.replicas` is `nil`) |

## Events

| Rule | Detail |
| ---- | ------ |
| `type = Warning` | All Warning events regardless of reason or source |
| **Excluded:** `type = Normal` | Normal events are informational; only Warnings indicate problems |
