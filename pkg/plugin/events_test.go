package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestEventNeedsAttention(t *testing.T) {
	cases := []struct {
		name  string
		event corev1.Event
		want  bool
	}{
		{
			name:  "normal event",
			event: corev1.Event{Type: corev1.EventTypeNormal},
			want:  false,
		},
		{
			name:  "warning event",
			event: corev1.Event{Type: corev1.EventTypeWarning},
			want:  true,
		},
		{
			name:  "empty type",
			event: corev1.Event{},
			want:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := eventNeedsAttention(tc.event); got != tc.want {
				t.Fatalf("eventNeedsAttention() = %v, want %v", got, tc.want)
			}
		})
	}
}
