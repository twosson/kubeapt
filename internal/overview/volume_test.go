package overview

import (
	"github.com/stretchr/testify/assert"
	"github.com/twosson/kubeapt/internal/content"
	"k8s.io/kubernetes/pkg/apis/core"
	"testing"
)

func Test_summarizePersistentVolumeClaimVolumeSource(t *testing.T) {
	claim := &core.PersistentVolumeClaimVolumeSource{
		ClaimName: "my-claim",
	}

	section := &content.Section{}

	summarizePersistentVolumeClaimVolumeSource(section, claim)

	expected := &content.Section{}
	expected.AddText("Type", "PersistentVolumeClaim")
	expected.AddLink("ClaimName", "my-claim", "/content/overview/config-and-storage/persistent-volume-claims/my-claim")
	expected.AddText("ReadOnly", "false")

	assert.Equal(t, expected, section)
}
