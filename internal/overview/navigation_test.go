package overview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNavigationEntries(t *testing.T) {
	got, err := navigationEntries()
	require.NoError(t, err)

	assert.Equal(t, got.Title, "Overview")
}
