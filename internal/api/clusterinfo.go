package api

import (
	"encoding/json"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/log"
	"net/http"
)

type clusterInfo struct {
	infoClient cluster.InfoInterface
	logger     log.Logger
}

type clusterInfoResponse struct {
	Context string `json:"context,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Server  string `json:"server,omitempty"`
	User    string `json:"user,omitempty"`
}

func newClusterInfo(infoClient cluster.InfoInterface, logger log.Logger) clusterInfo {
	return clusterInfo{
		infoClient: infoClient,
		logger:     logger,
	}
}

// ServerHTTP implements http.Handler and returns details about the cluster connection
func (ci clusterInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := clusterInfoResponse{
		Context: ci.infoClient.Context(),
		Cluster: ci.infoClient.Cluster(),
		Server:  ci.infoClient.Server(),
		User:    ci.infoClient.User(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ci.logger.Errorf("encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
