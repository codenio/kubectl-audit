package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codenio/kubectl-audit/pkg/plugin"
	"github.com/codenio/kubectl-audit/pkg/utils"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

func newPodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pods",
		Aliases: []string{"pod", "po"},
		Short:   "List non-running pods",
		RunE: func(cmd *cobra.Command, args []string) error {
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			wide, _ := cmd.Flags().GetBool("wide")
			namespace, _ := cmd.Flags().GetString("namespace")

			pods, err := plugin.ListNonRunningPodsWithScope(KubernetesConfigFlags, allNamespaces)
			if err != nil {
				return err
			}
			return printPodsTable(pods, allNamespaces, wide, namespace)
		},
	}
	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")
	cmd.Flags().Bool("wide", false, "Output in wide format")
	return cmd
}

func printPodsTable(pods []corev1.Pod, allNamespaces bool, wide bool, namespace string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if len(pods) == 0 {
		switch {
		case allNamespaces:
			fmt.Fprintln(w, "No unhealthy pods found in cluster")
		case namespace != "":
			fmt.Fprintf(w, "No unhealthy pods found in namespace %q\n", namespace)
		default:
			fmt.Fprintln(w, "No unhealthy pods found in current namespace")
		}
		return w.Flush()
	}

	if allNamespaces {
		if wide {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE\tIP\tNODE")
		} else {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
		}
	} else {
		if wide {
			fmt.Fprintln(w, "NAME\tREADY\tSTATUS\tRESTARTS\tAGE\tIP\tNODE")
		} else {
			fmt.Fprintln(w, "NAME\tREADY\tSTATUS\tRESTARTS\tAGE")
		}
	}

	for _, pod := range pods {
		status := podStatus(pod)
		if allNamespaces {
			if wide {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\t%s\n", pod.Namespace, pod.Name, podReady(pod), status, podRestarts(pod), utils.ShortAge(pod.CreationTimestamp.Time), pod.Status.PodIP, pod.Spec.NodeName)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n", pod.Namespace, pod.Name, podReady(pod), status, podRestarts(pod), utils.ShortAge(pod.CreationTimestamp.Time))
			}
		} else {
			if wide {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n", pod.Name, podReady(pod), status, podRestarts(pod), utils.ShortAge(pod.CreationTimestamp.Time), pod.Status.PodIP, pod.Spec.NodeName)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", pod.Name, podReady(pod), status, podRestarts(pod), utils.ShortAge(pod.CreationTimestamp.Time))
			}
		}
	}
	return w.Flush()
}

func podReady(pod corev1.Pod) string {
	total := len(pod.Spec.Containers)
	ready := 0
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, total)
}

func podRestarts(pod corev1.Pod) int32 {
	var restarts int32
	for _, cs := range pod.Status.ContainerStatuses {
		restarts += cs.RestartCount
	}
	return restarts
}

func podStatus(pod corev1.Pod) string {
	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	for i, initStatus := range pod.Status.InitContainerStatuses {
		if initStatus.State.Terminated != nil && initStatus.State.Terminated.ExitCode == 0 {
			continue
		}
		if initStatus.State.Terminated != nil {
			if initStatus.State.Terminated.Reason != "" {
				return "Init:" + initStatus.State.Terminated.Reason
			}
			return fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
		}
		if initStatus.State.Waiting != nil && initStatus.State.Waiting.Reason != "" {
			return "Init:" + initStatus.State.Waiting.Reason
		}
		return fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
	}

	hasRunning := false
	for _, cs := range pod.Status.ContainerStatuses {
		switch {
		case cs.State.Waiting != nil && cs.State.Waiting.Reason != "":
			// Match kubectl-style behavior where later containers can override status.
			reason = cs.State.Waiting.Reason
		case cs.State.Terminated != nil && cs.State.Terminated.Reason != "":
			reason = cs.State.Terminated.Reason
		case cs.State.Terminated != nil && cs.State.Terminated.Reason == "":
			if cs.State.Terminated.Signal != 0 {
				reason = fmt.Sprintf("Signal:%d", cs.State.Terminated.Signal)
			} else {
				reason = fmt.Sprintf("ExitCode:%d", cs.State.Terminated.ExitCode)
			}
		case cs.Ready && cs.State.Running != nil:
			hasRunning = true
		}
	}

	if reason == "Completed" && hasRunning {
		reason = "Running"
	}
	if pod.DeletionTimestamp != nil {
		if pod.Status.Reason == "NodeLost" {
			return "Unknown"
		}
		return "Terminating"
	}
	return reason
}
