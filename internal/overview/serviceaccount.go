package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/apis/core"
)

type ServiceAccountSummary struct{}

var _ View = (*ServiceAccountSummary)(nil)

func NewServiceAccountSummary() *ServiceAccountSummary {
	return &ServiceAccountSummary{}
}

func (js *ServiceAccountSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	serviceAccount, err := retrieveServiceAccount(object)
	if err != nil {
		return nil, err
	}

	tokens := []*core.Secret{}

	// missingSecrets is the set of all secrets present in the
	// serviceAccount but not present in the set of existing secrets.
	missingSecrets := sets.NewString()
	secrets, err := listSecrets(serviceAccount.GetNamespace(), c)
	if err != nil {
		return nil, err
	}

	// existingSecrets is the set of all secrets remaining on a
	// service account that are not present in the "tokens" slice.
	existingSecrets := sets.NewString()

	for _, s := range secrets {
		if s.Type == core.SecretTypeServiceAccountToken {
			name, _ := s.Annotations[corev1.ServiceAccountNameKey]
			uid, _ := s.Annotations[corev1.ServiceAccountUIDKey]
			if name == serviceAccount.Name && uid == string(serviceAccount.UID) {
				tokens = append(tokens, s)
			}
		}
		existingSecrets.Insert(s.Name)
	}

	for _, s := range serviceAccount.Secrets {
		if !existingSecrets.Has(s.Name) {
			missingSecrets.Insert(s.Name)
		}
	}
	for _, s := range serviceAccount.ImagePullSecrets {
		if !existingSecrets.Has(s.Name) {
			missingSecrets.Insert(s.Name)
		}
	}

	detail, err := printServiceAccountSummary(serviceAccount, tokens, missingSecrets)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

func retrieveServiceAccount(object runtime.Object) (*core.ServiceAccount, error) {
	rc, ok := object.(*core.ServiceAccount)
	if !ok {
		return nil, errors.Errorf("expected object to be a ServiceAccount, it was %T", object)
	}

	return rc, nil
}
