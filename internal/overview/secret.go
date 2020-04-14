package overview

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
	"sort"
	"strings"
)

type SecretSummary struct{}

var _ View = (*SecretSummary)(nil)

func NewSecretSummary(prefix, namespace string, c clock.Clock) View {
	return &SecretSummary{}
}

func (js *SecretSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	secret, err := retrieveSecret(object)
	if err != nil {
		return nil, err
	}

	detail, err := printSecretSummary(secret)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type SecretData struct{}

var _ View = (*SecretData)(nil)

func NewSecretData(prefix, namespace string, c clock.Clock) View {
	return &SecretData{}
}

func (js *SecretData) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	secret, err := retrieveSecret(object)
	if err != nil {
		return nil, err
	}

	dataSection := content.NewSection()
	dataSection.Title = "Data"

	var keys []string
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		data := secret.Data[key]
		switch {
		case key == core.ServiceAccountTokenKey && secret.Type == core.SecretTypeServiceAccountToken:
			dataSection.AddText(key, strings.TrimSpace(string(data)))
		default:
			dataSection.AddText(key, fmt.Sprintf("%d bytes", len(data)))
		}
	}

	summary := content.NewSummary("Data", []content.Section{dataSection})

	return []content.Content{
		&summary,
	}, nil
}

func retrieveSecret(object runtime.Object) (*core.Secret, error) {
	rc, ok := object.(*core.Secret)
	if !ok {
		return nil, errors.Errorf("expected object to be a Secret, it was %T", object)
	}

	return rc, nil
}

func listSecrets(namespace string, c Cache) ([]*core.Secret, error) {
	key := CacheKey{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}

	return loadSecrets(key, c)
}

func loadSecrets(key CacheKey, c Cache) ([]*core.Secret, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*core.Secret

	for _, object := range objects {
		e := &core.Secret{}
		if err := scheme.Scheme.Convert(object, e, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(e, object); err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}
