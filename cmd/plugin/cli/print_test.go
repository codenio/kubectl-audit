package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWriteAuditSummary(t *testing.T) {
	var buf bytes.Buffer
	WriteAuditSummary(&buf, "Pod", 10, 7, 3)
	out := buf.String()
	if !strings.Contains(out, "Pod Audit summary: total = 10") {
		t.Fatalf("unexpected output: %q", out)
	}
	if !strings.Contains(out, "benign = 7") {
		t.Fatalf("unexpected output: %q", out)
	}
	if !strings.Contains(out, "attention = 3") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestPrintContainersOutputJSON(t *testing.T) {
	list := &plugin.ContainerList{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ContainerList"},
		Items: []plugin.Container{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "demo"},
				Pod:        "web-0",
				Status:     "Running",
			},
		},
	}

	out := captureStdout(t, func() error {
		return printContainersOutput(list, "json", false, false)
	})
	if !strings.Contains(out, `"kind": "ContainerList"`) {
		t.Fatalf("unexpected json: %s", out)
	}
	if !strings.Contains(out, `"pod": "web-0"`) {
		t.Fatalf("unexpected json: %s", out)
	}
}

func TestPrintContainersOutputName(t *testing.T) {
	list := &plugin.ContainerList{
		Items: []plugin.Container{
			{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "demo"}, Pod: "web-0"},
		},
	}
	out := captureStdout(t, func() error {
		return printContainersOutput(list, "name", false, false)
	})
	if strings.TrimSpace(out) != "demo/app" {
		t.Fatalf("got %q", out)
	}
}

func TestPrintContainersOutputUnsupported(t *testing.T) {
	list := &plugin.ContainerList{}
	err := printContainersOutput(list, "custom-columns=NAME:.metadata.name", false, false)
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expected unsupported format error, got %v", err)
	}
}

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	runErr := fn()
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	if runErr != nil {
		t.Fatal(runErr)
	}
	return buf.String()
}
