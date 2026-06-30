---
title: Troubleshooting kubectl-audit
description: Common errors and fixes when running kubectl-audit.
---

# Troubleshooting

## Plugin not found

**Error:**
```
error: unknown command "audit" for "kubectl"
```

**Causes and fixes:**

1. Plugin not installed — run `kubectl krew install codenio/audit`
2. `~/.krew/bin` not on `PATH`:
   ```bash
   export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
   ```
   Add this to your shell profile (`~/.zshrc`, `~/.bashrc`).
3. Verify the binary exists:
   ```bash
   ls ~/.krew/bin/kubectl-audit
   kubectl audit --help
   ```

---

## Empty output — nothing flagged

If `kubectl audit pods` shows `attention = 0` but you expect issues:

1. **Your pods are actually healthy** — the audit rules are conservative. Check [Audit Rules](../reference/audit-rules.md) for exact thresholds.
2. **Wrong namespace** — add `-A` to check all namespaces:
   ```bash
   kubectl audit pods -A
   ```
3. **Wrong context** — verify you're on the right cluster:
   ```bash
   kubectl config current-context
   kubectl audit pods -A --context=my-cluster
   ```

---

## Permission denied / Forbidden

**Error:**
```
Error from server (Forbidden): pods is forbidden: User "..." cannot list resource "pods"
```

You need read access to the resources being audited. See [RBAC & Permissions](../getting-started/rbac.md) for the minimum ClusterRole.

Quick check:
```bash
kubectl auth can-i list pods --all-namespaces
kubectl auth can-i list nodes
```

---

## Unexpected output for `containers`

**`-o custom-columns` not working:**  
The `containers` subcommand uses a custom printer — it does not support `custom-columns`, Go templates, or `jsonpath`. Use `-o json` or `-o yaml` for machine output, then pipe to `jq`.

```bash
kubectl audit containers -o json | jq '.items[] | {pod: .pod, name: .name, status: .status}'
```

---

## Suspended CronJobs not appearing

`kubectl audit cronjobs` only flags CronJobs where `spec.suspend = true`. If your CronJob is failing (not suspended), check Jobs instead:

```bash
kubectl audit jobs -A
```

---

## Service shows as healthy but traffic isn't routing

`kubectl audit service` only checks whether the Service's selector matches any Pods — it doesn't verify those Pods are `Ready` or actually serving traffic. After confirming the Service is matched:

```bash
kubectl audit pods -n <namespace> -l <service-selector>
kubectl audit containers -n <namespace> -l <service-selector>
```

---

## Output looks different from docs

The plugin uses the same column widths and formats as `kubectl get`. Column alignment can differ depending on terminal width and the length of names in your cluster. The audit summary line is always on stderr; redirect if needed:

```bash
kubectl audit pods 2>/dev/null    # suppress summary
kubectl audit pods 2>&1           # merge to stdout
```

---

## Filing a bug

If you hit an issue not covered here, open an issue at [github.com/codenio/kubectl-audit/issues](https://github.com/codenio/kubectl-audit/issues) with:

- `kubectl audit --version` output (or binary version)
- `kubectl version` output
- The exact command you ran and the full error
