package filters

import (
	"regexp"
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

var agentMemoPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^[\s#/*-]*changed?\s+(from|to)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*modified?\s+(from|to)?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*updated?\s+(from|to)?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*refactor(ed|ing)?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*moved?\s+(from|to)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*renamed?\s+(from|to)?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*replaced?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*removed?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*deleted?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*added?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*implemented?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*this\s+(implements?|adds?|removes?|changes?|fixes?)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*here\s+we\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*now\s+(we|this|it)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*previously\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*before\s+this\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*after\s+this\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*was\s+changed\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*implementation\s+(of|note)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*note:\s*\w`),
	regexp.MustCompile(`(?i)^[\s#/*-]*[a-z]+\s*->\s*[a-z]+`),
	regexp.MustCompile(`(?i)^[\s#/*-]*converted?\s+(from|to)\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*migrated?\s+(from|to)?\b`),
	regexp.MustCompile(`(?i)^[\s#/*-]*switched?\s+(from|to)\b`),
}

type AgentMemoFilter struct{}

func NewAgentMemoFilter() *AgentMemoFilter {
	return &AgentMemoFilter{}
}

func (f *AgentMemoFilter) IsAgentMemo(comment models.CommentInfo) bool {
	text := strings.TrimSpace(comment.Text)

	for _, prefix := range []string{"#", "//", "/*", "--", "*"} {
		if strings.HasPrefix(text, prefix) {
			text = strings.TrimSpace(text[len(prefix):])
		}
	}

	for _, pattern := range agentMemoPatterns {
		if pattern.MatchString(text) {
			return true
		}
	}

	return false
}
