package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestAuditSummaryResourceTitle(t *testing.T) {
	if got := auditSummaryResourceTitle("pods"); got != "Pod" {
		t.Fatalf("got %q", got)
	}
	if got := auditSummaryResourceTitle("unknown"); got != "Resource" {
		t.Fatalf("got %q", got)
	}
}

func TestWithNamespaceColumn(t *testing.T) {
	if withNamespaceColumn("pods", false) {
		t.Fatal("expected false without -A")
	}
	if !withNamespaceColumn("pods", true) {
		t.Fatal("expected true for pods with -A")
	}
	if withNamespaceColumn("nodes", true) {
		t.Fatal("expected false for cluster-scoped nodes")
	}
	if withNamespaceColumn("namespaces", true) {
		t.Fatal("expected false for namespaces")
	}
}

func TestIsHumanTableOutput(t *testing.T) {
	for _, out := range []string{"", "wide"} {
		if !isHumanTableOutput(out) {
			t.Fatalf("expected human table for %q", out)
		}
	}
	if isHumanTableOutput("json") {
		t.Fatal("expected false for json")
	}
}

func TestAuditObjectLen(t *testing.T) {
	if got := auditObjectLen(&corev1.PodList{Items: []corev1.Pod{{}, {}}}); got != 2 {
		t.Fatalf("got %d", got)
	}
	if got := auditObjectLen(&metav1.Table{Rows: []metav1.TableRow{{}, {}, {}}}); got != 3 {
		t.Fatalf("got %d", got)
	}
}

func TestAuditResourcePhrase(t *testing.T) {
	if got := auditResourcePhrase("namespaces"); got != "namespaces" {
		t.Fatalf("got %q", got)
	}
	if got := auditResourcePhrase("does-not-exist"); got != "resources" {
		t.Fatalf("got %q", got)
	}
}

func TestAllNamespacesInArgv(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"kubectl-audit", "pods", "-A"}
	if !allNamespacesInArgv() {
		t.Fatal("expected -A to be detected")
	}

	os.Args = []string{"kubectl-audit", "pods", "--all-namespaces=true"}
	if !allNamespacesInArgv() {
		t.Fatal("expected --all-namespaces=true")
	}

	os.Args = []string{"kubectl-audit", "pods", "--all-namespaces=false"}
	if allNamespacesInArgv() {
		t.Fatal("expected false for --all-namespaces=false")
	}

	os.Args = []string{"kubectl-audit", "pods"}
	if allNamespacesInArgv() {
		t.Fatal("expected false without flag")
	}
}

func TestWriteAuditEmptyMessage(t *testing.T) {
	ns := "team-a"
	cf := genericclioptions.NewConfigFlags(false)
	cf.Namespace = &ns

	var buf bytes.Buffer
	writeAuditEmptyMessage(&buf, "pods", false, cf, 0)
	if !strings.Contains(buf.String(), "team-a") {
		t.Fatalf("unexpected message: %q", buf.String())
	}

	buf.Reset()
	writeAuditEmptyMessage(&buf, "nodes", false, cf, 2)
	if !strings.Contains(buf.String(), "require attention") {
		t.Fatalf("unexpected message: %q", buf.String())
	}
}
