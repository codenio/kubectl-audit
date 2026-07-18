package plugin

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestAuditNamespacedResource(t *testing.T) {
	if !auditNamespacedResource("pods") {
		t.Fatal("expected pods to be namespaced")
	}
	if auditNamespacedResource("namespaces") {
		t.Fatal("expected namespaces to be cluster-scoped")
	}
}

func TestObjectKeysForFilter(t *testing.T) {
	cases := []struct {
		name       string
		obj        interface{}
		resource   string
		wantKeys   []string
		namespaced bool
		empty      bool
	}{
		{
			name:       "empty pod list",
			obj:        &corev1.PodList{},
			resource:   "pods",
			empty:      true,
			namespaced: true,
		},
		{
			name: "pod keys",
			obj: &corev1.PodList{Items: []corev1.Pod{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "a", Name: "p1"}},
			}},
			resource:   "pods",
			wantKeys:   []string{"a/p1"},
			namespaced: true,
		},
		{
			name: "node keys",
			obj: &corev1.NodeList{Items: []corev1.Node{
				{ObjectMeta: metav1.ObjectMeta{Name: "node-1"}},
			}},
			resource: "nodes",
			wantKeys: []string{"node-1"},
		},
		{
			name: "namespace keys",
			obj: &corev1.NamespaceList{Items: []corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "demo"}},
			}},
			resource: "namespaces",
			wantKeys: []string{"demo"},
		},
		{
			name: "deployment keys",
			obj: &appsv1.DeploymentList{Items: []appsv1.Deployment{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "ops", Name: "api"}},
			}},
			resource:   "deployments",
			wantKeys:   []string{"ops/api"},
			namespaced: true,
		},
		{
			name: "job keys",
			obj: &batchv1.JobList{Items: []batchv1.Job{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "batch", Name: "import"}},
			}},
			resource:   "jobs",
			wantKeys:   []string{"batch/import"},
			namespaced: true,
		},
		{
			name: "cronjob keys",
			obj: &batchv1.CronJobList{Items: []batchv1.CronJob{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "ops", Name: "backup"}},
			}},
			resource:   "cronjobs",
			wantKeys:   []string{"ops/backup"},
			namespaced: true,
		},
		{
			name: "service keys",
			obj: &corev1.ServiceList{Items: []corev1.Service{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web"}},
			}},
			resource:   "services",
			wantKeys:   []string{"default/web"},
			namespaced: true,
		},
		{
			name: "event keys",
			obj: &corev1.EventList{Items: []corev1.Event{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "ev-1"}},
			}},
			resource:   "events",
			wantKeys:   []string{"default/ev-1"},
			namespaced: true,
		},
		{
			name:     "pv keys",
			obj:      &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{{ObjectMeta: metav1.ObjectMeta{Name: "pv-1"}}}},
			resource: "pv",
			wantKeys: []string{"pv-1"},
		},
		{
			name: "pvc keys",
			obj: &corev1.PersistentVolumeClaimList{Items: []corev1.PersistentVolumeClaim{
				{ObjectMeta: metav1.ObjectMeta{Namespace: "data", Name: "claim-1"}},
			}},
			resource:   "pvc",
			wantKeys:   []string{"data/claim-1"},
			namespaced: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keys, namespaced, empty := objectKeysForFilter(tc.obj.(runtime.Object), tc.resource)
			if namespaced != tc.namespaced {
				t.Fatalf("namespaced = %v, want %v", namespaced, tc.namespaced)
			}
			if empty != tc.empty {
				t.Fatalf("empty = %v, want %v", empty, tc.empty)
			}
			for _, want := range tc.wantKeys {
				if _, ok := keys[want]; !ok {
					t.Fatalf("missing key %q in %#v", want, keys)
				}
			}
		})
	}
}

func TestTableRowKey(t *testing.T) {
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "a", Name: "p1"}}
	row := metav1.TableRow{Object: runtime.RawExtension{Object: pod}}
	if got := tableRowKey(&row, true); got != "a/p1" {
		t.Fatalf("got %q", got)
	}

	raw, err := json.Marshal(pod)
	if err != nil {
		t.Fatal(err)
	}
	rawRow := metav1.TableRow{Object: runtime.RawExtension{Raw: raw}}
	if got := tableRowKey(&rawRow, true); got != "a/p1" {
		t.Fatalf("got %q from raw", got)
	}
}

func TestNamespaceFromRowRaw(t *testing.T) {
	if got := namespaceFromRowRaw(nil); got != "" {
		t.Fatalf("got %q", got)
	}
	raw, err := json.Marshal(map[string]interface{}{
		"metadata": map[string]string{"namespace": "demo"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := namespaceFromRowRaw(raw); got != "demo" {
		t.Fatalf("got %q", got)
	}
}

func TestPrependNamespaceColumn(t *testing.T) {
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "demo", Name: "web"}}
	raw, err := json.Marshal(pod)
	if err != nil {
		t.Fatal(err)
	}
	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{{Name: "Name"}},
		Rows: []metav1.TableRow{
			{Cells: []interface{}{"web"}, Object: runtime.RawExtension{Raw: raw}},
		},
	}
	prependNamespaceColumn(table)
	if table.ColumnDefinitions[0].Name != "Namespace" {
		t.Fatalf("first column = %q", table.ColumnDefinitions[0].Name)
	}
	if table.Rows[0].Cells[0] != "demo" {
		t.Fatalf("namespace cell = %v", table.Rows[0].Cells[0])
	}
}

func TestEnsureNamespaceColumn(t *testing.T) {
	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{{Name: "Namespace", Priority: 1}},
	}
	ensureNamespaceColumn(table)
	if table.ColumnDefinitions[0].Priority != 0 {
		t.Fatalf("priority = %d", table.ColumnDefinitions[0].Priority)
	}
}

func TestTableAcceptHeader(t *testing.T) {
	got := tableAcceptHeader()
	if got == "" || !strings.Contains(got, "Table") || !strings.Contains(got, "application/json") {
		t.Fatalf("unexpected header: %q", got)
	}
}

func TestNamespacesWithWorkloads(t *testing.T) {
	cs := fake.NewSimpleClientset(
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "team-a", Name: "p1"}},
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "team-b", Name: "api"}},
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{Namespace: "team-c", Name: "import"}},
	)
	got, err := namespacesWithWorkloads(context.Background(), cs)
	if err != nil {
		t.Fatal(err)
	}
	for _, ns := range []string{"team-a", "team-b", "team-c"} {
		if _, ok := got[ns]; !ok {
			t.Fatalf("missing namespace %q in %#v", ns, got)
		}
	}
}
