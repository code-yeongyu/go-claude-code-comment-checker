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

func TestNormalizedText_GivenUppercaseComment_ReturnsLowercase(t *testing.T) {
	// given
	comment := CommentInfo{
		Text:        "// TODO: FIX THIS",
		LineNumber:  10,
		FilePath:    "main.go",
		CommentType: CommentTypeLine,
	}

	// when
	result := comment.NormalizedText()

	// then
	assert.Equal(t, "// todo: fix this", result)
}

func TestCommentInfo_GivenDocstring_HasIsDocstringTrue(t *testing.T) {
	// given
	comment := CommentInfo{
		Text:        `"""This is a docstring."""`,
		LineNumber:  5,
		FilePath:    "module.py",
		CommentType: CommentTypeDocstring,
		IsDocstring: true,
	}

	// when & then
	assert.True(t, comment.IsDocstring)
	assert.Equal(t, CommentTypeDocstring, comment.CommentType)
}

func TestCommentInfo_GivenMetadata_StoresMetadata(t *testing.T) {
	// given
	metadata := map[string]string{
		"language":  "python",
		"node_type": "comment",
	}
	comment := CommentInfo{
		Text:        "# comment",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: CommentTypeLine,
		Metadata:    metadata,
	}

	// when & then
	assert.Equal(t, "python", comment.Metadata["language"])
	assert.Equal(t, "comment", comment.Metadata["node_type"])
}
