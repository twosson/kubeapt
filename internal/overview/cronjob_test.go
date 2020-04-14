package overview

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/apimachinery/pkg/util/clock"
	"testing"
	"time"
)

func TestCronJobSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewCronJobSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestCronJobSummary(t *testing.T) {
	s := NewCronJobSummary("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()
	cache := NewMemoryCache()

	cronJob := loadFromFile(t, "cronjob-1.yaml")
	cronJob = convertToInternal(t, cronJob)

	storeFromFile(t, "job-1.yaml", cache)

	contents, err := s.Content(ctx, cronJob, cache)
	require.NoError(t, err)

	sections := []content.Section{
		{
			Items: []content.Item{
				content.TextItem("Name", "hello"),
				content.TextItem("Namespace", "default"),
				content.LabelsItem("Labels", map[string]string{"overview": "default"}),
				content.LabelsItem("Annotations", map[string]string{}),
				content.TimeItem("Create Time", "2018-09-18T12:30:09Z"),
				content.TextItem("Active", "0"),
				content.TextItem("Schedule", "*/1 * * * *"),
				content.TextItem("Suspend", "false"),
				content.TimeItem("Last Schedule", "2018-11-02T09:45:00Z"),
				content.TextItem("Concurrency Policy", "Allow"),
				content.TextItem("Starting Deadline Seconds", "<unset>"),
			},
		},
	}
	details := content.NewSummary("Details", sections)

	expected := []content.Content{
		&details,
	}

	assert.Equal(t, expected, contents)
}

func TestCronJobJobs(t *testing.T) {
	cjj := NewCronJobJobs("prefix", "ns", clock.NewFakeClock(time.Now()))

	ctx := context.Background()
	cache := NewMemoryCache()

	cronJob := loadFromFile(t, "cronjob-1.yaml")
	cronJob = convertToInternal(t, cronJob)

	storeFromFile(t, "job-1.yaml", cache)

	contents, err := cjj.Content(ctx, cronJob, cache)
	require.NoError(t, err)

	jobColumns := tableCols("Name", "Desired", "Successful", "Age", "Containers",
		"Images", "Selector", "Labels")

	activeTable := content.NewTable("Active Jobs", "No active jobs")
	activeTable.Columns = jobColumns

	inactiveTable := content.NewTable("Inactive Jobs", "No inactive jobs")
	inactiveTable.Columns = jobColumns
	inactiveTable.AddRow(content.TableRow{
		"Age":        content.NewStringText("1d"),
		"Containers": content.NewStringText("hello"),
		"Desired":    content.NewStringText("1"),
		"Images":     content.NewStringText("busybox"),
		"Labels":     content.NewStringText("controller-uid=f20be17b-de8b-11e8-889a-025000000001,job-name=hello-1541155320"),
		"Name":       content.NewLinkText("hello-1541155320", "/content/overview/workloads/jobs/hello-1541155320"),
		"Selector":   content.NewStringText("controller-uid=f20be17b-de8b-11e8-889a-025000000001"),
		"Successful": content.NewStringText("1"),
	})

	expected := []content.Content{
		&activeTable,
		&inactiveTable,
	}

	assert.Equal(t, expected, contents)
}
