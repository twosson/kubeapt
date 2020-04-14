package api

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster"
	clusterfake "github.com/twosson/kubeapt/internal/cluster/fake"
	"github.com/twosson/kubeapt/internal/log"
	"net/http/httptest"
	"testing"
)

func Test_namespaces_list(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/namespaces", nil)

	tests := []struct {
		name     string
		nsClient cluster.NamespaceInterface
		expected []string
	}{
		{
			name:     "general",
			nsClient: clusterfake.NewNamespaceClient([]string{"default", "other"}, nil, "default"),
			expected: []string{"default", "other"},
		},
		{
			name:     "cannot list due to rbac error",
			nsClient: clusterfake.NewNamespaceClient(nil, errors.New("rbac error"), "initial-namespace"),
			expected: []string{"initial-namespace"},
		},
	}

	for _, tc := range tests {
		handler := newNamespaces(tc.nsClient, log.NopLogger())
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		var nr namespacesResponse
		err := json.NewDecoder(resp.Body).Decode(&nr)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, nr.Namespaces)
	}
}
