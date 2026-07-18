package version

import (
	"fmt"
	"strings"
)

var (
	// Version is the plugin release tag or "dev" for local builds.
	Version = "dev"
	// GitCommit is the git SHA embedded at build time (optional).
	GitCommit = ""
	// BuildDate is the build timestamp in UTC (optional).
	BuildDate = ""
)

// Info returns a single-line version string for CLI output.
func Info() string {
	var b strings.Builder
	fmt.Fprintf(&b, "kubectl-audit version %s", Version)
	if GitCommit != "" {
		fmt.Fprintf(&b, " (commit %s", GitCommit)
		if BuildDate != "" {
			fmt.Fprintf(&b, ", built %s", BuildDate)
		}
		b.WriteString(")")
	} else if BuildDate != "" {
		fmt.Fprintf(&b, " (built %s)", BuildDate)
	}
	return b.String()
}
