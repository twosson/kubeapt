package module

import (
	"github.com/twosson/kubeapt/internal/apt"
	"net/http"
)

// Module is an apt plugin.
type Module interface {
	Name() string
	ContentPath() string
	Handler(root string) http.Handler
	Navigation(root string) (*apt.Navigation, error)
	SetNamespace(namespace string) error
	Start() error
	Stop()
}
