package overview

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/core"
	"testing"
	"time"
)

func TestConfigMapDetails_InvalidObject(t *testing.T) {
	cm := NewConfigMapDetails("prefix", "ns", clock.NewFakeClock(time.Now()))
	ctx := context.Background()

	object := &unstructured.Unstructured{}

	_, err := cm.Content(ctx, object, nil)
	require.Error(t, err)
}

func TestConfigMapDetails(t *testing.T) {
	cm := NewConfigMapDetails("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()
	object := &core.ConfigMap{
		Data: map[string]string{
			"test": "data",
		},
	}

	contents, err := cm.Content(ctx, object, nil)
	require.NoError(t, err)

	require.Len(t, contents, 1)

	table, ok := contents[0].(*content.Table)
	require.True(t, ok)
	require.Len(t, table.Rows, 1)

	expectedColumns := []string{"Key", "Value"}
	assert.Equal(t, expectedColumns, table.ColumnNames())

	expectedRow := content.TableRow{
		"Key":   content.NewStringText("test"),
		"Value": content.NewStringText("data"),
	}
	assert.Equal(t, expectedRow, table.Rows[0])
}
