package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestContainerAuditBenign(t *testing.T) {
	cases := []struct {
		name   string
		pod    corev1.Pod
		isInit bool
		cs     *corev1.ContainerStatus
		want   bool
	}{
		{
			name: "nil status",
			pod:  corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}},
			cs:   nil,
			want: false,
		},
		{
			name: "running ready",
			pod:  corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}},
			cs: &corev1.ContainerStatus{
				Name:  "c",
				Ready: true,
				State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
			},
			want: true,
		},
		{
			name: "image pull back off",
			pod:  corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodPending}},
			cs: &corev1.ContainerStatus{
				Name:  "c",
				Ready: false,
				State: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{Reason: "ImagePullBackOff"},
				},
			},
			want: false,
		},
		{
			name:   "init completed",
			pod:    corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}},
			isInit: true,
			cs: &corev1.ContainerStatus{
				Name:  "init",
				Ready: false,
				State: corev1.ContainerState{
					Terminated: &corev1.ContainerStateTerminated{ExitCode: 0},
				},
			},
			want: true,
		},
		{
			name: "high restarts",
			pod:  corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}},
			cs: &corev1.ContainerStatus{
				Name:         "c",
				Ready:        true,
				RestartCount: podRestartAttentionThreshold,
				State:        corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
			},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := containerAuditBenign(&tc.pod, tc.isInit, tc.cs); got != tc.want {
				t.Fatalf("containerAuditBenign() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFormatContainerPorts(t *testing.T) {
	if got := formatContainerPorts(nil); got != "<none>" {
		t.Fatalf("got %q", got)
	}
	ports := []corev1.ContainerPort{
		{ContainerPort: 7654, Protocol: corev1.ProtocolTCP},
		{ContainerPort: 9852, Protocol: corev1.ProtocolTCP},
	}
	got := formatContainerPorts(ports)
	if got != "7654/TCP,9852/TCP" {
		t.Fatalf("got %q", got)
	}
}
