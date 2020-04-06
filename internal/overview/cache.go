package overview

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

// MemoryCache stores a cache of Kubernetes objects in memory.
type MemoryCache struct {
	store map[CacheKey]*unstructured.Unstructured
	mu    sync.Mutex
}

var _ Cache = (*MemoryCache)(nil)

// NewMemoryCache creates on instance of MemoryCache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[CacheKey]*unstructured.Unstructured),
	}
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

	m.store[key] = obj
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

	return nil
}
