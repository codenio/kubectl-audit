package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func TestServiceHasMatchingPod(t *testing.T) {
	sel := labels.Set(map[string]string{"app": "web"}).AsSelectorPreValidated()
	pods := []corev1.Pod{
		{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "api"}}},
		{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "web", "tier": "fe"}}},
	}
	if serviceHasMatchingPod(pods[:1], sel) {
		t.Fatal("expected no match")
	}
	if !serviceHasMatchingPod(pods, sel) {
		t.Fatal("expected match")
	}
}
