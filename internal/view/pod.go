package view

import (
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
	"time"
)

type PodCondition struct{}

func NewPodCondition() *PodCondition {
	return &PodCondition{}
}

func (pc *PodCondition) Content(ctx context.Context, object runtime.Object, clusterClient cluster.ClientInterface) ([]content.Content, error) {
	pod, err := retrievePod(object)
	if err != nil {
		return nil, err
	}

	conditions := pod.Status.Conditions

	table := content.NewTable("Conditions")
	table.Columns = []content.TableColumn{
		tableCol("Type"),
		tableCol("Status"),
		tableCol("Last probe time"),
		tableCol("Last transition time"),
		tableCol("Reason"),
		tableCol("Message"),
	}

	for _, condition := range conditions {

		lastProbeTime := condition.LastProbeTime.UTC().Format(time.RFC3339)
		lastTransitionTime := condition.LastTransitionTime.UTC().Format(time.RFC3339)

		row := content.TableRow{
			"Type":                 content.NewStringText(string(condition.Type)),
			"Status":               content.NewStringText(string(condition.Status)),
			"Last probe time":      content.NewStringText(lastProbeTime),
			"Last transition time": content.NewStringText(lastTransitionTime),
			"Reason":               content.NewStringText(condition.Reason),
			"Message":              content.NewStringText(condition.Message),
		}

		table.AddRow(row)
	}

	return []content.Content{&table}, nil
}

func retrievePod(object runtime.Object) (*core.Pod, error) {
	pod, ok := object.(*core.Pod)
	if !ok {
		return nil, errors.Errorf("expectect object to be a Pod, it was %T", object)
	}

	return pod, nil
}
