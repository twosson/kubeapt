package cluster

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/third_party/dynamicfake"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestNamespaceClient_Names(t *testing.T) {
	scheme := runtime.NewScheme()

	dc := dynamicfake.NewSimpleDynamicClient(
		scheme,
		newUnstructured("v1", "Namespace", "", "default"),
		newUnstructured("v1", "Namespace", "", "app-1"),
	)

	nc := newNamespaceClient(dc)

	got, err := nc.Names()
	require.NoError(t, err)

	expected := []string{"default", "app-1"}
	assert.Equal(t, expected, got)
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
