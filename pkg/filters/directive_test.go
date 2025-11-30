package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeongyu/comment-checker/pkg/models"
)

func TestDirectiveFilter_ShouldSkip_NoqaDirective_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "# noqa: F401",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDirectiveFilter_ShouldSkip_TsIgnoreDirective_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "// @ts-ignore",
		LineNumber:  1,
		FilePath:    "test.ts",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDirectiveFilter_ShouldSkip_PyrightDirective_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "# pyright: ignore",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDirectiveFilter_ShouldSkip_EslintDisable_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "// eslint-disable-next-line",
		LineNumber:  1,
		FilePath:    "test.js",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDirectiveFilter_ShouldSkip_TypeDirective_ReturnsTrue(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "# type: ignore",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestDirectiveFilter_ShouldSkip_RegularComment_ReturnsFalse(t *testing.T) {
	// given
	filter := NewDirectiveFilter()
	comment := models.CommentInfo{
		Text:        "# Regular comment",
		LineNumber:  5,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.False(t, result)
}
