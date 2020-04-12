package overview

import (
	"fmt"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/content"
	"github.com/twosson/kubeapt/internal/view"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"path"
	"reflect"
)

func resourceLink(sectionType, resourceType string) lookupFunc {
	return func(namespace, prefix string, cell interface{}) content.Text {
		name := fmt.Sprintf("%v", cell)
		resourcePath := path.Join("/content", "overview", sectionType, resourceType, name)
		return content.NewLinkText(name, resourcePath)
	}
}

type ResourceTitle struct {
	List   string
	Object string
}

type ResourceOptions struct {
	Path       string
	CacheKey   CacheKey
	ListType   interface{}
	ObjectType interface{}
	Titles     ResourceTitle
	Transforms map[string]lookupFunc
	Views      []view.View
}

type Resource struct {
	ResourceOptions
}

func NewResource(options ResourceOptions) *Resource {
	return &Resource{
		ResourceOptions: options,
	}
}

func (r *Resource) Describe(prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) ([]content.Content, string, error) {
	return r.List().Describe(prefix, namespace, clusterClient, options)
}

func (r *Resource) List() *ListDescriber {
	return NewListDescriber(
		r.Path,
		r.Titles.List,
		r.CacheKey,
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ListType).Elem().Type()).Interface()
		},
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		summaryFunc(r.Titles.List, r.Transforms),
	)
}

func (r *Resource) Object() *ObjectDescriber {
	return NewObjectDescriber(
		path.Join(r.Path, "(?P<name>.*?)"),
		r.Titles.Object,
		r.CacheKey,
		func() interface{} {
			return reflect.New(reflect.ValueOf(r.ObjectType).Elem().Type()).Interface()
		},
		summaryFunc(r.Titles.Object, r.Transforms),
		r.Views,
	)
}

func (r *Resource) PathFilters() []pathFilter {
	filters := []pathFilter{
		*newPathFilter(r.Path, r.List()),
		*newPathFilter(path.Join(r.Path, "(?P<name>.*?)"), r.Object()),
	}

	return filters
}

var defaultTransforms = map[string]lookupFunc{
	"Labels": func(namespace, prefix string, cell interface{}) content.Text {
		text := fmt.Sprintf("%v", cell)
		return content.NewStringText(text)
	},
}

func buildTransforms(transforms map[string]lookupFunc) map[string]lookupFunc {
	m := make(map[string]lookupFunc)
	for k, v := range defaultTransforms {
		m[k] = v
	}
	for k, v := range transforms {
		m[k] = v
	}

	return m
}

// summaryFunc creates an ObjectTransformFunc given a title and a lookup.
func summaryFunc(title string, m map[string]lookupFunc) ObjecTransformFunc {
	return func(namespace, prefix string, contents *[]content.Content) func(*metav1beta1.Table) error {
		return func(tbl *metav1beta1.Table) error {
			contentTable, err := printContentTable(title, namespace, prefix, tbl, m)
			if err != nil {
				return err
			}

			*contents = append(*contents, contentTable)
			return nil
		}
	}
}
