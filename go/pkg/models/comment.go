package models

import "strings"

// CommentType represents the type of comment.
type CommentType string

const (
	CommentTypeLine      CommentType = "line"
	CommentTypeBlock     CommentType = "block"
	CommentTypeDocstring CommentType = "docstring"
)

// CommentInfo holds information about a single comment in source code.
type CommentInfo struct {
	Text        string            `json:"text"`
	LineNumber  int               `json:"line_number"`
	FilePath    string            `json:"file_path"`
	CommentType CommentType       `json:"comment_type"`
	IsDocstring bool              `json:"is_docstring"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// NormalizedText returns the comment text stripped of whitespace and lowercased.
func (c *CommentInfo) NormalizedText() string {
	return strings.ToLower(strings.TrimSpace(c.Text))
}
