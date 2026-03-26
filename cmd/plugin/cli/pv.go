package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/codenio/kubectl-audit/pkg/utils"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

func newPVCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pv",
		Aliases: []string{"pvs", "persistentvolume", "persistentvolumes"},
		Short:   "List non-Bound PVs",
		RunE: func(cmd *cobra.Command, args []string) error {
			pvs, err := plugin.ListNonBoundPVs(KubernetesConfigFlags)
			if err != nil {
				return err
			}
			return printPVTable(pvs)
		},
	}
	return cmd
}

func printPVTable(pvs []corev1.PersistentVolume) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if len(pvs) == 0 {
		fmt.Fprintln(w, "No non-Bound PVs found in cluster")
		return w.Flush()
	}

	fmt.Fprintln(w, "NAME\tCAPACITY\tACCESS MODES\tRECLAIM POLICY\tSTATUS\tCLAIM\tSTORAGECLASS\tREASON\tAGE")
	for _, pv := range pvs {
		capacity := ""
		if v, ok := pv.Spec.Capacity[corev1.ResourceStorage]; ok {
			capacity = v.String()
		}
		claim := "<none>"
		if pv.Spec.ClaimRef != nil {
			claim = fmt.Sprintf("%s/%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			pv.Name,
			capacity,
			pvAccessModes(pv.Spec.AccessModes),
			string(pv.Spec.PersistentVolumeReclaimPolicy),
			string(pv.Status.Phase),
			claim,
			pv.Spec.StorageClassName,
			pv.Status.Reason,
			utils.ShortAge(pv.CreationTimestamp.Time),
		)
	}
	return w.Flush()
}

func pvAccessModes(modes []corev1.PersistentVolumeAccessMode) string {
	if len(modes) == 0 {
		return ""
	}
	values := make([]string, 0, len(modes))
	for _, mode := range modes {
		values = append(values, string(mode))
	}
	return strings.Join(values, ",")
}
