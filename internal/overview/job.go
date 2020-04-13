package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/batch"
)

type JobSummary struct{}

var _ View = (*JobSummary)(nil)

func NewJobSummary() *JobSummary {
	return &JobSummary{}
}

func (js *JobSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	job, err := retrieveJob(object)
	if err != nil {
		return nil, err
	}

	s := job.Spec.Selector
	s.MatchLabels["job-name"] = job.Labels["job-name"]

	pods, err := listPods(job.GetNamespace(), s, job.GetUID(), c)
	if err != nil {
		return nil, err
	}

	detail, err := printJobSummary(job, pods)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

func retrieveJob(object runtime.Object) (*batch.Job, error) {
	job, ok := object.(*batch.Job)
	if !ok {
		return nil, errors.Errorf("expected object to be a Job, it was %T", object)
	}

	return job, nil
}
