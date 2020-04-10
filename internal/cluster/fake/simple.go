package fake

import (
	"github.com/twosson/kubeapt/internal/apt"
	"net/http"
)

// SimpleClusterOverview is a fake that implements overview.Interface.
type SimpleClusterOverview struct{}

// NewSimpleClusterOverview creates an instance of SimpleClusterOverview.
func NewSimpleClusterOverview() *SimpleClusterOverview {
	return &SimpleClusterOverview{}
}

// Name is the module name.
func (sco *SimpleClusterOverview) Name() string {
	return "overview"
}

// Handler returns a nil HTTP handler.
func (sco *SimpleClusterOverview) Handler(prefix string) http.Handler {
	return nil
}

// ContentPath returns the content path for mounting this module.
func (sco *SimpleClusterOverview) ContentPath() string {
	return "/overview"
}

// Navigation is a no-op.
func (sco *SimpleClusterOverview) Navigation(root string) (*apt.Navigation, error) {
	return nil, nil
}

// SetNamespace sets the namespace for this module. It is a no-op
func (sco *SimpleClusterOverview) SetNamespace(namespace string) error {
	return nil
}

// Starts starts the module. It ia a no-op
func (sco *SimpleClusterOverview) Start() error {
	return nil
}

// Stop stops the module. It is a no-op
func (sco *SimpleClusterOverview) Stop() {
}
