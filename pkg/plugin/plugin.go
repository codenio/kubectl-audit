package plugin

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// podRestartAttentionThreshold is the minimum RestartCount on any container (or init container)
// that flags a pod for attention even when phase is Running and all containers are Ready.
const podRestartAttentionThreshold int32 = 5

// AuditOptions configures listing and filtering for audit commands.
type AuditOptions struct {
	AllNamespaces bool
	LabelSelector string
	// PodName, if set, restricts listing to pods with this exact metadata.name (used by audit containers).
	PodName string
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
	lo := metav1.ListOptions{LabelSelector: o.LabelSelector}
	if o.PodName != "" {
		lo.FieldSelector = fields.OneTermEqualSelector("metadata.name", o.PodName).String()
	}
	return lo
}

// AuditPods returns pods that need attention as a PodList. totalInScope is len(Items) from the
// unfiltered list; benignInScope counts running/ready pods without high restart counts.
func AuditPods(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PodList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, 0, 0, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list pods: %w", err)
	}

	totalInScope := len(pods.Items)
	benignInScope := 0
	filtered := make([]corev1.Pod, 0, totalInScope)
	for i := range pods.Items {
		if podIsBenign(pods.Items[i]) {
			benignInScope++
		} else {
			filtered = append(filtered, pods.Items[i])
		}
	}
	return &corev1.PodList{Items: filtered}, totalInScope, benignInScope, nil
}

// podIsHealthy reports Running phase with all containers ready (desired steady state).
func podIsHealthy(pod corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, cs := range pod.Status.ContainerStatuses {
		if !cs.Ready {
			return false
		}
	}
	return true
}

func podHasHighContainerRestarts(pod corev1.Pod) bool {
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.RestartCount >= podRestartAttentionThreshold {
			return true
		}
	}
	for _, cs := range pod.Status.InitContainerStatuses {
		if cs.RestartCount >= podRestartAttentionThreshold {
			return true
		}
	}
	return false
}

func podIsBenign(pod corev1.Pod) bool {
	return podIsHealthy(pod) && !podHasHighContainerRestarts(pod)
}

func podNeedsAttention(pod corev1.Pod) bool {
	return !podIsBenign(pod)
}

// AuditNodes returns NotReady or unschedulable nodes as a NodeList. benignInScope counts nodes that
// are Ready and schedulable (healthy).
func AuditNodes(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.NodeList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list nodes: %w", err)
	}

	totalInScope := len(nodes.Items)
	benignInScope := 0
	filtered := make([]corev1.Node, 0, totalInScope)
	for i := range nodes.Items {
		if nodeIsHealthy(nodes.Items[i]) {
			benignInScope++
		} else {
			filtered = append(filtered, nodes.Items[i])
		}
	}
	return &corev1.NodeList{Items: filtered}, totalInScope, benignInScope, nil
}

func nodeIsHealthy(node corev1.Node) bool {
	return !isNotReadyNode(node) && !node.Spec.Unschedulable
}

// AuditPV returns non-Bound persistent volumes as a PersistentVolumeList. benignInScope counts
// Bound volumes (healthy in-use storage).
func AuditPV(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PersistentVolumeList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	pvs, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list pvs: %w", err)
	}

	totalInScope := len(pvs.Items)
	benignInScope := 0
	filtered := make([]corev1.PersistentVolume, 0, totalInScope)
	for i := range pvs.Items {
		if pvs.Items[i].Status.Phase == corev1.VolumeBound {
			benignInScope++
		} else {
			filtered = append(filtered, pvs.Items[i])
		}
	}
	return &corev1.PersistentVolumeList{Items: filtered}, totalInScope, benignInScope, nil
}

// AuditPVC returns non-Bound claims as a PersistentVolumeClaimList. benignInScope counts Bound claims.
func AuditPVC(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.PersistentVolumeClaimList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, 0, 0, err
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list pvcs: %w", err)
	}

	totalInScope := len(pvcs.Items)
	benignInScope := 0
	filtered := make([]corev1.PersistentVolumeClaim, 0, totalInScope)
	for i := range pvcs.Items {
		if pvcs.Items[i].Status.Phase == corev1.ClaimBound {
			benignInScope++
		} else {
			filtered = append(filtered, pvcs.Items[i])
		}
	}
	return &corev1.PersistentVolumeClaimList{Items: filtered}, totalInScope, benignInScope, nil
}

// AuditJobs returns failed or backoff/deadline problem jobs as a JobList. benignInScope counts jobs
// with no failure/backoff/deadline problems.
func AuditJobs(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*batchv1.JobList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, 0, 0, err
	}

	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list jobs: %w", err)
	}

	totalInScope := len(jobs.Items)
	benignInScope := 0
	filtered := make([]batchv1.Job, 0, totalInScope)
	for i := range jobs.Items {
		if !isProblemJob(jobs.Items[i]) {
			benignInScope++
		} else {
			filtered = append(filtered, jobs.Items[i])
		}
	}
	return &batchv1.JobList{Items: filtered}, totalInScope, benignInScope, nil
}

// AuditCronJobs returns suspended cron jobs as a CronJobList (batch/v1). benignInScope counts
// CronJobs that are not suspended (desired to run on schedule).
func AuditCronJobs(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*batchv1.CronJobList, int, int, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, o.AllNamespaces)
	if err != nil {
		return nil, 0, 0, err
	}

	cronJobs, err := clientset.BatchV1().CronJobs(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list cronjobs: %w", err)
	}

	totalInScope := len(cronJobs.Items)
	benignInScope := 0
	filtered := make([]batchv1.CronJob, 0, totalInScope)
	for i := range cronJobs.Items {
		s := cronJobs.Items[i].Spec.Suspend
		if s == nil || !*s {
			benignInScope++
		} else {
			filtered = append(filtered, cronJobs.Items[i])
		}
	}
	return &batchv1.CronJobList{Items: filtered}, totalInScope, benignInScope, nil
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
