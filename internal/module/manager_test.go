package module

import (
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestManager(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy3"),
	}

	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	manager := NewManager(clusterClient, "default")

	modules, err := manager.Load()
	require.NoError(t, err)
	require.Len(t, modules, 1)

	manager.SetNamespace("other")
	manager.Unload()
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"namespace": namespace,
			"name":      name,
		},
	}}
}
