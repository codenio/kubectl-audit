package plugin

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// AuditOptions configures listing and filtering for audit commands.
type AuditOptions struct {
	AllNamespaces bool
	LabelSelector string
	// Pod filter flags (optional). If all false, default is phase != Running.
	PodPending  bool
	PodFailed   bool
	PodNotReady bool
}

func namespaceForQuery(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) (string, error) {
	if allNamespaces {
		return "", nil
	}

	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", fmt.Errorf("failed to resolve namespace: %w", err)
	}
	return namespace, nil
}

func auditMetav1ListOptions(o AuditOptions) metav1.ListOptions {
	return metav1.ListOptions{LabelSelector: o.LabelSelector}
}

// AuditPods returns pods matching audit criteria as a PodList.
func AuditPods(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PodList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	filtered := make([]corev1.Pod, 0)
	for _, pod := range pods.Items {
		if podMatchesAudit(pod, o) {
			filtered = append(filtered, pod)
		}
	}
	return &corev1.PodList{Items: filtered}, nil
}

func podMatchesAudit(pod corev1.Pod, o AuditOptions) bool {
	if !o.PodPending && !o.PodFailed && !o.PodNotReady {
		return pod.Status.Phase != corev1.PodRunning
	}
	if o.PodPending && pod.Status.Phase == corev1.PodPending {
		return true
	}
	if o.PodFailed && pod.Status.Phase == corev1.PodFailed {
		return true
	}
	if o.PodNotReady && podIsRunningButNotReady(pod) {
		return true
	}
	return false
}

func podIsRunningButNotReady(pod corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	total := len(pod.Spec.Containers)
	if total == 0 {
		return false
	}
	ready := 0
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			ready++
		}
	}
	return ready < total
}

// AuditNodes returns NotReady or unschedulable nodes as a NodeList.
func AuditNodes(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.NodeList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	filtered := make([]corev1.Node, 0)
	for _, node := range nodes.Items {
		if isNotReadyNode(node) || node.Spec.Unschedulable {
			filtered = append(filtered, node)
		}
	}
	return &corev1.NodeList{Items: filtered}, nil
}

// AuditPV returns non-Bound persistent volumes as a PersistentVolumeList.
func AuditPV(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PersistentVolumeList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list pvs: %w", err)
	}

	filtered := make([]corev1.PersistentVolume, 0)
	for _, pv := range pvs.Items {
		if pv.Status.Phase != corev1.VolumeBound {
			filtered = append(filtered, pv)
		}
	}
	return &corev1.PersistentVolumeList{Items: filtered}, nil
}

// AuditPVC returns non-Bound claims as a PersistentVolumeClaimList.
func AuditPVC(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PersistentVolumeClaimList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, err
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list pvcs: %w", err)
	}

	filtered := make([]corev1.PersistentVolumeClaim, 0)
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase != corev1.ClaimBound {
			filtered = append(filtered, pvc)
		}
	}
	return &corev1.PersistentVolumeClaimList{Items: filtered}, nil
}

// AuditJobs returns failed or backoff/deadline problem jobs as a JobList.
func AuditJobs(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*batchv1.JobList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, err
	}

	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	filtered := make([]batchv1.Job, 0)
	for _, job := range jobs.Items {
		if isProblemJob(job) {
			filtered = append(filtered, job)
		}
	}
	return &batchv1.JobList{Items: filtered}, nil
}

// AuditCronJobs returns suspended cron jobs as a CronJobList (batch/v1).
func AuditCronJobs(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*batchv1.CronJobList, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, err
	}

	cronJobs, err := clientset.BatchV1().CronJobs(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, fmt.Errorf("failed to list cronjobs: %w", err)
	}

	filtered := make([]batchv1.CronJob, 0)
	for _, cj := range cronJobs.Items {
		if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
			filtered = append(filtered, cj)
		}
	}
	return &batchv1.CronJobList{Items: filtered}, nil
}

func isNotReadyNode(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status != corev1.ConditionTrue
		}
	}
	return true
}

func isProblemJob(job batchv1.Job) bool {
	if job.Status.Failed > 0 {
		return true
	}

	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return true
		}
		if c.Reason == "BackoffLimitExceeded" || c.Reason == "DeadlineExceeded" {
			return true
		}
	}
	return false
}
