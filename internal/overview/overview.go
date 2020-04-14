package overview

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/log"
	"k8s.io/client-go/restmapper"
	"net/http"
	"os"
	"sync"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client    cluster.ClientInterface
	mu        sync.Mutex
	namespace string
	logger    log.Logger
	cache     Cache
	stopCh    chan struct{}
	generator *realGenerator
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client cluster.ClientInterface, namespace string, logger log.Logger) *ClusterOverview {
	stopCh := make(chan struct{})

	var opts []InformerCacheOpt

	if os.Getenv("DASH_VERBOSE_CACHE") != "" {
		ch := make(chan CacheNotification)

		go func() {
			for notif := range ch {
				spew.Dump(notif)
			}
		}()

		opts = append(opts, InformerCacheNotificationOpt(ch, stopCh))
	}

	dynamicClient, err := client.DynamicClient()
	if err != nil {
		// TODO error handling
		return nil
	}
	di, err := client.DiscoveryClient()
	if err != nil {
		// TODO error handling
		return nil
	}

	groupResources, err := restmapper.GetAPIGroupResources(di)
	if err != nil {
		logger.Errorf("discovering APIGroupResources: %v", err)
		// TODO error handling
		return nil
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	opts = append(opts, InformerCacheLoggerOpt(logger))
	cache := NewInformerCache(stopCh, dynamicClient, rm, opts...)

	var pathFilters []pathFilter
	pathFilters = append(pathFilters, rootDescriber.PathFilters(namespace)...)
	pathFilters = append(pathFilters, eventsDescriber.PathFilters(namespace)...)

	g := newGenerator(cache, pathFilters, client)

	co := &ClusterOverview{
		namespace: namespace,
		client:    client,
		logger:    logger,
		cache:     cache,
		generator: g,
		stopCh:    stopCh,
	}
	return co
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
	return newHandler(prefix, c.generator, stream, c.logger)
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
	c.logger.With("namespace", namespace, "module", "overview").Debugf("setting namespace")
	c.namespace = namespace
	return nil
}

// Start starts overview.
func (c *ClusterOverview) Start() error {
	return nil
}

// Stop stops overview.
func (c *ClusterOverview) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.stopCh)
	c.stopCh = nil
}
