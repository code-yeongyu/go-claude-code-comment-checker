package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeongyu/comment-checker/pkg/models"
)

func Test_ShouldSkip_GivenKeyword_ReturnsTrue(t *testing.T) {
	// given
	filter := NewBDDFilter()
	comment := models.CommentInfo{
		Text:        "# given",
		LineNumber:  1,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func Test_ShouldSkip_WhenThenKeyword_ReturnsTrue(t *testing.T) {
	// given
	filter := NewBDDFilter()
	comment := models.CommentInfo{
		Text:        "// when & then",
		LineNumber:  1,
		FilePath:    "test.js",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func Test_ShouldSkip_ArrangeActAssert_ReturnsTrue(t *testing.T) {
	// given
	filter := NewBDDFilter()
	comment := models.CommentInfo{
		Text:        "-- arrange",
		LineNumber:  1,
		FilePath:    "test.sql",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func Test_ShouldSkip_RegularComment_ReturnsFalse(t *testing.T) {
	// given
	filter := NewBDDFilter()
	comment := models.CommentInfo{
		Text:        "# This is a regular comment",
		LineNumber:  5,
		FilePath:    "test.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.False(t, result)
}

func Test_ShouldSkip_WhenAmpersandThen_ReturnsTrue(t *testing.T) {
	// given
	filter := NewBDDFilter()
	comment := models.CommentInfo{
		Text:        "// when&then",
		LineNumber:  1,
		FilePath:    "test.js",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}
