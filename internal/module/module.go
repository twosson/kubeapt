package module

import (
	"github.com/twosson/kubeapt/internal/apt"
	"net/http"
)

// Module is an apt plugin.
type Module interface {
	ContentPath() string
	Handler(root string) http.Handler
	Navigation(root string) (*apt.Navigation, error)
	Start() error
	Stop()
}
