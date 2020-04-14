package module

import (
	"github.com/pkg/errors"
	"github.com/twosson/kubeapt/internal/cluster"
	"github.com/twosson/kubeapt/internal/log"
	"github.com/twosson/kubeapt/internal/overview"
)

// ManagerInterface is an interface for managing module lifecycle.
type ManagerInterface interface {
	Modules() []Module
	SetNamespace(namespace string)
	GetNamespace() string
}

// Manager manages module lifecycle.
type Manager struct {
	clusterClient cluster.ClientInterface
	namespace     string
	logger        log.Logger
	loadedModules []Module
}

var _ ManagerInterface = (*Manager)(nil)

// NewManager creates an instance of Manager.
func NewManager(clusterClient cluster.ClientInterface, namespace string, logger log.Logger) (*Manager, error) {
	manager := &Manager{
		clusterClient: clusterClient,
		namespace:     namespace,
		logger:        logger,
	}

	if err := manager.Load(); err != nil {
		return nil, err
	}

	return manager, nil
}

// Load loads modules.
func (m *Manager) Load() error {
	overviewModule, err := overview.NewClusterOverview(m.clusterClient, m.namespace, m.logger)
	if err != nil {
		return errors.Wrap(err, "loading overview module")
	}

	modules := []Module{overviewModule}

	for _, module := range modules {
		if err := module.Start(); err != nil {
			return errors.Wrapf(err, "%s module failed to start", module.Name())
		}
	}

	m.loadedModules = modules

	return nil
}

// Modules returns a list of modules.
func (m *Manager) Modules() []Module {
	return m.loadedModules
}

// Unload unloads modules.
func (m *Manager) Unload() {
	for _, module := range m.loadedModules {
		module.Stop()
	}
}

// SetNamespace sets the current namespace.
func (m *Manager) SetNamespace(namespace string) {
	m.namespace = namespace
	for _, module := range m.loadedModules {
		if err := module.SetNamespace(namespace); err != nil {
			m.logger.Errorf("setting namespace for module %q: %v", module.Name(), err)
		}
	}
}

// GetNamespace gets the current namespace.
func (m *Manager) GetNamespace() string {
	return m.namespace
}
