package overview

import (
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"log"
	"net/http"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client *cluster.Cluster
	stopFn func()
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client *cluster.Cluster) *ClusterOverview {
	return &ClusterOverview{client: client}
}

// ContentPath returns the content path for overview.
func (c *ClusterOverview) ContentPath() string {
	return "/overview"
}

// Handler returns a handler for serving overview HTTP content.
func (c *ClusterOverview) Handler(prefix string) http.Handler {
	return newHandler(prefix)
}

func (c *ClusterOverview) Namespaces() ([]string, error) {
	nsClient, err := c.client.NamespaceClient()
	if err != nil {
		return nil, err
	}

	return nsClient.Names()
}

func (c *ClusterOverview) Navigation(root string) (*apt.Navigation, error) {
	return navigationEntries(root)
}

func (c *ClusterOverview) Content() error {
	return nil
}

// Start starts overview.
func (c *ClusterOverview) Start() error {
	log.Printf("Starting cluster overview")
	return nil
}

// Stop stops overview.
func (c *ClusterOverview) Stop() {
	if c.stopFn != nil {
		log.Printf("Stopping cluster overview")
		c.stopFn()
	}
}
