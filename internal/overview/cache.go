package overview

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
	"sync"
)

// Cache stores Kubernetes objects.
type Cache interface {
	Store(obj *unstructured.Unstructured) error
	Retrieve(key CacheKey) ([]*unstructured.Unstructured, error)
	Delete(obj *unstructured.Unstructured) error
}

// CacheKey is a key for the cache.
type CacheKey struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
}

// MemoryCacheOpt is an option for configuring memory cache.
type MemoryCacheOpt func(*MemoryCache)

// CacheAction is a cache action.
type CacheAction string

const (
	// CacheStore is a store action.
	CacheStore CacheAction = "store"
	// CacheDelete is a delete action.
	CacheDelete CacheAction = "delete"
)

// CacheNotification is a notifcation for a cache.
type CacheNotification struct {
	CacheKey CacheKey
	Action   CacheAction
}

// CacheNotificationOpt sets a channel that will receive a notification
// every time cache performs an add/delete.
func CacheNotificationOpt(ch chan<- CacheNotification) MemoryCacheOpt {
	return func(c *MemoryCache) {
		c.notifyCh = ch
	}
}

// MemoryCache stores a cache of Kubernetes objects in memory.
type MemoryCache struct {
	store    map[CacheKey]*unstructured.Unstructured
	mu       sync.Mutex
	notifyCh chan<- CacheNotification
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates on instance of MemoryCache.
func NewMemoryCache(opts ...MemoryCacheOpt) *MemoryCache {
	mc := &MemoryCache{
		store: make(map[CacheKey]*unstructured.Unstructured),
	}

	for _, opt := range opts {
		opt(mc)
	}

	return mc
}

// Reset resets the cache.
func (m *MemoryCache) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k := range m.store {
		delete(m.store, k)
	}
}

// Store stores an object to the object.
func (m *MemoryCache) Store(obj *unstructured.Unstructured) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := CacheKey{
		Namespace:  obj.GetNamespace(),
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Name:       obj.GetName(),
	}

	log.Printf("cache: store %+v", key)

	m.store[key] = obj
	m.notify(CacheStore, key)

	return nil
}

// Retrieve retrieves an object from the cache.
func (m *MemoryCache) Retrieve(key CacheKey) ([]*unstructured.Unstructured, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var objs []*unstructured.Unstructured

	for k, v := range m.store {
		if k.Namespace != key.Namespace {
			continue
		}

		if key.APIVersion == "" {
			objs = append(objs, v)
			continue
		}

		if k.APIVersion == key.APIVersion {
			if key.Kind == "" {
				objs = append(objs, v)
				continue
			}

			if k.Kind == key.Kind {
				if key.Name == "" {
					objs = append(objs, v)
					continue
				}

				if k.Name == key.Name {
					objs = append(objs, v)
				}
			}
		}
	}

	return objs, nil
}

// Delete deletes an object from the cache.
func (m *MemoryCache) Delete(obj *unstructured.Unstructured) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	namespace := obj.GetNamespace()
	apiVersion := obj.GetAPIVersion()
	kind := obj.GetKind()
	name := obj.GetName()

	key := CacheKey{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	delete(m.store, key)

	log.Printf("cache: delete %+v", key)
	m.notify(CacheDelete, key)

	return nil
}

func (m *MemoryCache) notify(action CacheAction, key CacheKey) {
	if m.notifyCh == nil {
		return
	}

	m.notifyCh <- CacheNotification{Action: action, CacheKey: key}
}
