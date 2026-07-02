---
title: Kubernetes Cluster Health Plugin
description: kubectl plugin to find unhealthy pods, nodes, storage, services, deployments and warning events from the command line.
hide:
  - navigation
  - toc
---

<div class="hero">
  <img src="assets/logo-home.png" alt="kubectl-audit Kubernetes health check plugin" class="hero-logo" />
  <h1>kubectl-audit</h1>
  <p class="tagline">Kubernetes cluster health, straight from your terminal.</p>
  <div class="install-pill">kubectl krew install codenio/audit</div>
</div>

`kubectl-audit` is a [`kubectl` plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) that surfaces Kubernetes resources failing common health checks — without switching tools or writing custom scripts.

Use it for **incident triage**, **SRE runbooks**, **pre-deploy smoke checks**, and **CI/CD gates**.

![kubectl audit pods — terminal demo](assets/audit-pods.gif){ .demo-gif }
<p class="demo-gif-caption"><code>kubectl audit pods</code> surfaces unhealthy pods with a summary line and attention table.</p>

---

## What it checks

| Resource | Flagged when… |
| -------- | -------------- |
| **Pods** | Phase ≠ `Running`, any container not `Ready`, or restart count ≥ 5 |
| **Containers** | Per-container: image pull errors, crash loops, not ready, high restarts |
| **Nodes** | `NotReady` or `SchedulingDisabled` (cordoned) |
| **PV / PVC** | Not in `Bound` phase (`Pending`, `Lost`, `Released`, `Failed`) |
| **Jobs** | Failed — including backoff and deadline limit failures |
| **CronJobs** | `suspend: true` |
| **Services** | Non-empty selector with no matching Pods in the namespace |
| **Deployments** | `spec.replicas: 0`, or `readyReplicas` < desired |
| **Events** | `type: Warning` |

See [Audit Rules](reference/audit-rules.md) for exact detection thresholds.

---

## Navigate the docs

<div class="grid-cards">
  <a href="getting-started/install/" class="grid-card">
    <div class="card-icon">🚀</div>
    <h3>Installation</h3>
    <p>Krew, from source, verify</p>
  </a>
  <a href="getting-started/quickstart/" class="grid-card">
    <div class="card-icon">⚡</div>
    <h3>Quick Start</h3>
    <p>Common workflows and first commands</p>
  </a>
  <a href="reference/audit-rules/" class="grid-card">
    <div class="card-icon">📋</div>
    <h3>Audit Rules</h3>
    <p>Every check, threshold, and detection logic</p>
  </a>
  <a href="reference/subcommands/" class="grid-card">
    <div class="card-icon">🔧</div>
    <h3>Subcommands</h3>
    <p>Full subcommand and flag reference</p>
  </a>
  <a href="guides/troubleshooting/" class="grid-card">
    <div class="card-icon">🔍</div>
    <h3>Troubleshooting</h3>
    <p>Common errors and fixes</p>
  </a>
  <a href="guides/cicd/" class="grid-card">
    <div class="card-icon">🔄</div>
    <h3>CI/CD Integration</h3>
    <p>GitHub Actions, exit codes, scripting</p>
  </a>
</div>
