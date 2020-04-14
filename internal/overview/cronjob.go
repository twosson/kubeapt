package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type CronJobSummary struct{}

var _ View = (*CronJobSummary)(nil)

func NewCronJobSummary() *CronJobSummary {
	return &CronJobSummary{}
}

func (rss *CronJobSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	cronJob, err := retrieveCronJob(object)
	if err != nil {
		return nil, err
	}

	return rss.summary(cronJob, c)
}

func (rss *CronJobSummary) summary(cronJob *batch.CronJob, c Cache) ([]content.Content, error) {
	jobs, err := listJobs(cronJob.GetNamespace(), cronJob.GetUID(), c)
	if err != nil {
		return nil, err
	}

	section, err := printCronJobSummary(cronJob, jobs)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	contents := []content.Content{
		&summary,
	}

	return contents, nil
}

type CronJobJobs struct{}

var _ View = (*CronJobJobs)(nil)

func NewCronJobJobs() *CronJobJobs {
	return &CronJobJobs{}
}

func (j *CronJobJobs) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	cronJob, err := retrieveCronJob(object)
	if err != nil {
		return nil, err
	}

	return j.jobs(cronJob, c)
}

func (j *CronJobJobs) jobs(cronJob *batch.CronJob, c Cache) ([]content.Content, error) {
	jobs, err := listJobs(cronJob.GetNamespace(), cronJob.GetUID(), c)
	if err != nil {
		return nil, err
	}

	var active, inactive batch.JobList
	for _, job := range jobs {
		if job.Status.Active == 1 {
			active.Items = append(active.Items, *job)
		} else {
			inactive.Items = append(inactive.Items, *job)
		}
	}

	var contents []content.Content

	err = printContentObject(
		"Active Jobs",
		"ns",
		"prefix",
		jobTransforms,
		&active,
		&contents,
	)
	if err != nil {
		return nil, err
	}

	err = printContentObject(
		"Inactive Jobs",
		"ns",
		"prefix",
		jobTransforms,
		&inactive,
		&contents,
	)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func retrieveCronJob(object runtime.Object) (*batch.CronJob, error) {
	replicaSet, ok := object.(*batch.CronJob)
	if !ok {
		return nil, errors.Errorf("expected object to be a CronJob, it was %T", object)
	}

	return replicaSet, nil
}

func listJobs(namespace string, uid types.UID, c Cache) ([]*batch.Job, error) {
	key := CacheKey{
		Namespace:  namespace,
		APIVersion: "batch/v1",
		Kind:       "Job",
	}

	jobs, err := loadJobs(key, c)
	if err != nil {
		return nil, err
	}

	var owned []*batch.Job
	for _, job := range jobs {
		controllerRef := metav1.GetControllerOf(job)
		if controllerRef == nil || controllerRef.UID != uid {
			continue
		}

		owned = append(owned, job)
	}

	return owned, nil
}

func loadJobs(key CacheKey, c Cache) ([]*batch.Job, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*batch.Job

	for _, object := range objects {
		job := &batch.Job{}
		if err := scheme.Scheme.Convert(object, job, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(job, object); err != nil {
			return nil, err
		}

		list = append(list, job)
	}

	return list, nil
}
