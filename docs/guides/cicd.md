---
title: CI/CD Integration
description: Integrate kubectl-audit into GitHub Actions, GitLab CI, and shell scripts as a pre-deploy health gate.
---

# CI/CD Integration

`kubectl-audit` works well as a **pre-deploy gate** or **post-deploy health check**. The audit summary goes to stderr; stdout is clean JSON/YAML for scripting.

## Exit codes

| Exit code | Meaning |
| --------- | ------- |
| `0` | Command ran successfully (even if issues were found — the plugin reports, not enforces) |
| non-zero | Command failed (API error, permission denied, invalid flag) |

!!! warning "Enforcing as a gate"
    `kubectl-audit` returns exit code `0` even when issues are found — it is a reporter, not a policy enforcer. To gate on findings, check the count via `-o json` and `jq` (see examples below).

## Shell scripting

### Gate on pod health

```bash
# Fail the script if any unhealthy pods exist in a namespace
ISSUES=$(kubectl audit pods -n my-namespace -o json 2>/dev/null | jq '.items | length')
if [ "$ISSUES" -gt 0 ]; then
  echo "❌ $ISSUES unhealthy pod(s) found"
  kubectl audit pods -n my-namespace
  exit 1
fi
echo "✅ All pods healthy"
```

### Post-deploy check

```bash
# Wait a few seconds after deploy, then audit
sleep 15
kubectl audit pods -n production -o json 2>/dev/null \
  | jq -r '.items[] | "\(.metadata.namespace)/\(.metadata.name): \(.status.phase)"'
```

## GitHub Actions

### Pre-deploy gate

```yaml title=".github/workflows/deploy.yml"
name: deploy
on:
  push:
    branches: [main]

jobs:
  pre-deploy-audit:
    runs-on: ubuntu-latest
    steps:
      - name: Set up kubectl
        uses: azure/setup-kubectl@v4

      - name: Install kubectl-audit
        run: |
          kubectl krew index add codenio https://github.com/codenio/codenio-krew-index.git
          kubectl krew install codenio/audit

      - name: Pre-deploy health check
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG }}
        run: |
          ISSUES=$(kubectl audit pods -n production -o json 2>/dev/null | jq '.items | length')
          if [ "$ISSUES" -gt 0 ]; then
            echo "::error::$ISSUES unhealthy pod(s) in production — aborting deploy"
            kubectl audit pods -n production
            exit 1
          fi

  deploy:
    needs: pre-deploy-audit
    runs-on: ubuntu-latest
    steps:
      - name: Deploy
        run: echo "Deploying..."
```

### Post-deploy audit comment on PR

```yaml title=".github/workflows/audit-comment.yml"
name: audit-comment
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - name: Install kubectl-audit
        run: |
          kubectl krew index add codenio https://github.com/codenio/codenio-krew-index.git
          kubectl krew install codenio/audit

      - name: Run audit
        id: audit
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG }}
        run: |
          OUTPUT=$(kubectl audit pods -A 2>&1)
          echo "result<<EOF" >> $GITHUB_OUTPUT
          echo "$OUTPUT" >> $GITHUB_OUTPUT
          echo "EOF" >> $GITHUB_OUTPUT

      - name: Comment on PR
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `## Cluster Audit\n\`\`\`\n${{ steps.audit.outputs.result }}\n\`\`\``
            })
```

## GitLab CI

```yaml title=".gitlab-ci.yml"
audit:
  stage: validate
  image: bitnami/kubectl:latest
  before_script:
    - kubectl krew index add codenio https://github.com/codenio/codenio-krew-index.git
    - kubectl krew install codenio/audit
  script:
    - |
      ISSUES=$(kubectl audit pods -n $CI_ENVIRONMENT_SLUG -o json 2>/dev/null | jq '.items | length')
      if [ "$ISSUES" -gt 0 ]; then
        kubectl audit pods -n $CI_ENVIRONMENT_SLUG
        exit 1
      fi
  environment:
    name: production
```

## Parsing output with jq

```bash
# Names of unhealthy pods
kubectl audit pods -A -o json 2>/dev/null | jq -r '.items[].metadata.name'

# Pods grouped by namespace
kubectl audit pods -A -o json 2>/dev/null \
  | jq 'group_by(.metadata.namespace) | map({namespace: .[0].metadata.namespace, count: length})'

# Services with no backing pods
kubectl audit service -A -o json 2>/dev/null \
  | jq -r '.items[] | "\(.metadata.namespace)/\(.metadata.name)"'
```
