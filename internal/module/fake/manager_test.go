package fake

import (
	"github.com/stretchr/testify/assert"
	"github.com/twosson/kubeapt/internal/module"
	"testing"
)

func TestStubManager(t *testing.T) {
	m := NewModule("module")

	sm := NewStubManager("default", []module.Module{m})

	assert.Equal(t, []module.Module{m}, sm.Modules())
	assert.Equal(t, "default", sm.GetNamespace())

	sm.SetNamespace("other")
	assert.Equal(t, "other", sm.GetNamespace())
}
