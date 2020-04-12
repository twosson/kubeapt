package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heptio/go-telemetry/pkg/telemetry"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/module"
	"log"
	"net/http"
	"path"
	"time"
)

// Service is an API service.
type Service interface {
	RegisterModule(module.Module) error
	Handler() *mux.Router
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Printf("encoding response error: %v", err)
	}
}

// API is the API for the dashboard client
type API struct {
	nsClient        cluster.NamespaceInterface
	moduleManager   module.ManagerInterface
	sections        []*apt.Navigation
	prefix          string
	telemetryClient telemetry.Interface

	modules map[string]http.Handler
}

func (a *API) telemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		msDuration := int64(time.Since(startTime) / time.Millisecond)
		go a.telemetryClient.With(telemetry.Labels{"endpoint": r.URL.Path, "client.useragent": r.Header.Get("User-Agent")}).SendEvent("dash.api", telemetry.Measurements{"count": 1, "duration": msDuration})
	})
}

// New creates an instance of API.
func New(prefix string, nsClient cluster.NamespaceInterface, moduleManager module.ManagerInterface, telemetryClient telemetry.Interface) *API {
	return &API{
		prefix:          prefix,
		nsClient:        nsClient,
		moduleManager:   moduleManager,
		modules:         make(map[string]http.Handler),
		telemetryClient: telemetryClient,
	}
}

// Handler returns a HTTP handler for the service.
func (a *API) Handler() *mux.Router {
	router := mux.NewRouter()
	router.Use(a.telemetryMiddleware)
	s := router.PathPrefix(a.prefix).Subrouter()

	namespacesService := newNamespaces(a.nsClient)
	s.Handle("/namespaces", namespacesService).Methods(http.MethodGet)

	navigationService := newNavigation(a.sections)
	s.Handle("/navigation", navigationService).Methods(http.MethodGet)

	namespaceUpdateService := newNamespace(a.moduleManager)
	s.HandleFunc("/namespace", namespaceUpdateService.update).Methods(http.MethodPost)
	s.HandleFunc("/namespace", namespaceUpdateService.read).Methods(http.MethodGet)

	for p, h := range a.modules {
		s.PathPrefix(p).Handler(h)
	}

	s.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("api handler not found: %s", r.URL.String())
		respondWithError(w, http.StatusNotFound, "not found")
	})

	return router
}

// RegisterModule registers a module with the API service.
func (a *API) RegisterModule(m module.Module) error {
	contentPath := path.Join("/content", m.ContentPath())
	log.Printf("Registering content path %s", contentPath)
	a.modules[contentPath] = m.Handler(path.Join(a.prefix, contentPath))

	nav, err := m.Navigation(contentPath)
	if err != nil {
		return err
	}

	a.sections = append(a.sections, nav)

	return nil
}
