package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeongyu/comment-checker/pkg/models"
)

func TestShebangFilter_ShouldSkip_ShebangLine_ReturnsTrue(t *testing.T) {
	// given
	filter := NewShebangFilter()
	comment := models.CommentInfo{
		Text:        "#!/usr/bin/env python",
		LineNumber:  1,
		FilePath:    "script.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestShebangFilter_ShouldSkip_ShebangWithSpace_ReturnsTrue(t *testing.T) {
	// given
	filter := NewShebangFilter()
	comment := models.CommentInfo{
		Text:        "#! /usr/bin/env python",
		LineNumber:  1,
		FilePath:    "script.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestShebangFilter_ShouldSkip_ShebangNode_ReturnsTrue(t *testing.T) {
	// given
	filter := NewShebangFilter()
	comment := models.CommentInfo{
		Text:        "#!/bin/bash",
		LineNumber:  1,
		FilePath:    "script.sh",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.True(t, result)
}

func TestShebangFilter_ShouldSkip_RegularComment_ReturnsFalse(t *testing.T) {
	// given
	filter := NewShebangFilter()
	comment := models.CommentInfo{
		Text:        "# Regular comment",
		LineNumber:  5,
		FilePath:    "script.py",
		CommentType: models.CommentTypeLine,
	}

	// when
	result := filter.ShouldSkip(comment)

	// then
	assert.False(t, result)
}
