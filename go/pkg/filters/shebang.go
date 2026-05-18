package filters

import (
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

// ShebangFilter filters shebang lines.
type ShebangFilter struct{}

// NewShebangFilter creates a new ShebangFilter.
func NewShebangFilter() *ShebangFilter {
	return &ShebangFilter{}
}

// ShouldSkip returns true if the comment is a shebang.
func (f *ShebangFilter) ShouldSkip(comment models.CommentInfo) bool {
	stripped := strings.TrimSpace(comment.Text)
	return strings.HasPrefix(stripped, "#!")
}
