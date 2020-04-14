package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/rbac"
)

type RoleBindingSummary struct{}

var _ View = (*RoleBindingSummary)(nil)

func NewRoleBindingSummary(prefix, namespace string, c clock.Clock) View {
	return &RoleBindingSummary{}
}

func (js *RoleBindingSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	roleBinding, err := retrieveRoleBinding(object)
	if err != nil {
		return nil, err
	}

	role, err := getRole(roleBinding.GetNamespace(), roleBinding.RoleRef.Name, c)
	if err != nil {
		return nil, err
	}

	detail, err := printRoleBindingSummary(roleBinding, role)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type RoleBindingSubjects struct{}

var _ View = (*RoleBindingSubjects)(nil)

func NewRoleBindingSubjects(prefix, namespace string, c clock.Clock) View {
	return &RoleBindingSubjects{}
}

func (js *RoleBindingSubjects) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	roleBinding, err := retrieveRoleBinding(object)
	if err != nil {
		return nil, err
	}

	subjectsTable, err := printRoleBindingSubjects(roleBinding)
	if err != nil {
		return nil, err
	}

	return []content.Content{
		&subjectsTable,
	}, nil
}

func retrieveRoleBinding(object runtime.Object) (*rbac.RoleBinding, error) {
	rc, ok := object.(*rbac.RoleBinding)
	if !ok {
		return nil, errors.Errorf("expected object to be a RoleBinding, it was %T", object)
	}

	return rc, nil
}
