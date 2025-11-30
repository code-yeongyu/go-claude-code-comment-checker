package filters

import (
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

// DocstringFilter filters docstrings (Python, JSDoc, Javadoc).
type DocstringFilter struct{}

// NewDocstringFilter creates a new DocstringFilter.
func NewDocstringFilter() *DocstringFilter {
	return &DocstringFilter{}
}

// ShouldSkip returns true if the comment is a docstring.
func (f *DocstringFilter) ShouldSkip(comment models.CommentInfo) bool {
	// Check is_docstring flag
	if comment.IsDocstring {
		return true
	}

	// Check JSDoc/Javadoc pattern: /** ... */
	if strings.HasPrefix(strings.TrimSpace(comment.Text), "/**") {
		return true
	}

	return false
}
