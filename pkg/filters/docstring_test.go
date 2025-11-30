package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeongyu/comment-checker/pkg/models"
)

func TestDocstringFilter_ShouldSkip_PythonDocstring_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDocstringFilter()
	comment := models.CommentInfo{
		Text:        `"""This is a docstring."""`,
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: models.CommentTypeDocstring,
		IsDocstring: true,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDocstringFilter_ShouldSkip_JSDocComment_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDocstringFilter()
	comment := models.CommentInfo{
		Text:        "/** This is a JSDoc comment */",
		LineNumber:  1,
		FilePath:    "test.js",
		CommentType: models.CommentTypeBlock,
		IsDocstring: false,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDocstringFilter_ShouldSkip_RegularComment_ReturnsFalse(t *testing.T) {
	// given
	filter := NewDocstringFilter()
	comment := models.CommentInfo{
		Text:        "// This is a regular comment",
		LineNumber:  1,
		FilePath:    "test.js",
		CommentType: models.CommentTypeLine,
		IsDocstring: false,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.False(t, result)
}
