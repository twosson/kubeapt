package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
)

type ReplicationControllerSummary struct{}

var _ View = (*ReplicationControllerSummary)(nil)

func NewReplicationControllerSummary() *ReplicationControllerSummary {
	return &ReplicationControllerSummary{}
}

func (js *ReplicationControllerSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	rc, err := retrieveReplicationController(object)
	if err != nil {
		return nil, err
	}

	s := &metav1.LabelSelector{
		MatchLabels: rc.Spec.Selector,
	}

	pods, err := listPods(rc.GetNamespace(), s, rc.GetUID(), c)
	if err != nil {
		return nil, err
	}

	detail, err := printReplicationControllerSummary(rc, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	contents := []content.Content{
		&summary,
	}

	podTemplate, err := printPodTemplate(rc.Spec.Template, nil)
	if err != nil {
		return nil, err
	}
	contents = append(contents, podTemplate...)

	return contents, nil
}

func retrieveReplicationController(object runtime.Object) (*core.ReplicationController, error) {
	rc, ok := object.(*core.ReplicationController)
	if !ok {
		return nil, errors.Errorf("expected object to be a ReplicationController, it was %T", object)
	}

	return rc, nil
}
