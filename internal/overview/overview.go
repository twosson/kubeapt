package overview

import (
	"fmt"
	"github.com/twosson/kubeapt/internal/cluster"
)

type notImplemented struct {
	name string
}

func (e *notImplemented) Error() string {
	return fmt.Sprintf("%s not implemented", e.name)
}

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client *cluster.Cluster
}

var _ Interface = (*ClusterOverview)(nil)

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client *cluster.Cluster) *ClusterOverview {
	return &ClusterOverview{client: client}
}

func (c *ClusterOverview) Namespaces() ([]string, error) {
	return c.client.Namespace.Names()
}

func (c *ClusterOverview) Navigation() error {
	return &notImplemented{name: "Navigation"}
}

func (c *ClusterOverview) Content(path string) error {
	return &notImplemented{name: "Content"}
}
