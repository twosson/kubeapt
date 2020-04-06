package fake

import (
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"net/http"
)

// SimpleClusterOverview is a fake that implements overview.Interface.
type SimpleClusterOverview struct{}

// NewSimpleClusterOverview creates an instance of SimpleClusterOverview.
func NewSimpleClusterOverview() *SimpleClusterOverview {
	return &SimpleClusterOverview{}
}

func (sco *SimpleClusterOverview) Handler(prefix string) http.Handler {
	return nil
}

func (sco *SimpleClusterOverview) ContentPath() string {
	return "/overview"
}

// Navigation is a no-op.
func (sco *SimpleClusterOverview) Navigation(root string) (*apt.Navigation, error) {
	return nil, nil
}

func (sco *SimpleClusterOverview) Start() error {
	return nil
}

func (sco *SimpleClusterOverview) Stop() {
}

// NamespaceClient is a fake that implements cluster.NamespaceInterface.
type NamespaceClient struct{}

var _ cluster.NamespaceInterface = (*NamespaceClient)(nil)

// NewNamespaceClient creates an instance of NamespaceClient.
func NewNamespaceClient() *NamespaceClient {
	return &NamespaceClient{}
}

func (nc *NamespaceClient) Names() ([]string, error) {
	names := []string{"default"}
	return names, nil
}
