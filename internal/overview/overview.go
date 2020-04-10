package overview

import (
	"fmt"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"log"
	"net/http"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client       cluster.ClientInterface
	namespace    string
	watchFactory func(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch2
	cache        Cache
	stopFn       func()
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client cluster.ClientInterface, namespace string) *ClusterOverview {
	return &ClusterOverview{
		client:       client,
		namespace:    namespace,
		watchFactory: watchFactory,
	}
}

// Name returns the name for this module.
func (c *ClusterOverview) Name() string {
	return "overview"
}

// ContentPath returns the content path for overview.
func (c *ClusterOverview) ContentPath() string {
	return fmt.Sprintf("/%s", c.Name())
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

// SetNamespace sets the current namespace.
func (c *ClusterOverview) SetNamespace(namespace string) error {
	log.Printf("Setting namespace for overview to %q", namespace)
	if c.stopFn != nil {
		c.stopFn()
	}

	c.namespace = namespace
	return c.Start()
}

// Start starts overview.
func (c *ClusterOverview) Start() error {
	if c.namespace == "" {
		return nil
	}

	log.Printf("Starting cluster overview")

	stopFn, err := c.watch(c.namespace)
	if err != nil {
		return err
	}

	c.stopFn = stopFn

	return nil
}

// Stop stops overview.
func (c *ClusterOverview) Stop() {
	if c.stopFn != nil {
		log.Printf("Stopping cluster overview")
		c.stopFn()
	}
}

func (c *ClusterOverview) watch(namespace string) (StopFunc, error) {
	log.Printf("Watching namespace %s", namespace)

	cache := NewMemoryCache()
	c.cache = cache

	watch := c.watchFactory(namespace, c.client, c.cache)
	return watch.Start()
}

func watchFactory(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch2 {
	return NewWatch(namespace, clusterClient, cache)
}
