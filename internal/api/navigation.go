package api

import (
	"encoding/json"
	"github.com/twosson/kubeapt/internal/apt"
	"github.com/twosson/kubeapt/internal/log"
	"net/http"
)

type navigationResponse struct {
	Sections []*apt.Navigation `json:"sections,omitempty"`
}

type navigation struct {
	sections []*apt.Navigation
	logger   log.Logger
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(sections []*apt.Navigation, logger log.Logger) *navigation {
	return &navigation{
		sections: sections,
		logger:   logger,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := navigationResponse{
		Sections: n.sections,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		n.logger.Errorf("encoding navigation error: %v", err)
	}
}
