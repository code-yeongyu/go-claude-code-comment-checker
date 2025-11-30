package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

func Test_Detect_PythonLineComment_ReturnsCommentInfo(t *testing.T) {
	// given
	detector := NewCommentDetector()
	code := "# This is a comment\nprint('hello')"

	// when
	comments := detector.Detect(code, "test.py", false)

	// then
	assert.Len(t, comments, 1)
	assert.Equal(t, "# This is a comment", comments[0].Text)
	assert.Equal(t, 1, comments[0].LineNumber)
	assert.Equal(t, "test.py", comments[0].FilePath)
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
	assert.False(t, comments[0].IsDocstring)
}

func Test_Detect_TypeScriptBlockComment_ReturnsCommentInfo(t *testing.T) {
	// given
	detector := NewCommentDetector()
	code := `/* This is a block comment */
const x = 1;`

	// when
	comments := detector.Detect(code, "test.ts", false)

	// then
	assert.Len(t, comments, 1)
	assert.Equal(t, "/* This is a block comment */", comments[0].Text)
	assert.Equal(t, 1, comments[0].LineNumber)
	assert.Equal(t, "test.ts", comments[0].FilePath)
	assert.Equal(t, models.CommentTypeBlock, comments[0].CommentType)
	assert.False(t, comments[0].IsDocstring)
}

func Test_Detect_PythonDocstring_ReturnsDocstring(t *testing.T) {
	// given
	detector := NewCommentDetector()
	code := `"""This is a module docstring."""
def hello():
    pass`

	// when
	comments := detector.Detect(code, "module.py", true)

	// then
	assert.NotEmpty(t, comments)
	// Find the docstring in results
	var foundDocstring bool
	for _, c := range comments {
		if c.IsDocstring {
			foundDocstring = true
			assert.Equal(t, models.CommentTypeDocstring, c.CommentType)
			assert.Contains(t, c.Text, "module docstring")
			break
		}
	}
	assert.True(t, foundDocstring, "Expected to find a docstring")
}

func Test_Detect_UnsupportedExtension_ReturnsEmptyList(t *testing.T) {
	// given
	detector := NewCommentDetector()
	code := "some random content"

	// when
	comments := detector.Detect(code, "test.xyz", false)

	// then
	assert.Empty(t, comments)
}

func Test_Detect_GoComment_ReturnsCommentInfo(t *testing.T) {
	// given
	detector := NewCommentDetector()
	code := `// This is a Go comment
package main

func main() {}`

	// when
	comments := detector.Detect(code, "main.go", false)

	// then
	assert.Len(t, comments, 1)
	assert.Equal(t, "// This is a Go comment", comments[0].Text)
	assert.Equal(t, 1, comments[0].LineNumber)
	assert.Equal(t, "main.go", comments[0].FilePath)
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
	assert.False(t, comments[0].IsDocstring)
}
