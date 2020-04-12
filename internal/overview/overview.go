package overview

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/log"
	"net/http"
	"os"
	"sync"
)

// ClusterOverview is an API for generating a cluster overview.
type ClusterOverview struct {
	client       cluster.ClientInterface
	mu           sync.Mutex
	namespace    string
	logger       log.Logger
	watchFactory func(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch
	cache        Cache
	stopFn       func()
	generator    *realGenerator
}

// NewClusterOverview creates an instance of ClusterOverview.
func NewClusterOverview(client cluster.ClientInterface, namespace string, logger log.Logger) *ClusterOverview {
	var opts []MemoryCacheOpt

	if os.Getenv("DASH_VERBOSE_CACHE") != "" {
		ch := make(chan CacheNotification)

		go func() {
			for notif := range ch {
				spew.Dump(notif)
			}
		}()

		opts = append(opts, CacheNotificationOpt(ch, nil))
	}

	cache := NewMemoryCache(opts...)

	var pathFilters []pathFilter
	pathFilters = append(pathFilters, rootDescriber.PathFilters()...)
	pathFilters = append(pathFilters, eventsDescriber.PathFilters()...)

	g := newGenerator(cache, pathFilters, client)

	co := &ClusterOverview{
		namespace: namespace,
		client:    client,
		logger:    logger,
		cache:     cache,
		generator: g,
	}

	co.watchFactory = co.defaultWatchFactory
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
	c.logger.With("namespace", namespace, "module", "overview").Debugf("stopping")
	c.Stop()

	c.logger.With("namespace", namespace, "module", "overview").Debugf("setting namespace")
	c.namespace = namespace
	return c.Start()
}

// Start starts overview.
func (c *ClusterOverview) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.namespace == "" {
		return nil
	}

	if c.stopFn != nil {
		return errors.New("synchronization error - residual state detected")
	}

	stopFn, err := c.watch(c.namespace)
	if err != nil {
		return err
	}

	c.stopFn = stopFn

	return nil
}

// Stop stops overview.
func (c *ClusterOverview) Stop() {
	c.mu.Lock()
	stopFn := c.stopFn
	c.stopFn = nil
	c.mu.Unlock()

	if stopFn != nil {
		go func() {
			stopFn()
		}()
	}
}

func (c *ClusterOverview) watch(namespace string) (StopFunc, error) {
	c.logger.With("namespace", namespace, "module", "overview").Debugf("watching namespace")

	watch := c.watchFactory(namespace, c.client, c.cache)
	return watch.Start()
}

func (c *ClusterOverview) defaultWatchFactory(namespace string, clusterClient cluster.ClientInterface, cache Cache) Watch {
	return NewWatch(namespace, clusterClient, cache, c.logger)
}
