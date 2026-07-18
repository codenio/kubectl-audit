package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestNamespaceNeedsAttention(t *testing.T) {
	cases := []struct {
		name         string
		phase        corev1.NamespacePhase
		hasWorkloads bool
		want         bool
	}{
		{name: "active with workloads", phase: corev1.NamespaceActive, hasWorkloads: true, want: false},
		{name: "empty phase with workloads", phase: "", hasWorkloads: true, want: false},
		{name: "active without workloads", phase: corev1.NamespaceActive, hasWorkloads: false, want: true},
		{name: "empty phase without workloads", phase: "", hasWorkloads: false, want: true},
		{name: "terminating with workloads", phase: corev1.NamespaceTerminating, hasWorkloads: true, want: true},
		{name: "terminating without workloads", phase: corev1.NamespaceTerminating, hasWorkloads: false, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ns := corev1.Namespace{Status: corev1.NamespaceStatus{Phase: tc.phase}}
			if got := namespaceNeedsAttention(ns, tc.hasWorkloads); got != tc.want {
				t.Fatalf("namespaceNeedsAttention() = %v, want %v", got, tc.want)
			}
		})
	}
}
