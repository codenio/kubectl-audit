package plugin

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func ListNonRunningPods(configFlags *genericclioptions.ConfigFlags) ([]corev1.Pod, error) {
	return ListNonRunningPodsWithScope(configFlags, true)
}

func ListNonRunningPodsWithScope(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) ([]corev1.Pod, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, allNamespaces)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	nonRunningPods := make([]corev1.Pod, 0)
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			nonRunningPods = append(nonRunningPods, pod)
		}
	}

	return nonRunningPods, nil
}

func ListProblemNodes(configFlags *genericclioptions.ConfigFlags) ([]corev1.Node, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	problemNodes := make([]corev1.Node, 0)
	for _, node := range nodes.Items {
		if isNotReadyNode(node) || node.Spec.Unschedulable {
			problemNodes = append(problemNodes, node)
		}
	}

	return problemNodes, nil
}

func ListNonBoundPVCs(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) ([]corev1.PersistentVolumeClaim, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, allNamespaces)
	if err != nil {
		return nil, err
	}

	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pvcs: %w", err)
	}

	nonBoundPVCs := make([]corev1.PersistentVolumeClaim, 0)
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase != corev1.ClaimBound {
			nonBoundPVCs = append(nonBoundPVCs, pvc)
		}
	}

	return nonBoundPVCs, nil
}

func ListProblemJobs(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) ([]batchv1.Job, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, allNamespaces)
	if err != nil {
		return nil, err
	}

	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	problemJobs := make([]batchv1.Job, 0)
	for _, job := range jobs.Items {
		if isProblemJob(job) {
			problemJobs = append(problemJobs, job)
		}
	}
	return problemJobs, nil
}

func ListSuspendedCronJobs(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) ([]batchv1beta1.CronJob, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespace, err := namespaceForQuery(configFlags, allNamespaces)
	if err != nil {
		return nil, err
	}

	cronJobs, err := clientset.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cronjobs: %w", err)
	}

	suspended := make([]batchv1beta1.CronJob, 0)
	for _, cj := range cronJobs.Items {
		if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
			suspended = append(suspended, cj)
		}
	}
	return suspended, nil
}

func ListNonBoundPVs(configFlags *genericclioptions.ConfigFlags) ([]corev1.PersistentVolume, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	pvs, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pvs: %w", err)
	}

	nonBound := make([]corev1.PersistentVolume, 0)
	for _, pv := range pvs.Items {
		if pv.Status.Phase != corev1.VolumeBound {
			nonBound = append(nonBound, pv)
		}
	}
	return nonBound, nil
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
