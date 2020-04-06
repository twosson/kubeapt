package fake

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleClusterOverview_Navigation(t *testing.T) {
	sco := NewSimpleClusterOverview()
	_, err := sco.Navigation("/root")
	require.NoError(t, err)
}
