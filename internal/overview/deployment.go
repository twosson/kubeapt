package overview

import (
	"context"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/content"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
	"log"
	"reflect"
	"sort"
)

type DeploymentSummary struct{}

var _ View = (*DeploymentSummary)(nil)

func NewDeploymentSummary() *DeploymentSummary {
	return &DeploymentSummary{}
}

func (ds *DeploymentSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	deployment, err := retrieveDeployment(object)
	if err != nil {
		return nil, err
	}

	return ds.summary(deployment)
}

func (ds *DeploymentSummary) summary(deployment *extensions.Deployment) ([]content.Content, error) {
	section, err := printDeploymentSummary(deployment)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{section})
	return []content.Content{
		&summary,
	}, nil
}

type DeploymentReplicaSets struct{}

var _ View = (*DeploymentReplicaSets)(nil)

func NewDeploymentReplicaSets() *DeploymentReplicaSets {
	return &DeploymentReplicaSets{}
}

func (drs *DeploymentReplicaSets) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	var contents []content.Content

	deployment, err := retrieveDeployment(object)
	if err != nil {
		log.Printf("wtf: %v", err)
		return nil, err
	}

	replicaSetContent, err := drs.replicaSets(deployment, c)
	if err != nil {
		log.Printf("wtf2: %v", err)
		return nil, err
	}
	contents = append(contents, replicaSetContent...)

	return contents, nil
}

func (drs *DeploymentReplicaSets) replicaSets(deployment *extensions.Deployment, c Cache) ([]content.Content, error) {
	contents := []content.Content{}

	replicaSets, err := listReplicaSets(deployment, c)
	if err != nil {
		return nil, err
	}

	newReplicaSet := findNewReplicaSet(deployment, replicaSets)

	err = printContentObject(
		"New Replica Set",
		"",
		"",
		replicaSetTransforms,
		newReplicaSet,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print new replica set")
	}

	oldList := &extensions.ReplicaSetList{}
	for _, rs := range findOldReplicaSets(deployment, replicaSets) {
		oldList.Items = append(oldList.Items, *rs)
	}

	err = printContentObject(
		"Old Replica Sets",
		"",
		"",
		replicaSetTransforms,
		oldList,
		&contents,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to print old replica sets")
	}

	return contents, nil
}

func printContentObject(title, namespace, prefix string, transforms map[string]lookupFunc, object runtime.Object, contents *[]content.Content) error {
	if reflect.ValueOf(object).IsNil() {
		return errors.New("unable to print a nil object")
	}

	otf := summaryFunc(title, transforms)
	transformed := otf(namespace, prefix, contents)
	return printObject(object, transformed)
}

func retrieveDeployment(object runtime.Object) (*extensions.Deployment, error) {
	deployment, ok := object.(*extensions.Deployment)
	if !ok {
		return nil, errors.Errorf("expected object to be a Deployment, it was %T", object)
	}

	return deployment, nil
}

func listReplicaSets(deployment *extensions.Deployment, c Cache) ([]*extensions.ReplicaSet, error) {
	key := CacheKey{
		Namespace:  deployment.GetNamespace(),
		APIVersion: deployment.APIVersion,
		Kind:       "ReplicaSet",
	}

	replicaSets, err := loadReplicaSets(key, c, deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	var owned []*extensions.ReplicaSet
	for _, rs := range replicaSets {
		if metav1.IsControlledBy(rs, deployment) {
			owned = append(owned, rs)
		}
	}

	return owned, nil
}

func loadReplicaSets(key CacheKey, c Cache, selector *metav1.LabelSelector) ([]*extensions.ReplicaSet, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*extensions.ReplicaSet

	for _, object := range objects {
		rs := &extensions.ReplicaSet{}
		if err := scheme.Scheme.Convert(object, rs, 0); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(rs, object); err != nil {
			return nil, err
		}

		if selector == nil {
			list = append(list, rs)
		} else if isEqualSelector(selector, rs.Spec.Selector) {
			list = append(list, rs)
		}
	}

	return list, nil
}

// extraKeys are keys that should be ignored in labels. These keys are added
// by tools or by Kubernetes itself.
var extraKeys = []string{
	"statefulset.kubernetes.io/pod-name",
	extensions.DefaultDeploymentUniqueLabelKey,
	"controller-revision-hash",
	"pod-template-generation",
}

func isEqualSelector(s1, s2 *metav1.LabelSelector) bool {
	s1Copy := s1.DeepCopy()
	s2Copy := s2.DeepCopy()

	for _, key := range extraKeys {
		delete(s1Copy.MatchLabels, key)
		delete(s2Copy.MatchLabels, key)
	}

	return apiequality.Semantic.DeepEqual(s1Copy, s2Copy)
}

func equalIgnoreHash(template1, template2 *core.PodTemplateSpec) bool {
	t1Copy := template1.DeepCopy()
	t2Copy := template2.DeepCopy()

	for _, key := range extraKeys {
		delete(t1Copy.Labels, key)
		delete(t2Copy.Labels, key)
	}

	return apiequality.Semantic.DeepEqual(*t1Copy, *t2Copy)
}

func findNewReplicaSet(deployment *extensions.Deployment, rsList []*extensions.ReplicaSet) *extensions.ReplicaSet {
	sort.Sort(replicaSetsByCreationTimestamp(rsList))
	for i := range rsList {
		if equalIgnoreHash(&rsList[i].Spec.Template, &deployment.Spec.Template) {
			// In rare cases, such as after cluster upgrades, Deployment may end up with
			// having more than one new ReplicaSets that have the same template as its template,
			// see https://github.com/kubernetes/kubernetes/issues/40415
			// We deterministically choose the oldest new ReplicaSet.
			return rsList[i]
		}
	}

	// new ReplicaSet does not exist.
	return nil
}

// findOldReplicaSets returns the old replica sets targeted by the given Deployment, with the given slice of RSes.
// Note that the first set of old replica sets doesn't include the ones with no pods, and the second set of old replica sets include all old replica sets.
func findOldReplicaSets(deployment *extensions.Deployment, rsList []*extensions.ReplicaSet) []*extensions.ReplicaSet {
	var requiredRSs []*extensions.ReplicaSet
	newRS := findNewReplicaSet(deployment, rsList)
	for _, rs := range rsList {
		// Filter out new replica set
		if newRS != nil && rs.UID == newRS.UID {
			continue
		}
		if rs.Spec.Replicas != 0 {
			requiredRSs = append(requiredRSs, rs)
		}
	}
	return requiredRSs
}

// replicaSetsByCreationTimestamp sorts a list of ReplicaSet by creation timestamp, using their names as a tie breaker.
type replicaSetsByCreationTimestamp []*extensions.ReplicaSet

func (o replicaSetsByCreationTimestamp) Len() int      { return len(o) }
func (o replicaSetsByCreationTimestamp) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o replicaSetsByCreationTimestamp) Less(i, j int) bool {
	if o[i].CreationTimestamp.Equal(&o[j].CreationTimestamp) {
		return o[i].Name < o[j].Name
	}
	return o[i].CreationTimestamp.Before(&o[j].CreationTimestamp)
}
