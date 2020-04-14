package overview

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/cluster"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"regexp"
	"sync"
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
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewCronJobSummary,
					NewCronJobJobs,
					NewEventList,
				},
				Title: "Summary",
			},
			{
				Views: []ViewFactory{
					NewCronJobSummary,
				},
				Title: "Resource Viewer",
			},
		},
	})

	workloadsDaemonSets = NewResource(ResourceOptions{
		Path:       "/workloads/daemon-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "DaemonSet"},
		ListType:   &extensions.DaemonSetList{},
		ObjectType: &extensions.DaemonSet{},
		Titles:     ResourceTitle{List: "Daemon Sets", Object: "Daemon Set"},
		Transforms: daemonSetTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewDaemonSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})

	workloadsDeployments = NewResource(ResourceOptions{
		Path:       "/workloads/deployments",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "Deployment"},
		ListType:   &extensions.DeploymentList{},
		ObjectType: &extensions.Deployment{},
		Titles:     ResourceTitle{List: "Deployments", Object: "Deployment"},
		Transforms: deploymentTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewDeploymentSummary,
					NewContainerSummary,
					NewDeploymentReplicaSets,
					NewEventList,
				},
			},
		},
	})

	workloadsJobs = NewResource(ResourceOptions{
		Path:       "/workloads/jobs",
		CacheKey:   CacheKey{APIVersion: "batch/v1", Kind: "Job"},
		ListType:   &batch.JobList{},
		ObjectType: &batch.Job{},

		Titles:     ResourceTitle{List: "Jobs", Object: "Job"},
		Transforms: jobTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewJobSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})

	workloadsPods = NewResource(ResourceOptions{
		Path:       "/workloads/pods",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Pod"},
		ListType:   &core.PodList{},
		ObjectType: &core.Pod{},
		Titles:     ResourceTitle{List: "Pods", Object: "Pod"},
		Transforms: podTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewPodSummary,
					NewPodContainer,
					NewPodCondition,
					NewPodVolume,
					NewEventList,
				},
			},
		},
	})

	workloadsReplicaSets = NewResource(ResourceOptions{
		Path:       "/workloads/replica-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "ReplicaSet"},
		ListType:   &extensions.ReplicaSetList{},
		ObjectType: &extensions.ReplicaSet{},
		Titles:     ResourceTitle{List: "Replica Sets", Object: "Replica Set"},
		Transforms: replicaSetTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewReplicaSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})

	workloadsReplicationControllers = NewResource(ResourceOptions{
		Path:       "/workloads/replication-controllers",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "ReplicationController"},
		ListType:   &core.ReplicationControllerList{},
		ObjectType: &core.ReplicationController{},
		Titles:     ResourceTitle{List: "Replication Controllers", Object: "Replication Controller"},
		Transforms: replicationControllerTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewReplicationControllerSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})
	workloadsStatefulSets = NewResource(ResourceOptions{
		Path:       "/workloads/stateful-sets",
		CacheKey:   CacheKey{APIVersion: "apps/v1", Kind: "StatefulSet"},
		ListType:   &apps.StatefulSetList{},
		ObjectType: &apps.StatefulSet{},
		Titles:     ResourceTitle{List: "Stateful Sets", Object: "Stateful Set"},
		Transforms: statefulSetTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewStatefulSetSummary,
					NewContainerSummary,
					NewPodList,
					NewEventList,
				},
			},
		},
	})

	workloadsDescriber = NewSectionDescriber(
		"/workloads",
		"Workloads",
		workloadsCronJobs,
		workloadsDaemonSets,
		workloadsDeployments,
		workloadsJobs,
		workloadsPods,
		workloadsReplicaSets,
		workloadsReplicationControllers,
		workloadsStatefulSets,
	)

	dlbIngresses = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/ingresses",
		CacheKey:   CacheKey{APIVersion: "extensions/v1beta1", Kind: "Ingress"},
		ListType:   &v1beta1.IngressList{},
		ObjectType: &v1beta1.Ingress{},
		Titles:     ResourceTitle{List: "Ingresses", Object: "Ingress"},
		Transforms: ingressTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewIngressSummary,
					NewIngressDetails,
					NewEventList,
				},
			},
		},
	})

	dlbServices = NewResource(ResourceOptions{
		Path:       "/discovery-and-load-balancing/services",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Service"},
		ListType:   &core.ServiceList{},
		ObjectType: &core.Service{},
		Titles:     ResourceTitle{List: "Services", Object: "Service"},
		Transforms: serviceTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewServiceSummary,
					NewServicePort,
					NewServiceEndpoints,
					NewEventList,
				},
			},
		},
	})

	discoveryAndLoadBalancingDescriber = NewSectionDescriber(
		"/discovery-and-load-balancing",
		"Discovery and Load Balancing",
		dlbIngresses,
		dlbServices,
	)

	csConfigMaps = NewResource(ResourceOptions{
		Path:       "/config-and-storage/config-maps",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "ConfigMap"},
		ListType:   &core.ConfigMapList{},
		ObjectType: &core.ConfigMap{},
		Titles:     ResourceTitle{List: "Config Maps", Object: "Config Map"},
		Transforms: configMapTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewConfigMapSummary,
					NewConfigMapDetails,
					NewEventList,
				},
			},
		},
	})

	csPVCs = NewResource(ResourceOptions{
		Path:       "/config-and-storage/persistent-volume-claims",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "PersistentVolumeClaim"},
		ListType:   &core.PersistentVolumeClaimList{},
		ObjectType: &core.PersistentVolumeClaim{},
		Titles:     ResourceTitle{List: "Persistent Volume Claims", Object: "Persistent Volume Claim"},
		Transforms: pvcTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewPersistentVolumeClaimSummary,
					NewEventList,
				},
			},
		},
	})

	csSecrets = NewResource(ResourceOptions{
		Path:       "/config-and-storage/secrets",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Secret"},
		ListType:   &core.SecretList{},
		ObjectType: &core.Secret{},
		Titles:     ResourceTitle{List: "Secrets", Object: "Secret"},
		Transforms: secretTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewSecretSummary,
					NewSecretData,
					NewEventList,
				},
			},
		},
	})

	csServiceAccounts = NewResource(ResourceOptions{
		Path:       "/config-and-storage/service-accounts",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "ServiceAccount"},
		ListType:   &core.ServiceAccountList{},
		ObjectType: &core.ServiceAccount{},
		Titles:     ResourceTitle{List: "Service Accounts", Object: "Service Account"},
		Transforms: serviceAccountTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewServiceAccountSummary,
					NewEventList,
				},
			},
		},
	})

	configAndStorageDescriber = NewSectionDescriber(
		"/config-and-storage",
		"Config and Storage",
		csConfigMaps,
		csPVCs,
		csSecrets,
		csServiceAccounts,
	)

	customResourcesDescriber = NewSectionDescriber(
		"/custom-resources",
		"Custom Resources",
	)

	rbacRoles = NewResource(ResourceOptions{
		Path:       "/rbac/roles",
		CacheKey:   CacheKey{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "Role"},
		ListType:   &rbac.RoleList{},
		ObjectType: &rbac.Role{},
		Titles:     ResourceTitle{List: "Roles", Object: "Role"},
		Transforms: roleTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewRoleSummary,
					NewRoleRule,
					NewEventList,
				},
			},
		},
	})

	rbacRoleBindings = NewResource(ResourceOptions{
		Path:       "/rbac/role-bindings",
		CacheKey:   CacheKey{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "RoleBinding"},
		ListType:   &rbac.RoleBindingList{},
		ObjectType: &rbac.RoleBinding{},
		Titles:     ResourceTitle{List: "Role Bindings", Object: "Role Binding"},
		Transforms: roleBindingTransforms,
		Sections: []ContentSection{
			{
				Views: []ViewFactory{
					NewRoleBindingSummary,
					NewRoleBindingSubjects,
					NewEventList,
				},
			},
		},
	})

	rbacDescriber = NewSectionDescriber(
		"/rbac",
		"RBAC",
		rbacRoles,
		rbacRoleBindings,
	)

	rootDescriber = NewSectionDescriber(
		"/",
		"Overview",
		workloadsDescriber,
		discoveryAndLoadBalancingDescriber,
		configAndStorageDescriber,
		customResourcesDescriber,
		rbacDescriber,
	)

	eventsDescriber = NewResource(ResourceOptions{
		Path:       "/events",
		CacheKey:   CacheKey{APIVersion: "v1", Kind: "Event"},
		ListType:   &core.EventList{},
		ObjectType: &core.Event{},
		Titles:     ResourceTitle{List: "Events", Object: "Event"},
		Transforms: roleBindingTransforms,
	})
)

var contentNotFound = errors.Errorf("content not found")

type generator interface {
	Generate(ctx context.Context, path, prefix, namespace string) (ContentResponse, error)
}

type realGenerator struct {
	cache         Cache
	pathFilters   []pathFilter
	clusterClient cluster.ClientInterface

	mu sync.Mutex
}

func newGenerator(cache Cache, pathFilters []pathFilter, clusterClient cluster.ClientInterface) *realGenerator {
	return &realGenerator{
		cache:         cache,
		pathFilters:   pathFilters,
		clusterClient: clusterClient,
	}
}

func (g *realGenerator) Generate(ctx context.Context, path, prefix, namespace string) (ContentResponse, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)
		options := DescriberOptions{
			Cache:  g.cache,
			Fields: fields,
		}

		cResponse, err := pf.describer.Describe(ctx, prefix, namespace, g.clusterClient, options)
		if err != nil {
			return emptyContentResponse, err
		}

		return cResponse, nil
	}

	return emptyContentResponse, contentNotFound
}
