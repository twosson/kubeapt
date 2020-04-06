package overview

import (
	"github.com/twosson/kubeapt/internal/apt"
	"path"
)

func navigationEntries(root string) (*apt.Navigation, error) {
	n := &apt.Navigation{
		Title: "Overview",
		Path:  path.Join(root, "/"),
		Children: []*apt.Navigation{
			{
				Title: "Workloads",
				Path:  path.Join(root, "workloads"),
				Children: []*apt.Navigation{
					{
						Title: "Cron Jobs",
						Path:  path.Join(root, "workloads/cron-jobs"),
					},
					{
						Title: "Daemon Sets",
						Path:  path.Join(root, "workloads/daemon-sets"),
					},
					{
						Title: "Deployments",
						Path:  path.Join(root, "workloads/deployments"),
					},
					{
						Title: "Jobs",
						Path:  path.Join(root, "workloads/jobs"),
					},
					{
						Title: "Pods",
						Path:  path.Join(root, "workloads/pods"),
					},
					{
						Title: "Replica Sets",
						Path:  path.Join(root, "workloads/replica-sets"),
					},
					{
						Title: "Replication Controllers",
						Path:  path.Join(root, "workloads/replication-controllers"),
					},
					{
						Title: "Stateful Sets",
						Path:  path.Join(root, "workloads/stateful-sets"),
					},
				},
			},
			{
				Title: "Discovery and Load Balancing",
				Path:  path.Join(root, "discovery-and-load-balancing"),
				Children: []*apt.Navigation{
					{
						Title: "Ingresses",
						Path:  path.Join(root, "discovery-and-load-balancing/ingresses"),
					},
					{
						Title: "Services",
						Path:  path.Join(root, "discovery-and-load-balancing/services"),
					},
				},
			},
			{
				Title: "Config and Storage",
				Path:  path.Join(root, "config-and-storage"),
				Children: []*apt.Navigation{
					{
						Title: "Config Maps",
						Path:  path.Join(root, "config-and-storage/config-maps"),
					},
					{
						Title: "Persistent Volume Claims",
						Path:  path.Join(root, "config-and-storage/persistent-volume-claims"),
					},
					{
						Title: "Secrets",
						Path:  path.Join(root, "config-and-storage/secrets"),
					},
				},
			},
			{
				Title: "Custom Resources",
				Path:  path.Join(root, "custom-resources"),
			},
			{
				Title: "RBAC",
				Path:  path.Join(root, "rbac"),
				Children: []*apt.Navigation{
					{
						Title: "Roles",
						Path:  path.Join(root, "rbac/roles"),
					},
					{
						Title: "Role Bindings",
						Path:  path.Join(root, "rbac/role-bindings"),
					},
				},
			},
			{
				Title: "Events",
				Path:  path.Join(root, "events"),
			},
		},
	}

	return n, nil
}
