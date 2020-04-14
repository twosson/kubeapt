package overview

import (
	"github.com/twosson/kubeapt/internal/content"
)

type ContentResponse struct {
	Views []Content `json:"views,omitempty"`
}

var emptyContentResponse = ContentResponse{}

type Content struct {
	Contents []content.Content `json:"contents,omitempty"`
	Title    string            `json:"title,omitempty"`
}
