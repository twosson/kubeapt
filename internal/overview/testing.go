package overview

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertViewInvalidObject(t *testing.T, v View) {
	ctx := context.Background()
	_, err := v.Content(ctx, nil, nil)
	require.Error(t, err)
}
