package fake

import "github.com/twosson/kubeapt/internal/overview"

// SimpleClusterOverview is a fake that implements overview.Interface.
type SimpleClusterOverview struct {
}

var _ overview.Interface = (*SimpleClusterOverview)(nil)

// NewSimpleClusterOverview creates an instance of SimpleClusterOverview.
func NewSimpleClusterOverview() *SimpleClusterOverview {
	return &SimpleClusterOverview{}
}

func (s *SimpleClusterOverview) Namespaces() ([]string, error) {
	names := []string{"default"}
	return names, nil
}

func (s *SimpleClusterOverview) Navigation() (*overview.Navigation, error) {
	return nil, nil
}

func (s *SimpleClusterOverview) Content(path string) error {
	return nil
}
