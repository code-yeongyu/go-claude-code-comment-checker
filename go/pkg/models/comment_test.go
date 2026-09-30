package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedText_GivenCommentWithWhitespace_ReturnsStrippedLowercase(t *testing.T) {
	// given
	comment := CommentInfo{
		Text:        "  # Hello World  ",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: CommentTypeLine,
	}

	// when
	result := comment.NormalizedText()

	// then
	assert.Equal(t, "# hello world", result)
}
