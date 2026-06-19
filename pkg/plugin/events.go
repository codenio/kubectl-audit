package plugin

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// AuditEvents returns Warning events as an EventList. benignInScope counts Normal events.
func AuditEvents(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*corev1.EventList, int, int, error) {
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

	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), auditMetav1ListOptions(o))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list events: %w", err)
	}

	totalInScope := len(events.Items)
	benignInScope := 0
	filtered := make([]corev1.Event, 0, totalInScope)
	for i := range events.Items {
		if !eventNeedsAttention(events.Items[i]) {
			benignInScope++
		} else {
			filtered = append(filtered, events.Items[i])
		}
	}
	return &corev1.EventList{Items: filtered}, totalInScope, benignInScope, nil
}

func eventNeedsAttention(event corev1.Event) bool {
	return event.Type == corev1.EventTypeWarning
}
