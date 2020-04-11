package overview

import (
	"github.com/twosson/kubeapt/internal/content"
)

type contentResponse struct {
	Contents []content.Content `json:"contents,omitempty"`
}
