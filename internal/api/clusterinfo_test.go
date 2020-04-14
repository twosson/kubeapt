package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"github.com/twosson/kubeapt/internal/log"
	"net/http/httptest"
	"testing"
)

func Test_clusterInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/info", nil)

	tests := []struct {
		name       string
		infoClient cluster.InfoInterface
		expected   clusterInfoResponse
	}{
		{
			name: "general",
			infoClient: fake.ClusterInfo{
				ContextVal: "main-context",
				ClusterVal: "my-cluster",
				ServerVal:  "https://localhost:6443",
				UserVal:    "me-of-course",
			},
			expected: clusterInfoResponse{
				Context: "main-context",
				Cluster: "my-cluster",
				Server:  "https://localhost:6443",
				User:    "me-of-course",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

		})
		handler := newClusterInfo(tc.infoClient, log.NopLogger())
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		var ciResp clusterInfoResponse
		err := json.NewDecoder(resp.Body).Decode(&ciResp)
		require.NoError(t, err)

		assert.Equal(t, tc.expected, ciResp)
	}
}
