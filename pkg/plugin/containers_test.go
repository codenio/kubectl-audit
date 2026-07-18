package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		{Name: "metrics", ContainerPort: 7654, Protocol: corev1.ProtocolTCP},
		{ContainerPort: 9852, Protocol: corev1.ProtocolTCP},
	}
	got := formatContainerPorts(ports)
	if got != "metrics:7654/TCP,9852/TCP" {
		t.Fatalf("got %q", got)
	}
}

func TestFindContainerStatus(t *testing.T) {
	statuses := []corev1.ContainerStatus{
		{Name: "app", Ready: true},
		{Name: "sidecar", Ready: false},
	}
	if got := findContainerStatus("missing", statuses); got != nil {
		t.Fatal("expected nil")
	}
	if got := findContainerStatus("sidecar", statuses); got == nil || got.Name != "sidecar" {
		t.Fatalf("unexpected status: %#v", got)
	}
}

func TestBuildContainerRow(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "web-0", Namespace: "demo"},
	}
	spec := corev1.Container{
		Name:            "app",
		Image:           "nginx:1.27",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports:           []corev1.ContainerPort{{ContainerPort: 8080, Protocol: corev1.ProtocolTCP}},
	}
	cs := &corev1.ContainerStatus{
		Name:  "app",
		Ready: true,
		State: corev1.ContainerState{Running: &corev1.ContainerStateRunning{}},
	}
	row := buildContainerRow(pod, spec, false, cs)
	if row.Pod != "web-0" || row.Namespace != "demo" || row.Type != "container" {
		t.Fatalf("unexpected row: %#v", row)
	}
	if row.Ready != "1" || row.Status != "Running" {
		t.Fatalf("unexpected ready/status: %s / %s", row.Ready, row.Status)
	}
}

func TestIsTerminalWaitingReason(t *testing.T) {
	if !isTerminalWaitingReason("CrashLoopBackOff") {
		t.Fatal("expected terminal reason")
	}
	if isTerminalWaitingReason("ContainerCreating") {
		t.Fatal("expected non-terminal reason")
	}
}

func TestContainerDisplayHelpers(t *testing.T) {
	if readyString(nil) != "0" {
		t.Fatal("expected not ready")
	}
	if readyString(&corev1.ContainerStatus{Ready: true}) != "1" {
		t.Fatal("expected ready")
	}
	if containerStatusString(nil) != "Unknown" {
		t.Fatalf("got %q", containerStatusString(nil))
	}
	if got := containerStatusString(&corev1.ContainerStatus{
		State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"}},
	}); got != "CrashLoopBackOff" {
		t.Fatalf("got %q", got)
	}
	if restartsString(nil) != "0" {
		t.Fatalf("got %q", restartsString(nil))
	}
	if containerAgeString(nil) != "-" {
		t.Fatalf("got %q", containerAgeString(nil))
	}
}

func TestContainerDeepCopy(t *testing.T) {
	list := &ContainerList{
		Items: []Container{{ObjectMeta: metav1.ObjectMeta{Name: "app"}, Pod: "web-0"}},
	}
	copy := list.DeepCopy()
	if copy == list || copy.Items[0].Pod != "web-0" {
		t.Fatal("unexpected deep copy")
	}
	copy.Items[0].Pod = "changed"
	if list.Items[0].Pod != "web-0" {
		t.Fatal("expected independent copy")
	}
}
