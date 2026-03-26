package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/codenio/kubectl-audit/pkg/utils"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

func newNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nodes",
		Aliases: []string{"node", "no"},
		Short:   "List NotReady or SchedulingDisabled nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			wide, _ := cmd.Flags().GetBool("wide")
			nodes, err := plugin.ListProblemNodes(KubernetesConfigFlags)
			if err != nil {
				return err
			}
			return printNodesTable(nodes, wide)
		},
	}
	cmd.Flags().Bool("wide", false, "Output in wide format")
	return cmd
}

func printNodesTable(nodes []corev1.Node, wide bool) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if wide {
		fmt.Fprintln(w, "NAME\tSTATUS\tROLES\tAGE\tVERSION\tINTERNAL-IP\tOS-IMAGE\tKERNEL-VERSION\tCONTAINER-RUNTIME")
	} else {
		fmt.Fprintln(w, "NAME\tSTATUS\tROLES\tAGE\tVERSION")
	}
	for _, node := range nodes {
		if wide {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", node.Name, nodeStatus(node), nodeRoles(node), utils.ShortAge(node.CreationTimestamp.Time), node.Status.NodeInfo.KubeletVersion, nodeAddress(node), node.Status.NodeInfo.OSImage, node.Status.NodeInfo.KernelVersion, node.Status.NodeInfo.ContainerRuntimeVersion)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", node.Name, nodeStatus(node), nodeRoles(node), utils.ShortAge(node.CreationTimestamp.Time), node.Status.NodeInfo.KubeletVersion)
		}
	}
	return w.Flush()
}

func nodeStatus(node corev1.Node) string {
	status := "Unknown"
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			if c.Status == corev1.ConditionTrue {
				status = "Ready"
			} else {
				status = "NotReady"
			}
			break
		}
	}
	if node.Spec.Unschedulable {
		status = status + ",SchedulingDisabled"
	}
	return status
}

func nodeRoles(node corev1.Node) string {
	roles := make([]string, 0)
	for key := range node.Labels {
		if strings.HasPrefix(key, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(key, "node-role.kubernetes.io/")
			if role == "" {
				role = "<none>"
			}
			roles = append(roles, role)
		}
	}
	if len(roles) == 0 {
		return "<none>"
	}
	sort.Strings(roles)
	return strings.Join(roles, ",")
}

func nodeAddress(node corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
