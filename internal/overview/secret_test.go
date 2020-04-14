package overview

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/util/clock"
	"testing"
	"time"
)

func TestSecretData_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewSecretData("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestSecretData(t *testing.T) {
	v := NewSecretData("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()
	cache := NewMemoryCache()

	secret := loadFromFile(t, "secret-1.yaml")
	secret = convertToInternal(t, secret)

	got, err := v.Content(ctx, secret, cache)
	require.NoError(t, err)

	dataSection := content.NewSection()
	dataSection.AddText("ca.crt", "1025 bytes")
	dataSection.AddText("namespace", "8 bytes")
	dataSection.AddText("token", "token")

	dataSummary := content.NewSummary("Data", []content.Section{dataSection})

	expected := []content.Content{
		&dataSummary,
	}

	assert.Equal(t, got, expected)
}
