package fake

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/third_party/dynamicfake"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
)

var (
	scheme         = runtime.NewScheme()
	codecs         = serializer.NewCodecFactory(scheme)
	parameterCodec = runtime.NewParameterCodec(scheme)
)

// Client implements cluster Interface.
type Client struct {
	FakeDynamic   *dynamicfake.FakeDynamicClient
	FakeDiscovery *fakediscovery.FakeDiscovery
}

// NewClient creates an instance of Client.
func NewClient(scheme *runtime.Scheme, objects []runtime.Object) (*Client, error) {
	client := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
	if !ok {
		return nil, errors.New("couldn't convert Discovery() to *FakeDiscovery")
	}

	dynamicClient := NewSimpleDynamicClient(scheme, fakeDiscovery, objects...)

	return &Client{
		FakeDynamic:   dynamicClient,
		FakeDiscovery: fakeDiscovery,
	}, nil
}

// DynamicClient returns a dynamic client or an error.
func (c *Client) DynamicClient() (dynamic.Interface, error) {
	return c.FakeDynamic, nil
}

// DiscoveryClient returns a discovery client or an error.
func (c *Client) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return c.FakeDiscovery, nil
}

// NamespaceClient returns a namespace client or an error.
func (c *Client) NamespaceClient() (cluster.NamespaceInterface, error) {
	return &NamespaceClient{}, nil
}

func kindFor(discoveryClient discovery.DiscoveryInterface, gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	l, err := discoveryClient.ServerResourcesForGroupVersion(gvr.GroupVersion().String())
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	for _, r := range l.APIResources {
		if r.Name != gvr.Resource {
			continue
		}
		return schema.GroupVersionKind{
			Group:   r.Group,
			Version: r.Version,
			Kind:    r.Kind,
		}, nil
	}
	return schema.GroupVersionKind{}, errors.New("not found")
}

// NewSimpleDynamicClient creates a FakeDynamicClient which fixes behavior from dynamicfake.NewSimpleDynamicClient - we properly forward
// ADDED events for preexisting objects when adding watches.
func NewSimpleDynamicClient(scheme *runtime.Scheme, discoveryClient discovery.DiscoveryInterface, objects ...runtime.Object) *dynamicfake.FakeDynamicClient {
	// In order to use List with this client, you have to have the v1.List registered in your scheme. Neat thing though
	// it does NOT have to be the *same* list
	scheme.AddKnownTypeWithName(schema.GroupVersionKind{Group: "fake-dynamic-client-group", Version: "v1", Kind: "List"}, &unstructured.UnstructuredList{})

	codecs := serializer.NewCodecFactory(scheme)
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &dynamicfake.FakeDynamicClient{}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		w, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}

		gvk, err := kindFor(discoveryClient, gvr)
		if err != nil {
			return false, nil, fmt.Errorf("no registered kind for resource: %v", gvr.String())
		}

		l, err := o.List(gvr, gvk, ns)
		if err != nil {
			return false, nil, errors.Wrap(err, "listing existing objects")
		}

		// Replay existing objects
		rfw, ok := w.(*watch.RaceFreeFakeWatcher)
		if !ok {
			return false, nil, fmt.Errorf("unexpected watch type: %T", w)
		}

		ul, ok := l.(*unstructured.UnstructuredList)
		if !ok {
			return false, nil, errors.Errorf("wrong type for list: %T\n", l)
		}

		err = ul.EachListItem(func(obj runtime.Object) error {
			rfw.Add(obj)
			return nil
		})

		return true, w, nil
	})

	return cs
}
