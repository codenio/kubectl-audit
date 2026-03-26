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

func newPVCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pvc",
		Aliases: []string{"pvcs", "persistentvolumeclaim", "persistentvolumeclaims"},
		Short:   "List non-Bound PVCs",
		RunE: func(cmd *cobra.Command, args []string) error {
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			wide, _ := cmd.Flags().GetBool("wide")
			namespace, _ := cmd.Flags().GetString("namespace")

			pvcs, err := plugin.ListNonBoundPVCs(KubernetesConfigFlags, allNamespaces)
			if err != nil {
				return err
			}
			return printPVCTable(pvcs, allNamespaces, wide, namespace)
		},
	}
	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")
	cmd.Flags().Bool("wide", false, "Output in wide format")
	return cmd
}

func printPVCTable(pvcs []corev1.PersistentVolumeClaim, allNamespaces bool, wide bool, namespace string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if len(pvcs) == 0 {
		switch {
		case allNamespaces:
			fmt.Fprintln(w, "No non-Bound PVCs found in cluster")
		case namespace != "":
			fmt.Fprintf(w, "No non-Bound PVCs found in namespace %q\n", namespace)
		default:
			fmt.Fprintln(w, "No non-Bound PVCs found in current namespace")
		}
		return w.Flush()
	}

	if allNamespaces {
		if wide {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tSTATUS\tVOLUME\tCAPACITY\tACCESS MODES\tSTORAGECLASS\tAGE\tVOLUMEMODE")
		} else {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tSTATUS\tVOLUME\tCAPACITY\tACCESS MODES\tSTORAGECLASS\tAGE")
		}
	} else {
		if wide {
			fmt.Fprintln(w, "NAME\tSTATUS\tVOLUME\tCAPACITY\tACCESS MODES\tSTORAGECLASS\tAGE\tVOLUMEMODE")
		} else {
			fmt.Fprintln(w, "NAME\tSTATUS\tVOLUME\tCAPACITY\tACCESS MODES\tSTORAGECLASS\tAGE")
		}
	}
	for _, pvc := range pvcs {
		storageClass := ""
		if pvc.Spec.StorageClassName != nil {
			storageClass = *pvc.Spec.StorageClassName
		}

		capacity := ""
		if v, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
			capacity = v.String()
		}

		volumeMode := ""
		if wide && pvc.Spec.VolumeMode != nil {
			volumeMode = string(*pvc.Spec.VolumeMode)
		}
		if allNamespaces {
			if wide {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pvc.Namespace, pvc.Name, string(pvc.Status.Phase), pvc.Spec.VolumeName, capacity, accessModes(pvc.Spec.AccessModes), storageClass, utils.ShortAge(pvc.CreationTimestamp.Time), volumeMode)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pvc.Namespace, pvc.Name, string(pvc.Status.Phase), pvc.Spec.VolumeName, capacity, accessModes(pvc.Spec.AccessModes), storageClass, utils.ShortAge(pvc.CreationTimestamp.Time))
			}
		} else {
			if wide {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pvc.Name, string(pvc.Status.Phase), pvc.Spec.VolumeName, capacity, accessModes(pvc.Spec.AccessModes), storageClass, utils.ShortAge(pvc.CreationTimestamp.Time), volumeMode)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pvc.Name, string(pvc.Status.Phase), pvc.Spec.VolumeName, capacity, accessModes(pvc.Spec.AccessModes), storageClass, utils.ShortAge(pvc.CreationTimestamp.Time))
			}
		}
	}
	return w.Flush()
}

func accessModes(modes []corev1.PersistentVolumeAccessMode) string {
	if len(modes) == 0 {
		return ""
	}
	out := make([]string, 0, len(modes))
	for _, m := range modes {
		out = append(out, string(m))
	}
	return strings.Join(out, ",")
}
