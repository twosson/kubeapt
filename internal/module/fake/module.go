package fake

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/twosson/kubeapt/internal/apt"
	"log"
	"net/http"
)

// Module is a fake module.
type Module struct {
	name string
}

// NewModule creates an instance of Module.
func NewModule(name string) *Module {
	return &Module{
		name: name,
	}
}

// Name is the name of the module.
func (m *Module) Name() string {
	return m.name
}

// ContentPath is the path to the module's content.
func (m *Module) ContentPath() string {
	return fmt.Sprintf("/%s", m.name)
}

// Handler is a HTTP handler for the module.
func (m *Module) Handler(prefix string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix(prefix).Subrouter()
	s.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "root")
	}))

	s.Handle("/nested", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, m.name)
	}))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("fake module path not found: %s", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	return router

}

// Navigation returns navigation entries for the module.
func (m *Module) Navigation(prefix string) (*apt.Navigation, error) {
	nav := &apt.Navigation{
		Path:  prefix,
		Title: m.name,
	}

	return nav, nil
}

// SetNamespace sets the current namespace.
func (m *Module) SetNamespace(namespace string) error {
	return nil
}

// Start doesn't do anything.
func (m *Module) Start() error {
	return nil
}

// Stop doesn't do anything.
func (m *Module) Stop() {
}
