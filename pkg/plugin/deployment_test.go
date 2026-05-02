package plugin

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
)

func TestDeploymentScaledToZero(t *testing.T) {
	z := int32(0)
	two := int32(2)
	cases := []struct {
		name string
		d    appsv1.Deployment
		want bool
	}{
		{name: "nil replicas", d: appsv1.Deployment{}, want: false},
		{name: "zero", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &z}}, want: true},
		{name: "non-zero", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := deploymentScaledToZero(tc.d); got != tc.want {
				t.Fatalf("deploymentScaledToZero() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDeploymentDesiredReplicas(t *testing.T) {
	z, two := int32(0), int32(2)
	cases := []struct {
		name string
		d    appsv1.Deployment
		want int32
	}{
		{name: "nil spec replicas defaults to 1", d: appsv1.Deployment{}, want: 1},
		{name: "explicit zero", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &z}}, want: 0},
		{name: "explicit two", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}}, want: 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := deploymentDesiredReplicas(tc.d); got != tc.want {
				t.Fatalf("deploymentDesiredReplicas() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDeploymentReadyBelowDesired(t *testing.T) {
	z, two, three := int32(0), int32(2), int32(3)
	cases := []struct {
		name string
		d    appsv1.Deployment
		want bool
	}{
		{name: "desired zero not under", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &z}, Status: appsv1.DeploymentStatus{ReadyReplicas: 0}}, want: false},
		{name: "two desired two ready", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}, Status: appsv1.DeploymentStatus{ReadyReplicas: 2}}, want: false},
		{name: "two desired one ready", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}, Status: appsv1.DeploymentStatus{ReadyReplicas: 1}}, want: true},
		{name: "three desired zero ready", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &three}, Status: appsv1.DeploymentStatus{ReadyReplicas: 0}}, want: true},
		{name: "nil replicas means desired 1 ready 0", d: appsv1.Deployment{Status: appsv1.DeploymentStatus{ReadyReplicas: 0}}, want: true},
		{name: "nil replicas ready 1", d: appsv1.Deployment{Status: appsv1.DeploymentStatus{ReadyReplicas: 1}}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := deploymentReadyBelowDesired(tc.d); got != tc.want {
				t.Fatalf("deploymentReadyBelowDesired() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDeploymentNeedsAttention(t *testing.T) {
	z, two := int32(0), int32(2)
	cases := []struct {
		name string
		d    appsv1.Deployment
		want bool
	}{
		{name: "scaled to zero", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &z}}, want: true},
		{name: "under ready", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}, Status: appsv1.DeploymentStatus{ReadyReplicas: 1}}, want: true},
		{name: "healthy two of two", d: appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: &two}, Status: appsv1.DeploymentStatus{ReadyReplicas: 2}}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := deploymentNeedsAttention(tc.d); got != tc.want {
				t.Fatalf("deploymentNeedsAttention() = %v, want %v", got, tc.want)
			}
		})
	}
}
