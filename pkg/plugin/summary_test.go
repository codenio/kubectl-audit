package plugin

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestSummarizeAudit(t *testing.T) {
	list := &corev1.PodList{
		Items: make([]corev1.Pod, 3),
	}
	total, benign, attention := SummarizeAudit(10, 7, list)
	if total != 10 || benign != 7 || attention != 3 {
		t.Fatalf("got total=%d benign=%d attention=%d", total, benign, attention)
	}
}

func TestSummarizeAudit_EmptyAttention(t *testing.T) {
	list := &corev1.PodList{}
	total, benign, attention := SummarizeAudit(5, 4, list)
	if total != 5 || benign != 4 || attention != 0 {
		t.Fatalf("got total=%d benign=%d attention=%d", total, benign, attention)
	}
}

func TestSummarizeAudit_AllHealthyNoAttentionRows(t *testing.T) {
	list := &corev1.PodList{}
	total, benign, attention := SummarizeAudit(1, 1, list)
	if total != 1 || benign != 1 || attention != 0 {
		t.Fatalf("got total=%d benign=%d attention=%d", total, benign, attention)
	}
}

func TestSummarizeAudit_InScopeExceedsBenignEmptyTable(t *testing.T) {
	list := &corev1.PodList{}
	total, benign, attention := SummarizeAudit(2, 1, list)
	if total != 2 || benign != 1 || attention != 0 {
		t.Fatalf("got total=%d benign=%d attention=%d", total, benign, attention)
	}
}

func TestSummarizeAudit_EmptyNamespace(t *testing.T) {
	list := &corev1.PodList{}
	total, benign, attention := SummarizeAudit(0, 0, list)
	if total != 0 || benign != 0 || attention != 0 {
		t.Fatalf("got total=%d benign=%d attention=%d", total, benign, attention)
	}
}
