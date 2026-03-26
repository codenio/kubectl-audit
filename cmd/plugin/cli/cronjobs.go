package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/codenio/kubectl-audit/pkg/utils"
	"github.com/spf13/cobra"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

func newCronJobsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cronjobs",
		Aliases: []string{"cronjob", "cj"},
		Short:   "List suspended cronjobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			namespace, _ := cmd.Flags().GetString("namespace")

			cronjobs, err := plugin.ListSuspendedCronJobs(KubernetesConfigFlags, allNamespaces)
			if err != nil {
				return err
			}
			return printCronJobsTable(cronjobs, allNamespaces, namespace)
		},
	}
	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")
	return cmd
}

func printCronJobsTable(cronjobs []batchv1beta1.CronJob, withNamespace bool, namespace string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if len(cronjobs) == 0 {
		switch {
		case withNamespace:
			fmt.Fprintln(w, "No suspended cronjobs found in cluster")
		case namespace != "":
			fmt.Fprintf(w, "No suspended cronjobs found in namespace %q\n", namespace)
		default:
			fmt.Fprintln(w, "No suspended cronjobs found in current namespace")
		}
		return w.Flush()
	}

	if withNamespace {
		fmt.Fprintln(w, "NAMESPACE\tNAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE")
	} else {
		fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE")
	}

	for _, cj := range cronjobs {
		last := "<none>"
		if cj.Status.LastScheduleTime != nil {
			last = utils.ShortAge(cj.Status.LastScheduleTime.Time)
		}
		suspend := "False"
		if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
			suspend = "True"
		}

		if withNamespace {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", cj.Namespace, cj.Name, cj.Spec.Schedule, suspend, len(cj.Status.Active), last, utils.ShortAge(cj.CreationTimestamp.Time))
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", cj.Name, cj.Spec.Schedule, suspend, len(cj.Status.Active), last, utils.ShortAge(cj.CreationTimestamp.Time))
		}
	}
	return w.Flush()
}
