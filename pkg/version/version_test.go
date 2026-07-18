package version

import (
	"strings"
	"testing"
)

func TestInfoDefault(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, GitCommit, BuildDate
	t.Cleanup(func() {
		Version, GitCommit, BuildDate = oldVersion, oldCommit, oldDate
	})

	Version, GitCommit, BuildDate = "dev", "", ""
	got := Info()
	if got != "kubectl-audit version dev" {
		t.Fatalf("Info() = %q", got)
	}
}

func TestInfoRelease(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, GitCommit, BuildDate
	t.Cleanup(func() {
		Version, GitCommit, BuildDate = oldVersion, oldCommit, oldDate
	})

	Version = "v1.2.3"
	GitCommit = "abc1234"
	BuildDate = "2026-07-18T00:00:00Z"
	got := Info()
	for _, want := range []string{"kubectl-audit version v1.2.3", "commit abc1234", "built 2026-07-18T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Info() = %q, want substring %q", got, want)
		}
	}
}
