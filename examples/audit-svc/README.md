# Example: `kubectl audit service`

Manifests: [`demo.yaml`](demo.yaml) — namespace **`kubectl-audit-service-demo`**, label **`demo=kubectl-audit-service`**.

- **`orphan-svc`** — selector matches no Pods → **attention**.
- **`healthy-svc`** + **`backend`** Deployment — selector matches Pods → not listed.
- **`external-only`** — `ExternalName` → skipped by the audit.
- **`no-selector`** — empty `spec.selector` → skipped.

## Try it

**1. Apply**

```bash
kubectl apply -f examples/audit-svc/demo.yaml
```

**2. Wait** so `healthy-svc` has backends before you audit.

```bash
kubectl wait --for=condition=available deployment/backend -n kubectl-audit-service-demo --timeout=120s
```

**3. Audit**

```bash
kubectl audit service -n kubectl-audit-service-demo
kubectl audit svc -A -l demo=kubectl-audit-service -o wide
```

**4. Optional**

```bash
kubectl audit service -n kubectl-audit-service-demo -o json
kubectl audit service -n kubectl-audit-service-demo -o name
```

**5. Cleanup**

```bash
kubectl delete -f examples/audit-svc/demo.yaml
```

You should see **`orphan-svc`** only in the attention table (steps 3–4).
