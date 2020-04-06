package api

import (
	"encoding/json"
	"github.com/twosson/kubeapt/internal/overview"
	"log"
	"net/http"
)

type navigationResponse struct {
	Sections []*overview.Navigation `json:"sections,omitempty"`
}

type navigationsResponse struct {
	Navigation []navigationResponse `json:"navigation,omitempty"`
}

type navigation struct {
	overview overview.Interface
}

var _ http.Handler = (*navigation)(nil)

func newNavigation(o overview.Interface) *navigation {
	return &navigation{overview: o}
}

func (n *navigation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	overviewNav, err := n.overview.Navigation()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nr := &navigationResponse{Sections: []*overview.Navigation{overviewNav}}

	if err := json.NewEncoder(w).Encode(nr); err != nil {
		log.Printf("encoding navigation error: %v", err)
	}
}
