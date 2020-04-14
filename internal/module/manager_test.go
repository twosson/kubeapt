package module

import (
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"github.com/twosson/kubeapt/internal/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestManager(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy3"),
	}

	clusterClient, err := fake.NewClient(scheme, nil, objects)
	require.NoError(t, err)

	manager, err := NewManager(clusterClient, "default", log.NopLogger())
	require.NoError(t, err)

	modules := manager.Modules()
	require.NoError(t, err)
	require.Len(t, modules, 1)

	manager.SetNamespace("other")
	manager.Unload()
}

func TestManager_badClient(t *testing.T) {
	var badClient cluster.ClientInterface
	_, err := NewManager(badClient, "default", log.NopLogger())
	require.Error(t, err)
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}
