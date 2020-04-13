package overview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/cluster/fake"
	"github.com/twosson/kubeapt/internal/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
	"time"
)

var resources = []*metav1.APIResourceList{
	{
		GroupVersion: "apps/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "deployments",
				SingularName: "deployment",
				Group:        "apps",
				Version:      "v1",
				Kind:         "Deployment",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "bar/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "bars",
				SingularName: "bar",
				Group:        "bar",
				Version:      "v1",
				Kind:         "Bar",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
	{
		GroupVersion: "foo/v1",
		APIResources: []metav1.APIResource{
			metav1.APIResource{
				Name:         "kinds",
				SingularName: "kind",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Kind",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "foos",
				SingularName: "foo",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Foo",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
			metav1.APIResource{
				Name:         "others",
				SingularName: "other",
				Group:        "foo",
				Version:      "v1",
				Kind:         "Other",
				Namespaced:   true,
				Verbs:        metav1.Verbs{"list", "watch"},
				Categories:   []string{"all"},
			},
		},
	},
}

func newCache(t *testing.T, objects []runtime.Object) (*InformerCache, error) {
	// func NewInformerCache(stopCh <-chan struct{}, client dynamic.Interface, restMapper meta.RESTMapper, opts ...InformerCacheOpt) *InformerCache {
	scheme := newScheme()

	client, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)
	if err != nil {
		return nil, err
	}
	// notifyCh := make(chan CacheNotification)
	stopCh := make(chan struct{})

	restMapper, err := client.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")
	return NewInformerCache(stopCh, client.FakeDynamic, restMapper), nil
}

func TestInformerCache_Retrieve(t *testing.T) {
	objects := []runtime.Object{}
	for _, u := range genObjectsSeed() {
		objects = append(objects, u)
	}

	cases := []struct {
		name        string
		key         CacheKey
		expectedLen int
		expectErr   bool
	}{
		{
			name: "ns, apiVersion, kind, name",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
				Name:       "foo1",
			},
			expectedLen: 1,
		},
		{
			name: "ns, apiVersion, kind",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
				Kind:       "Kind",
			},
			expectedLen: 2,
		},
		{
			name: "ns, apiVersion: error because we require kind",
			key: CacheKey{
				Namespace:  "default",
				APIVersion: "foo/v1",
			},
			expectErr: true,
		},
		/*
		   {
		       name: "ns, apiVersion",
		       key: CacheKey{
		           Namespace:  "default",
		           APIVersion: "foo/v1",
		       },
		       expectedLen: 3,
		   },
		   {
		       name: "ns",
		       key: CacheKey{
		           Namespace: "default",
		       }, expectedLen: 4,
		   },
		*/
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := newCache(t, objects)
			require.NoError(t, err)

			objs, err := c.Retrieve(tc.key)
			if err != nil {

			}
			hadErr := (err != nil)
			assert.Equalf(t, tc.expectErr, hadErr, "error mismatch: %v", err)
			assert.Len(t, objs, tc.expectedLen)
		})
	}
}

func TestInformerCache_Watch(t *testing.T) {
	scheme := newScheme()

	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy3"),
	}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	discoveryClient := clusterClient.FakeDiscovery
	discoveryClient.Resources = resources

	dynamicClient := clusterClient.FakeDynamic

	notifyCh := make(chan CacheNotification)
	notifyDone := make(chan struct{})

	restMapper, err := clusterClient.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")

	cache := NewInformerCache(notifyDone, dynamicClient, restMapper,
		InformerCacheNotificationOpt(notifyCh, notifyDone),
		InformerCacheLoggerOpt(log.TestLogger(t)),
	)

	defer func() {
		close(notifyDone)
	}()

	// verify predefined objects are present
	cacheKey := CacheKey{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
	found, err := cache.Retrieve(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 1)

	// define new object
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("deploy2")
	obj.SetNamespace("default")

	res := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	resClient := dynamicClient.Resource(res).Namespace("default")

	// create object
	_, err = resClient.Create(obj)
	require.NoError(t, err)

	// wait for cache to store an item before proceeding.
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timed out wating for create object to notify")
	case <-notifyCh:
	}

	found, err = cache.Retrieve(cacheKey)
	require.NoError(t, err)

	// 2 == initial + the new object
	require.Len(t, found, 2)

	annotations := map[string]string{"update": "update"}
	obj.SetAnnotations(annotations)

	// update object
	_, err = resClient.Update(obj)
	require.NoError(t, err)

	// wait for cache to store an item before proceeding.
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("timed out wating for update object to notify")
	case <-notifyCh:
	}

	found, err = cache.Retrieve(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 2)

	// Find the object we updated
	var match bool
	for _, u := range found {
		if u.GetName() == obj.GetName() && u.GroupVersionKind() == obj.GroupVersionKind() {
			match = true
			require.Equal(t, annotations, u.GetAnnotations())
		}
	}
	require.True(t, match, "unable to find object from fetched results")
}

func TestInformerCache_Watch_Stop(t *testing.T) {
	scheme := newScheme()

	objects := []runtime.Object{}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	discoveryClient := clusterClient.FakeDiscovery
	discoveryClient.Resources = resources

	dynamicClient := clusterClient.FakeDynamic

	notifyCh := make(chan CacheNotification)
	notifyDone := make(chan struct{})

	restMapper, err := clusterClient.RESTMapper()
	require.NoError(t, err, "fetching RESTMapper")

	cache := NewInformerCache(notifyDone, dynamicClient, restMapper,
		InformerCacheNotificationOpt(notifyCh, notifyDone),
		InformerCacheLoggerOpt(log.TestLogger(t)),
	)

	// verify predefined objects are present
	cacheKey := CacheKey{Namespace: "default", APIVersion: "apps/v1", Kind: "Deployment"}
	found, err := cache.Retrieve(cacheKey)
	require.NoError(t, err)

	require.Len(t, found, 0)

	// define new object
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("apps/v1")
	obj.SetKind("Deployment")
	obj.SetName("deploy2")
	obj.SetNamespace("default")

	res := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	resClient := dynamicClient.Resource(res).Namespace("default")

	// Stop notifications
	close(notifyDone)

	// Drain notifications
	closeDone := make(chan struct{})
	go func() {
		for range notifyCh {
		}
		close(closeDone)
	}()

	// Wait for informers to shutdown
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for notification channel to close")
	case <-closeDone:
	}

	// create object
	_, err = resClient.Create(obj)
	require.NoError(t, err)

	found, err = cache.Retrieve(cacheKey)
	require.NoError(t, err)

	// The second object is not seen because we shutdown the informer
	require.Len(t, found, 0)
}
