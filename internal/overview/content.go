package overview

import (
	"github.com/twosson/kubeapt/internal/content"
)

type ContentResponse struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}

var emptyContentResponse = ContentResponse{}
