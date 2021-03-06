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

func TestIngressSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestIngressDetails_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressDetails("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestIngressDetails(t *testing.T) {
	v := NewIngressDetails("prefix", "ns", clock.NewFakeClock(time.Now()))

	cache := NewMemoryCache()

	ingress := loadFromFile(t, "ingress-1.yaml")
	ingress = convertToInternal(t, ingress)

	ctx := context.Background()

	got, err := v.Content(ctx, ingress, cache)
	require.NoError(t, err)

	tlsTable := content.NewTable("TLS", "TLS is not configured for this Ingress")
	tlsTable.Columns = tableCols("Secret", "Hosts")

	rulesTable := content.NewTable("Rules", "Rules are not configured for this Ingress")
	rulesTable.Columns = tableCols("Host", "Path", "Backend")
	rulesTable.AddRow(content.TableRow{
		"Backend": content.NewLinkText("test:80", "/content/overview/discovery-and-load-balancing/services/test"),
		"Host":    content.NewStringText(""),
		"Path":    content.NewStringText("/testpath"),
	})

	expected := []content.Content{
		&tlsTable,
		&rulesTable,
	}

	assert.Equal(t, expected, got)
}
