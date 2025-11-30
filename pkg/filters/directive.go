package filters

import (
	"strings"

	"github.com/yeongyu/comment-checker/pkg/models"
)

// TypeCheckerPrefixes contains type checker and linter directive prefixes that should be skipped.
var TypeCheckerPrefixes = []string{
	"type:",
	"noqa",
	"pyright:",
	"ruff:",
	"mypy:",
	"pylint:",
	"flake8:",
	"pyre:",
	"pytype:",
	"eslint-disable",
	"eslint-ignore",
	"prettier-ignore",
	"ts-ignore",
	"ts-expect-error",
	"clippy:",
	"allow",
	"deny",
	"warn",
	"forbid",
}

// DirectiveFilter filters type checker and linter directives.
type DirectiveFilter struct{}

// NewDirectiveFilter creates a new DirectiveFilter.
func NewDirectiveFilter() *DirectiveFilter {
	return &DirectiveFilter{}
}

// ShouldSkip returns true if the comment is a directive.
func (f *DirectiveFilter) ShouldSkip(comment models.CommentInfo) bool {
	normalized := strings.ToLower(strings.TrimSpace(comment.Text))

	// Remove comment prefix (#, //, /*, --)
	for _, prefix := range []string{"#", "//", "/*", "--"} {
		if strings.HasPrefix(normalized, prefix) {
			normalized = strings.TrimSpace(normalized[len(prefix):])
			break
		}
	}

	// Remove @ symbol (TypeScript directives like @ts-ignore)
	if strings.HasPrefix(normalized, "@") {
		normalized = strings.TrimSpace(normalized[1:])
	}

	// Check if starts with any directive prefix
	for _, directive := range TypeCheckerPrefixes {
		if strings.HasPrefix(normalized, directive) {
			return true
		}
	}

	return false
}
