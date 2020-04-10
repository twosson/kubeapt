package overview

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/twosson/kubeapt/internal/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"sync"
)

// StopFunc tells a wwatch to stop watching a namespace.
type StopFunc func()

type Watch2 interface {
	Start() (StopFunc, error)
}

// Watch watches a namespace's objects.
type Watch struct {
	clusterClient cluster.ClientInterface
	cache         Cache
	namespace     string
}

// NewWatch creates an instance of Watch.
func NewWatch(namespace string, clusterClient cluster.ClientInterface, c Cache) *Watch {
	return &Watch{
		clusterClient: clusterClient,
		cache:         c,
		namespace:     namespace,
	}
}

// Start starts the watch. It returns a stop function and an error.
func (w *Watch) Start() (StopFunc, error) {
	resources, err := w.resources()
	if err != nil {
		return nil, err
	}

	var watchers []watch.Interface

	for _, resource := range resources {
		dc, err := w.clusterClient.DynamicClient()
		if err != nil {
			return nil, err
		}

		nri := dc.Resource(resource).Namespace(w.namespace)

		watcher, err := nri.Watch(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		watchers = append(watchers, watcher)
	}

	done := make(chan struct{})

	events, shutdownCh := consumeEvents(done, watchers)

	go func() {
		for event := range events {
			w.eventHandler(event)
		}
	}()

	stopFn := func() {
		done <- struct{}{}
		<-shutdownCh
	}

	return stopFn, nil
}

func (w *Watch) eventHandler(event watch.Event) {
	u, ok := event.Object.(*unstructured.Unstructured)
	if !ok {
		return
	}

	switch t := event.Type; t {
	case watch.Added:
		if err := w.cache.Store(u); err != nil {
			log.Printf("store object: %v", err)
		}
	case watch.Modified:
		if err := w.cache.Store(u); err != nil {
			log.Printf("store object: %v", err)
		}
	case watch.Deleted:
		if err := w.cache.Delete(u); err != nil {
			log.Printf("store object: %v", err)
		}
	case watch.Error:
		log.Printf("unknown log err : %s", spew.Sdump(event))
	default:
		log.Printf("unknown event %q", t)
	}
}

func (w *Watch) resources() ([]schema.GroupVersionResource, error) {
	discoveryClient, err := w.clusterClient.DiscoveryClient()
	if err != nil {
		return nil, err
	}

	lists, err := discoveryClient.ServerResources()
	if err != nil {
		return nil, err
	}

	var gvrs []schema.GroupVersionResource

	for _, list := range lists {
		// todo kubernetes api-resource: the server does not allow this method on the requested resource
		if list.GroupVersion == "metrics.k8s.io/v1beta1" {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, res := range list.APIResources {
			if !res.Namespaced {
				continue
			}
			if hasList(res) {
				gvr := schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: res.Name,
				}

				gvrs = append(gvrs, gvr)
			}
		}
	}

	return gvrs, nil
}

func hasList(res metav1.APIResource) bool {
	for _, verb := range res.Verbs {
		if verb == "list" {
			return true
		}
	}

	return false
}

func consumeEvents(done <-chan struct{}, watchers []watch.Interface) (chan watch.Event, chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(len(watchers))

	events := make(chan watch.Event)

	shutdownComplate := make(chan struct{})

	for _, watcher := range watchers {
		go func(watcher watch.Interface) {
			for event := range watcher.ResultChan() {
				events <- event
			}
		}(watcher)
	}

	go func() {
		<-done
		for _, watch := range watchers {
			watch.Stop()
			wg.Done()
		}

		shutdownComplate <- struct{}{}
	}()

	go func() {
		// wait for all watchers to exit.
		wg.Wait()
		close(events)
	}()

	return events, shutdownComplate
}
