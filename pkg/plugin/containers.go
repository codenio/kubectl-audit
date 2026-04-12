package plugin

import (
	"context"
	"fmt"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// Container is one workload container row for `kubectl audit containers` (not a core API kind).
type Container struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Pod        string `json:"pod"`
	Ready      string `json:"ready"`
	Status     string `json:"status"`
	Restarts   string `json:"restarts"`
	Age        string `json:"age"`
	Ports      string `json:"ports"`
	Image      string `json:"image"`
	PullPolicy string `json:"pullPolicy"`
	Type       string `json:"type"`
}

// DeepCopy returns a deep copy of c.
func (c *Container) DeepCopy() *Container {
	if c == nil {
		return nil
	}
	out := &Container{}
	out.TypeMeta = c.TypeMeta
	c.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Pod = c.Pod
	out.Ready = c.Ready
	out.Status = c.Status
	out.Restarts = c.Restarts
	out.Age = c.Age
	out.Ports = c.Ports
	out.Image = c.Image
	out.PullPolicy = c.PullPolicy
	out.Type = c.Type
	return out
}

// ContainerList is a list of Container rows.
type ContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Container `json:"items"`
}

// DeepCopy returns a deep copy of l.
func (l *ContainerList) DeepCopy() *ContainerList {
	if l == nil {
		return nil
	}
	out := &ContainerList{}
	out.TypeMeta = l.TypeMeta
	l.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = make([]Container, len(l.Items))
	for i := range l.Items {
		if x := l.Items[i].DeepCopy(); x != nil {
			out.Items[i] = *x
		}
	}
	return out
}

// AuditContainers lists container rows that need attention across pods (init + app containers).
// totalInScope counts every container row; benignInScope counts rows that pass the audit bar.
func AuditContainers(configFlags *genericclioptions.ConfigFlags, o AuditOptions) (*ContainerList, int, int, error) {
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

	var (
		totalInScope  int
		benignInScope int
		attention     []Container
	)

	for i := range pods.Items {
		pod := &pods.Items[i]
		// Init containers first (pod spec order), then app containers.
		for _, spec := range pod.Spec.InitContainers {
			cs := findContainerStatus(spec.Name, pod.Status.InitContainerStatuses)
			totalInScope++
			if containerAuditBenign(pod, true, cs) {
				benignInScope++
				continue
			}
			attention = append(attention, buildContainerRow(pod, spec, true, cs))
		}
		for _, spec := range pod.Spec.Containers {
			cs := findContainerStatus(spec.Name, pod.Status.ContainerStatuses)
			totalInScope++
			if containerAuditBenign(pod, false, cs) {
				benignInScope++
				continue
			}
			attention = append(attention, buildContainerRow(pod, spec, false, cs))
		}
	}

	sort.Slice(attention, func(i, j int) bool {
		if attention[i].Namespace != attention[j].Namespace {
			return attention[i].Namespace < attention[j].Namespace
		}
		if attention[i].Pod != attention[j].Pod {
			return attention[i].Pod < attention[j].Pod
		}
		// Init containers before workload containers (lexicographic "container" < "init-container" would invert that).
		rank := func(typ string) int {
			if typ == "init-container" {
				return 0
			}
			return 1
		}
		if ri, rj := rank(attention[i].Type), rank(attention[j].Type); ri != rj {
			return ri < rj
		}
		return attention[i].Name < attention[j].Name
	})

	list := &ContainerList{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ContainerList"},
		ListMeta: metav1.ListMeta{},
		Items:    attention,
	}
	return list, totalInScope, benignInScope, nil
}

// GetObjectKind implements runtime.Object.
func (in *ContainerList) GetObjectKind() schema.ObjectKind { return &in.TypeMeta }

// DeepCopyObject implements runtime.Object.
func (in *ContainerList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	return in.DeepCopy()
}

func findContainerStatus(name string, statuses []corev1.ContainerStatus) *corev1.ContainerStatus {
	for i := range statuses {
		if statuses[i].Name == name {
			return &statuses[i]
		}
	}
	return nil
}

func containerAuditBenign(pod *corev1.Pod, isInit bool, cs *corev1.ContainerStatus) bool {
	if cs == nil {
		return false
	}
	if cs.RestartCount >= podRestartAttentionThreshold {
		return false
	}
	if isInit {
		if t := cs.State.Terminated; t != nil && t.ExitCode == 0 {
			return true
		}
	}
	if w := cs.State.Waiting; w != nil {
		if isTerminalWaitingReason(w.Reason) {
			return false
		}
	}
	if t := cs.State.Terminated; t != nil && !isInit {
		if t.ExitCode != 0 {
			return false
		}
	}
	if !cs.Ready {
		switch pod.Status.Phase {
		case corev1.PodSucceeded, corev1.PodFailed:
			return true
		default:
			return false
		}
	}
	return true
}

func isTerminalWaitingReason(reason string) bool {
	switch reason {
	case "ImagePullBackOff", "ErrImagePull", "CrashLoopBackOff", "ImageInspectError",
		"ErrImageNeverPull", "InvalidImageName", "CreateContainerConfigError",
		"CreateContainerError", "RunContainerError":
		return true
	default:
		return false
	}
}

func buildContainerRow(pod *corev1.Pod, spec corev1.Container, isInit bool, cs *corev1.ContainerStatus) Container {
	ctype := "container"
	if isInit {
		ctype = "init-container"
	}
	row := Container{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Container"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: pod.Namespace,
		},
		Pod:        pod.Name,
		Ready:      readyString(cs),
		Status:     containerStatusString(cs),
		Restarts:   restartsString(cs),
		Age:        containerAgeString(cs),
		Ports:      formatContainerPorts(spec.Ports),
		Image:      spec.Image,
		PullPolicy: string(spec.ImagePullPolicy),
		Type:       ctype,
	}
	return row
}

func readyString(cs *corev1.ContainerStatus) string {
	if cs == nil || !cs.Ready {
		return "0"
	}
	return "1"
}

func containerStatusString(cs *corev1.ContainerStatus) string {
	if cs == nil {
		return "Unknown"
	}
	if cs.State.Running != nil {
		return "Running"
	}
	if w := cs.State.Waiting; w != nil {
		if w.Reason != "" {
			return w.Reason
		}
		return "Waiting"
	}
	if t := cs.State.Terminated; t != nil {
		if t.Reason != "" {
			return t.Reason
		}
		return "Terminated"
	}
	return "Unknown"
}

func restartsString(cs *corev1.ContainerStatus) string {
	if cs == nil {
		return "0"
	}
	if cs.RestartCount == 0 {
		return "0"
	}
	s := fmt.Sprintf("%d", cs.RestartCount)
	if cs.LastTerminationState.Terminated != nil {
		ft := cs.LastTerminationState.Terminated.FinishedAt
		if !ft.IsZero() {
			ago := duration.ShortHumanDuration(time.Since(ft.Time))
			s += fmt.Sprintf(" (%s ago)", ago)
		}
	}
	return s
}

func containerAgeString(cs *corev1.ContainerStatus) string {
	if cs == nil {
		return "-"
	}
	if cs.State.Running != nil && !cs.State.Running.StartedAt.IsZero() {
		return duration.ShortHumanDuration(time.Since(cs.State.Running.StartedAt.Time))
	}
	if t := cs.State.Terminated; t != nil && !t.StartedAt.IsZero() {
		return duration.ShortHumanDuration(time.Since(t.StartedAt.Time))
	}
	return "-"
}

func formatContainerPorts(ports []corev1.ContainerPort) string {
	if len(ports) == 0 {
		return "<none>"
	}
	parts := make([]string, 0, len(ports))
	for _, p := range ports {
		proto := string(p.Protocol)
		if proto == "" {
			proto = string(corev1.ProtocolTCP)
		}
		if p.Name != "" {
			parts = append(parts, fmt.Sprintf("%s:%d/%s", p.Name, p.ContainerPort, proto))
			continue
		}
		parts = append(parts, fmt.Sprintf("%d/%s", p.ContainerPort, proto))
	}
	s := parts[0]
	for i := 1; i < len(parts); i++ {
		s += "," + parts[i]
	}
	return s
}
