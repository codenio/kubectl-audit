package plugin

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// SummarizeAudit returns summary counts for the audit line:
//   - total: all resources in scope (same list as the audit)
//   - benign: healthy / desired by kind (from each Audit*)
//   - attention: rows in the filtered list (matches the printed table)
func SummarizeAudit(totalInScope, benignInScope int, obj runtime.Object) (total, benign, attention int) {
	attention = meta.LenList(obj)
	return totalInScope, benignInScope, attention
}
