package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

func Test_FormatHookMessage_SingleComment_ReturnsFormattedMessage(t *testing.T) {
	// given
	comments := []models.CommentInfo{
		{
			Text:        "# TODO: fix this",
			LineNumber:  10,
			FilePath:    "src/app.py",
			CommentType: models.CommentTypeLine,
		},
	}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Contains(t, result, "COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED")
	assert.Contains(t, result, `<comments file="src/app.py">`)
	assert.Contains(t, result, `<comment line-number="10"># TODO: fix this</comment>`)
	assert.Contains(t, result, "</comments>")
	assert.Contains(t, result, "PRIORITY-BASED ACTION GUIDELINES:")
}

func Test_FormatHookMessage_MultipleCommentsSingleFile_ReturnsGroupedXML(t *testing.T) {
	// given
	comments := []models.CommentInfo{
		{
			Text:        "# First comment",
			LineNumber:  5,
			FilePath:    "src/main.py",
			CommentType: models.CommentTypeLine,
		},
		{
			Text:        "# Second comment",
			LineNumber:  15,
			FilePath:    "src/main.py",
			CommentType: models.CommentTypeLine,
		},
	}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Contains(t, result, `<comments file="src/main.py">`)
	assert.Contains(t, result, `<comment line-number="5"># First comment</comment>`)
	assert.Contains(t, result, `<comment line-number="15"># Second comment</comment>`)
	assert.Equal(t, 1, strings.Count(result, `<comments file="src/main.py">`))
	assert.Equal(t, 1, strings.Count(result, "</comments>"))
}

func Test_FormatHookMessage_CommentsMultipleFiles_ReturnsSeparateXMLBlocks(t *testing.T) {
	// given
	comments := []models.CommentInfo{
		{
			Text:        "# Comment in file 1",
			LineNumber:  10,
			FilePath:    "src/file1.py",
			CommentType: models.CommentTypeLine,
		},
		{
			Text:        "# Comment in file 2",
			LineNumber:  20,
			FilePath:    "src/file2.py",
			CommentType: models.CommentTypeLine,
		},
	}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Contains(t, result, `<comments file="src/file1.py">`)
	assert.Contains(t, result, `<comments file="src/file2.py">`)
	assert.Contains(t, result, `<comment line-number="10"># Comment in file 1</comment>`)
	assert.Contains(t, result, `<comment line-number="20"># Comment in file 2</comment>`)
	assert.Equal(t, 2, strings.Count(result, "</comments>"))
}

func Test_FormatHookMessage_EmptyList_ReturnsEmptyString(t *testing.T) {
	// given
	comments := []models.CommentInfo{}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Equal(t, "", result)
}

func Test_FormatHookMessage_DocstringComment_ReturnsFormattedMessage(t *testing.T) {
	// given
	comments := []models.CommentInfo{
		{
			Text:        `"""Module docstring."""`,
			LineNumber:  1,
			FilePath:    "src/utils.py",
			CommentType: models.CommentTypeDocstring,
			IsDocstring: true,
		},
	}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Contains(t, result, "COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED")
	assert.Contains(t, result, `<comments file="src/utils.py">`)
	assert.Contains(t, result, `<comment line-number="1">"""Module docstring."""</comment>`)
	assert.Contains(t, result, "MANDATORY REQUIREMENT:")
}

func Test_FormatHookMessage_XMLUsesTabsForIndentation(t *testing.T) {
	// given
	comments := []models.CommentInfo{
		{
			Text:        "# Comment",
			LineNumber:  5,
			FilePath:    "src/test.py",
			CommentType: models.CommentTypeLine,
		},
	}

	// when
	result := FormatHookMessage(comments)

	// then
	assert.Contains(t, result, "\t<comment line-number=\"5\"># Comment</comment>")
}
