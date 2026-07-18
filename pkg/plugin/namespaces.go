package plugin

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// AuditNamespaces returns namespaces that need attention: Terminating phase, or Active with no
// workloads (pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, or CronJobs).
// benignInScope counts namespaces that are Active and contain at least one workload.
func AuditNamespaces(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.NamespaceList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	ctx := context.Background()
	list, err := clientset.CoreV1().Namespaces().List(ctx, auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list namespaces: %w", err)
	}

	withWorkloads, err := namespacesWithWorkloads(ctx, clientset)
	if err != nil {
		return nil, 0, 0, err
	}

	totalInScope := len(list.Items)
	benignInScope := 0
	filtered := make([]corev1.Namespace, 0, totalInScope)
	for i := range list.Items {
		ns := list.Items[i]
		_, hasWorkloads := withWorkloads[ns.Name]
		if namespaceNeedsAttention(ns, hasWorkloads) {
			filtered = append(filtered, ns)
		} else {
			benignInScope++
		}
	}
	return &corev1.NamespaceList{Items: filtered}, totalInScope, benignInScope, nil
}

func namespaceNeedsAttention(ns corev1.Namespace, hasWorkloads bool) bool {
	if ns.Status.Phase == corev1.NamespaceTerminating {
		return true
	}
	return !hasWorkloads
}

func namespacesWithWorkloads(ctx context.Context, cs kubernetes.Interface) (map[string]struct{}, error) {
	out := make(map[string]struct{})
	listOpts := metav1.ListOptions{}

	pods, err := cs.CoreV1().Pods("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	for i := range pods.Items {
		out[pods.Items[i].Namespace] = struct{}{}
	}

	deployments, err := cs.AppsV1().Deployments("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	for i := range deployments.Items {
		out[deployments.Items[i].Namespace] = struct{}{}
	}

	statefulSets, err := cs.AppsV1().StatefulSets("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulsets: %w", err)
	}
	for i := range statefulSets.Items {
		out[statefulSets.Items[i].Namespace] = struct{}{}
	}

	daemonSets, err := cs.AppsV1().DaemonSets("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list daemonsets: %w", err)
	}
	for i := range daemonSets.Items {
		out[daemonSets.Items[i].Namespace] = struct{}{}
	}

	replicaSets, err := cs.AppsV1().ReplicaSets("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list replicasets: %w", err)
	}
	for i := range replicaSets.Items {
		out[replicaSets.Items[i].Namespace] = struct{}{}
	}

	jobs, err := cs.BatchV1().Jobs("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	for i := range jobs.Items {
		out[jobs.Items[i].Namespace] = struct{}{}
	}

	cronJobs, err := cs.BatchV1().CronJobs("").List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list cronjobs: %w", err)
	}
	for i := range cronJobs.Items {
		out[cronJobs.Items[i].Namespace] = struct{}{}
	}

	return out, nil
}
