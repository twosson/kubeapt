package overview

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
	"testing"
	"time"
)

func TestSectionDescriber(t *testing.T) {
	namespace := "default"

	d := NewSectionDescriber(
		newStubDescriber(),
	)

	cache := NewMemoryCache()

	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	assert.Equal(t, stubbedContent, got)
}

func TestCronJobsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "cronjob.yaml")

	d := NewCronJobsDescriber()
	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, tbl.Title, "Cron Jobs")
	assert.Len(t, tbl.Rows, 1)
}

func TestCronJobDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "cronjob.yaml")
	loadUnstructured(t, cache, namespace, "event-1.yaml")

	fields := map[string]string{
		"name": "hello",
	}

	d := NewCronJobDescriber()
	got, err := d.Describe("/prefix", namespace, cache, fields)
	require.NoError(t, err)

	require.Len(t, got, 2)
	cjTable, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, cjTable.Title, "Cron Job")
	assert.Len(t, cjTable.Rows, 1)

	eventsTable, ok := got[1].(table)
	require.True(t, ok)

	assert.Equal(t, eventsTable.Title, "Events")
	assert.Len(t, eventsTable.Rows, 1)
}

func Test_printCronJob(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		path     string
		expected tableRow
	}{
		{
			name: "not scheduled",
			path: "cronjob.yaml",
			expected: tableRow{
				"active":        newStringText("0"),
				"age":           newStringText("<unknown>"),
				"labels":        newLabelsText(nil),
				"last_schedule": newStringText("<none>"),
				"name":          newLinkText("hello", "/api/workloads/cron-jobs/hello?namespace=default"),
				"schedule":      newStringText("*/1 * * * *"),
			},
		},
		{
			name: "scheduled",
			path: "cronjob-scheduled.yaml",
			expected: tableRow{
				"active":        newStringText("0"),
				"age":           newStringText("<unknown>"),
				"labels":        newLabelsText(nil),
				"last_schedule": newStringText("30s"),
				"name":          newLinkText("hello", "/api/workloads/cron-jobs/hello?namespace=default"),
				"schedule":      newStringText("*/1 * * * *"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cj, ok := loadType(t, tc.path).(*batchv1beta1.CronJob)
			require.True(t, ok)

			got := printCronJob(cj, "/api", "default", c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestEventsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "event-1.yaml")
	loadUnstructured(t, cache, namespace, "event-2.yaml")

	d := NewEventsDescriber()
	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, tbl.Title, "Events")
	assert.Len(t, tbl.Rows, 2)
}

func Test_printEvent(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		path     string
		expected tableRow
	}{
		{
			name: "event",
			path: "event-1.yaml",
			expected: tableRow{
				"message":    newStringText("(combined from similar events): Saw completed job: hello-1538868300"),
				"source":     newStringText("cronjob-controller"),
				"sub_object": newStringText(""),
				"count":      newStringText("24973"),
				"first_seen": newStringText("2018-09-18T08:40:18-04:00"),
				"last_seen":  newStringText("2018-10-06T19:25:55-04:00"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event, ok := loadType(t, tc.path).(*corev1.Event)
			require.True(t, ok)

			got := printEvent(event, "/api", "default", c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_translateTimestamp(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "zero",
			expected: "<unknown>",
		},
		{
			name:     "not zero",
			time:     time.Unix(1538828100, 0),
			expected: "30s",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := metav1.NewTime(tc.time)

			got := translateTimestamp(ts, c)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func loadType(t *testing.T, path string) runtime.Object {
	data, err := ioutil.ReadFile(filepath.Join("testdata", path))
	require.NoError(t, err)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(data, nil, nil)
	require.NoError(t, err)

	return obj
}

func loadUnstructured(t *testing.T, cache Cache, namespace, path string) {
	obj := loadType(t, path)
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	require.NoError(t, err)

	u := &unstructured.Unstructured{
		Object: m,
	}
	u.Object = m
	u.SetNamespace(namespace)

	err = cache.Store(u)
	require.NoError(t, err)
}
