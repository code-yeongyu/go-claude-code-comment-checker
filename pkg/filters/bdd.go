package filters

import (
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

// BDDKeywords contains BDD-style comment keywords that should be skipped.
var BDDKeywords = map[string]struct{}{
	"given":       {},
	"when":        {},
	"then":        {},
	"arrange":     {},
	"act":         {},
	"assert":      {},
	"when & then": {},
	"when&then":   {},
}

// BDDFilter filters BDD-style comments.
type BDDFilter struct{}

// NewBDDFilter creates a new BDDFilter.
func NewBDDFilter() *BDDFilter {
	return &BDDFilter{}
}

// ShouldSkip returns true if the comment is a BDD keyword.
func (f *BDDFilter) ShouldSkip(comment models.CommentInfo) bool {
	normalized := strings.ToLower(strings.TrimSpace(comment.Text))

	// Remove comment prefix (#, //, --)
	for _, prefix := range []string{"#", "//", "--"} {
		if strings.HasPrefix(normalized, prefix) {
			normalized = strings.TrimSpace(normalized[len(prefix):])
			break
		}
	}

	_, exists := BDDKeywords[normalized]
	return exists
}
