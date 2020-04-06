package api

import (
	"encoding/json"
	"github.com/twosson/kubeapt/internal/apt"
	"log"
	"net/http"
)

type navigationResponse struct {
	Sections []*apt.Navigation `json:"sections,omitempty"`
}

type navigationsResponse struct {
	Navigation []navigationResponse `json:"navigation,omitempty"`
}

type navigation struct {
	sections []*apt.Navigation
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(sections []*apt.Navigation) *navigation {
	return &navigation{
		sections: sections,
	}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nr := navigationResponse{
		Sections: n.sections,
	}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		log.Printf("encoding navigation error: %v", err)
	}
}
