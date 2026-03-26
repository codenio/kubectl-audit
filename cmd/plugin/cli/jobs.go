package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/codenio/kubectl-audit/pkg/utils"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
)

func newJobsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jobs",
		Aliases: []string{"job"},
		Short:   "List failed jobs or deadline/backoff failures",
		RunE: func(cmd *cobra.Command, args []string) error {
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			namespace, _ := cmd.Flags().GetString("namespace")

			jobs, err := plugin.ListProblemJobs(KubernetesConfigFlags, allNamespaces)
			if err != nil {
				return err
			}
			return printJobsTable(jobs, allNamespaces, namespace)
		},
	}
	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")
	return cmd
}

func printJobsTable(jobs []batchv1.Job, withNamespace bool, namespace string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if len(jobs) == 0 {
		switch {
		case withNamespace:
			fmt.Fprintln(w, "No failed/backoff/deadline-exceeded jobs found in cluster")
		case namespace != "":
			fmt.Fprintf(w, "No failed/backoff/deadline-exceeded jobs found in namespace %q\n", namespace)
		default:
			fmt.Fprintln(w, "No failed/backoff/deadline-exceeded jobs found in current namespace")
		}
		return w.Flush()
	}

	if withNamespace {
		fmt.Fprintln(w, "NAMESPACE\tNAME\tCOMPLETIONS\tDURATION\tAGE\tREASON")
	} else {
		fmt.Fprintln(w, "NAME\tCOMPLETIONS\tDURATION\tAGE\tREASON")
	}

	for _, job := range jobs {
		completions := fmt.Sprintf("%d/%d", job.Status.Succeeded, desiredCompletions(job))
		duration := jobDuration(job)
		reason := jobFailureReason(job)
		if withNamespace {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", job.Namespace, job.Name, completions, duration, utils.ShortAge(job.CreationTimestamp.Time), reason)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", job.Name, completions, duration, utils.ShortAge(job.CreationTimestamp.Time), reason)
		}
	}
	return w.Flush()
}

func desiredCompletions(job batchv1.Job) int32 {
	if job.Spec.Completions != nil {
		return *job.Spec.Completions
	}
	return 1
}

func jobDuration(job batchv1.Job) string {
	if job.Status.StartTime == nil {
		return "<none>"
	}
	end := time.Now()
	if job.Status.CompletionTime != nil {
		end = job.Status.CompletionTime.Time
	}
	if end.Before(job.Status.StartTime.Time) {
		return "<none>"
	}
	return shortDuration(end.Sub(job.Status.StartTime.Time))
}

func jobFailureReason(job batchv1.Job) string {
	for _, c := range job.Status.Conditions {
		if c.Reason != "" {
			return c.Reason
		}
	}
	if job.Status.Failed > 0 {
		return "Failed"
	}
	return "Unknown"
}

func shortDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
