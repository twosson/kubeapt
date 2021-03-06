package overview

import (
	"context"
	"fmt"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubernetes/pkg/apis/core"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type EventList struct {
	prefix    string
	namespace string
	clock     clock.Clock
}

func NewEventList(prefix, namespace string, cl clock.Clock) View {
	return &EventList{
		prefix:    prefix,
		namespace: namespace,
		clock:     cl,
	}
}

func (el *EventList) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}
	eventObjects, err := c.Events(&unstructured.Unstructured{Object: m})
	if err != nil {
		return nil, err
	}

	var events []*core.Event
	for _, obj := range eventObjects {
		event := &core.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	table, err := printEvents(el.prefix, el.namespace, object, events, el.clock)
	if err != nil {
		return nil, err
	}

	return []content.Content{&table}, nil
}

// EventsDescriber creates content for a list of events.
type EventsDescriber struct {
	*baseDescriber

	path      string
	title     string
	cacheKeys []CacheKey
}

// NewEventsDescriber creates an instance of EventsDescriber.
func NewEventsDescriber(p string) *EventsDescriber {
	return &EventsDescriber{
		baseDescriber: newBaseDescriber(),
		path:          p,
		title:         "Events",
		cacheKeys: []CacheKey{
			{
				APIVersion: "v1",
				Kind:       "Event",
			},
		},
	}
}

// Describe creates content.
func (d *EventsDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (ContentResponse, error) {
	objects, err := loadObjects(ctx, options.Cache, namespace, options.Fields, d.cacheKeys)
	if err != nil {
		return emptyContentResponse, err
	}

	var contents []content.Content

	t := newEventTable(d.title, nil)

	sort.Slice(objects, func(i, j int) bool {
		tsI := objects[i].GetCreationTimestamp()
		tsJ := objects[j].GetCreationTimestamp()

		return tsI.Before(&tsJ)
	})

	for _, object := range objects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, event)
		if err != nil {
			return emptyContentResponse, err
		}

		t.Rows = append(t.Rows, printEvent(event, prefix, namespace, d.clock()))
	}

	contents = append(contents, &t)

	return ContentResponse{
		Views: []Content{
			Content{Contents: contents, Title: "Events"},
		},
	}, nil
}

func (d *EventsDescriber) PathFilters(namespace string) []pathFilter {
	return []pathFilter{
		*newPathFilter(d.path, d),
	}
}

func newEventTable(namespace string, object runtime.Object) content.Table {
	emptyMessage := emptyEventsMessageForObject(namespace, object)
	t := content.NewTable("Events", emptyMessage)

	t.Columns = []content.TableColumn{
		{Name: "Message", Accessor: "message"},
		{Name: "Source", Accessor: "source"},
		{Name: "Sub-Object", Accessor: "sub_object"},
		{Name: "Count", Accessor: "count"},
		{Name: "First Seen", Accessor: "first_seen"},
		{Name: "Last Seen", Accessor: "last_seen"},
	}

	return t
}

func printEvent(event *corev1.Event, prefix, namespace string, c clock.Clock) content.TableRow {
	firstSeen := event.FirstTimestamp.UTC().Format(time.RFC3339)
	lastSeen := event.LastTimestamp.UTC().Format(time.RFC3339)

	return content.TableRow{
		"message":    content.NewStringText(event.Message),
		"source":     content.NewStringText(event.Source.Component),
		"sub_object": content.NewStringText(""), // TODO: where does this come from?
		"count":      content.NewStringText(fmt.Sprint(event.Count)),
		"first_seen": content.NewStringText(firstSeen),
		"last_seen":  content.NewStringText(lastSeen),
	}
}

func emptyEventsMessageForObject(namespace string, object runtime.Object) string {
	if object == nil {
		return fmt.Sprintf("Namespace %s does not contain any events", namespace)
	}
	objectKind := object.GetObjectKind()
	gvk := objectKind.GroupVersionKind()

	return fmt.Sprintf("Namespace %s does not contain any events for this %s",
		namespace, gvk.Kind)
}
