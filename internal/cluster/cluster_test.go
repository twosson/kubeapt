package cluster

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func Test_FromKubeConfig(t *testing.T) {
	kubeconfig := filepath.Join("testdata", "kubeconfig.yaml")
	_, err := FromKubeconfig(kubeconfig)
	require.NoError(t, err)
}
