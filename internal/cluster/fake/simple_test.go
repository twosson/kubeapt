package fake

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleClusterOverview_Namespaces(t *testing.T) {
	seo := NewSimpleClusterOverview()
	got, err := seo.Namespaces()
	require.NoError(t, err)

	expected := []string{"default"}
	assert.Equal(t, expected, got)
}

func TestSimpleClusterOverview_Navigation(t *testing.T) {
	seo := NewSimpleClusterOverview()
	err := seo.Navigation()
	require.NoError(t, err)
}

func TestSimpleClusterOverview_Content(t *testing.T) {
	seo := NewSimpleClusterOverview()
	err := seo.Content("/path")
	require.NoError(t, err)
}
