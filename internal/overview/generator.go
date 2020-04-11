package overview

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type pathFilter struct {
	path      string
	describer Describer

	re *regexp.Regexp
}

func newPathFilter(path string, describer Describer) *pathFilter {
	re := regexp.MustCompile(fmt.Sprintf("^%s/?$", path))

	return &pathFilter{
		re:        re,
		path:      path,
		describer: describer,
	}
}

func (pf *pathFilter) Match(path string) bool {
	return pf.re.MatchString(path)
}

func (pf *pathFilter) Fields(path string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(path)
	for i, name := range pf.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}

var (
	workloadsCronJobs = NewResource(ResourceOptions{
		Path:       "/workloads/cron-jobs",
		CacheKey:   CacheKey{APIVersion: "batch/v1beta1", Kind: "CronJob"},
		ListType:   &batch.CronJobList{},
		ObjectType: &batch.CronJob{},
		Titles:     ResourceTitle{List: "Cron Jobs", Object: "Cron Job"},
		Transforms: cronJobTransforms,
	})

	workloadsDaemonSets = NewResource(ResourceOptions{
		Path:       "/workloads/daemon-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:   &extensions.DaemonSetList{},
		ObjectType: &extensions.DaemonSet{},
		Titles:     ResourceTitle{List: "Daemon Sets", Object: "Daemon Set"},
		Transforms: daemonSetTransforms,
	})

	workloadsDeployments = NewResource(ResourceOptions{
		Path:       "/workloads/deployments",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:   &extensions.DeploymentList{},
		ObjectType: &extensions.Deployment{},
		Titles:     ResourceTitle{List: "Deployments", Object: "Deployment"},
		Transforms: deploymentTransforms,
	})

	workloadsJobs = NewResource(ResourceOptions{
		Path:       "/workloads/jobs",
		CacheKey:   CacheKey{APIVersion: "batch/v1", Kind: "Job"},
		ListType:   &batch.JobList{},
		ObjectType: &batch.Job{},
		Titles:     ResourceTitle{List: "Jobs", Object: "Job"},
		Transforms: jobTransforms,
	})

	workloadsPods = NewResource(ResourceOptions{
		Path:       "/workloads/pods",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Pod"},
		ListType:   &core.PodList{},
		ObjectType: &core.Pod{},
		Titles:     ResourceTitle{List: "Pods", Object: "Pod"},
		Transforms: podTransforms,
	})

	workloadsReplicaSets = NewResource(ResourceOptions{
		Path:       "/workloads/replica-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:   &extensions.ReplicaSetList{},
		ObjectType: &extensions.ReplicaSet{},
		Titles:     ResourceTitle{List: "Replica Sets", Object: "Replica Set"},
		Transforms: replicaSetTransforms,
	})

	workloadsReplicationControllers = NewResource(ResourceOptions{
		Path:       "/workloads/replication-controllers",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:   &core.ReplicationControllerList{},
		ObjectType: &core.ReplicationController{},
		Titles:     ResourceTitle{List: "Replication Controllers", Object: "Replication Controller"},
		Transforms: replicationControllerTransforms,
	})
	workloadsStatefulSets = NewResource(ResourceOptions{
		Path:       "/workloads/stateful-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:   &apps.StatefulSetList{},
		ObjectType: &apps.StatefulSet{},
		Titles:     ResourceTitle{List: "Stateful Sets", Object: "Stateful Set"},
		Transforms: statefulSetTransforms,
	})

	workloadsDescriber = NewSectionDescriber(
		"/workloads",
		workloadsCronJobs.List(),
		workloadsDaemonSets.List(),
		workloadsDeployments.List(),
		workloadsJobs.List(),
		workloadsPods.List(),
		workloadsReplicaSets.List(),
		workloadsReplicationControllers.List(),
		workloadsStatefulSets.List(),
	)

	dlbIngresses = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/ingresses",
		CacheKey:   CacheKey{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ListType:   &extensions.IngressList{},
		ObjectType: &extensions.Ingress{},
		Titles:     ResourceTitle{List: "Ingresses", Object: "Ingress"},
		Transforms: ingressTransforms,
	})

	dlbServices = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/services",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Service"},
		ListType:   &core.ServiceList{},
		ObjectType: &core.Service{},
		Titles:     ResourceTitle{List: "Services", Object: "Service"},
		Transforms: serviceTransforms,
	})

	discoveryAndLoadBalancingDescriber = NewSectionDescriber(
		"/discovery-and-load-balancing",
		dlbIngresses.List(),
		dlbServices.List(),
	)

	csConfigMaps = NewResource(ResourceOptions{
		Path:       "/config-and-storage/config-maps",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:   &core.ConfigMapList{},
		ObjectType: &core.ConfigMap{},
		Titles:     ResourceTitle{List: "Config Maps", Object: "Config Map"},
		Transforms: configMapTransforms,
	})

	csPVCs = NewResource(ResourceOptions{
		Path:       "/config-and-storage/persistent-volume-claims",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:   &core.PersistentVolumeClaimList{},
		ObjectType: &core.PersistentVolumeClaim{},
		Titles:     ResourceTitle{List: "Persistent Volume Claims", Object: "Persistent Volume Claim"},
		Transforms: pvcTransforms,
	})

	csSecrets = NewResource(ResourceOptions{
		Path:       "/config-and-storage/secrets",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Secret"},
		ListType:   &core.SecretList{},
		ObjectType: &core.Secret{},
		Titles:     ResourceTitle{List: "Secrets", Object: "Secret"},
		Transforms: secretTransforms,
	})

	configAndStorageDescriber = NewSectionDescriber(
		"/config-and-storage",
		csConfigMaps.List(),
		csPVCs.List(),
		csSecrets.List(),
	)

	customResourcesDescriber = NewSectionDescriber(
		"/custom-resources",
	)

	rbacDescriber = NewSectionDescriber(
		"/rbac",
	)

	rootDescriber = NewSectionDescriber(
		"/",
		workloadsDescriber,
		discoveryAndLoadBalancingDescriber,
		configAndStorageDescriber,
		customResourcesDescriber,
		rbacDescriber,
	)

	eventsDescriber = NewEventsDescriber("/events")
)

var navPaths = []string{
	"/rbac/roles",
	"/rbac/role-bindings",
}

var contentNotFound = errors.Errorf("content not found")

type generator interface {
	Generate(path, prefix, namespace string) ([]Content, error)
}

type realGenerator struct {
	cache       Cache
	pathFilters []pathFilter

	mu sync.Mutex
}

func newGenerator(cache Cache, pathFilters []pathFilter) *realGenerator {
	return &realGenerator{
		cache:       cache,
		pathFilters: pathFilters,
	}
}

func (g *realGenerator) Generate(path, prefix, namespace string) ([]Content, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if stringInSlice(path, navPaths) {
		return stubContent(path), nil
	}

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)

		return pf.describer.Describe(prefix, namespace, g.cache, fields)
	}

	return nil, contentNotFound
}

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

func stubContent(name string) []Content {
	t := newTable(name)
	t.Columns = []tableColumn{
		{Name: "foo", Accessor: "foo"},
		{Name: "bar", Accessor: "bar"},
		{Name: "baz", Accessor: "baz"},
	}

	t.Rows = []tableRow{
		{
			"foo": newStringText("r1c1"),
			"bar": newStringText("r1c2"),
			"baz": newStringText("r1c3"),
		},
		{
			"foo": newStringText("r2c1"),
			"bar": newStringText("r2c2"),
			"baz": newStringText("r2c3"),
		},
		{
			"foo": newStringText("r3c1"),
			"bar": newStringText("r3c2"),
			"baz": newStringText("r3c3"),
		},
	}

	return []Content{t}
}
