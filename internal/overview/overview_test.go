package overview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestClusterOverview(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy"),
	}

	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	o := NewClusterOverview(clusterClient, "default")

	namespaces, err := o.Namespaces()
	require.NoError(t, err)

	assert.Equal(t, "overview", o.Name())
	assert.Equal(t, "/overview", o.ContentPath())
	assert.Equal(t, []string{"default"}, namespaces)
}

func TestClusterOverview_SetNamespace(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy"),
	}

	clusterClient, err := fake.NewClient(scheme, objects)
	require.NoError(t, err)

	o := NewClusterOverview(clusterClient, "default")
	defer o.Stop()

	err = o.SetNamespace("ns2")
	require.NoError(t, err)
}
