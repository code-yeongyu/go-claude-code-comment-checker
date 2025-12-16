package tests

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/core"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/filters"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// FULL PIPELINE TESTS
// ============================================================================

func Test_FullPipeline_NoComments_ReturnsEmpty(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `package main

func main() {
	println("hello")
}`

	// when
	comments := detector.Detect(code, "main.go", false)

	// then
	assert.Empty(t, comments)
}

func Test_FullPipeline_WithBDDComment_FiltersOut(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# given
def test_something():
    pass`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.NotEmpty(t, comments, "Should detect BDD comment")
	assert.Empty(t, filtered, "BDD comment should be filtered out")
}

func Test_FullPipeline_WithRegularComment_ReturnsFormatted(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# This is a regular comment
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)
	message := output.FormatHookMessage(filtered)

	// then
	assert.Len(t, filtered, 1)
	assert.Contains(t, message, "COMMENT/DOCSTRING DETECTED")
	assert.Contains(t, message, "test.py")
	assert.Contains(t, message, "This is a regular comment")
}

func Test_FullPipeline_WithDocstring_DetectsAsCodeSmell(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `"""This is a docstring"""
def hello():
    pass`

	// when
	comments := detector.Detect(code, "test.py", true)
	filtered := applyFilterChain(comments)

	// then
	assert.NotEmpty(t, comments, "Should detect docstring")
	assert.NotEmpty(t, filtered, "Docstring should NOT be filtered out (now a code smell)")
}

func Test_FullPipeline_WithDirective_FiltersOut(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# noqa: E501
print("very long line")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.NotEmpty(t, comments, "Should detect directive comment")
	assert.Empty(t, filtered, "Directive should be filtered out")
}

func Test_FullPipeline_WithShebang_FiltersOut(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `#!/usr/bin/env python
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.NotEmpty(t, comments, "Should detect shebang")
	assert.Empty(t, filtered, "Shebang should be filtered out")
}

// ============================================================================
// CLI SUBPROCESS TESTS
// ============================================================================

func getBinaryPath(t *testing.T) string {
	// Get the project root directory (one level up from tests/)
	projectRoot := filepath.Join("..", "")
	binaryPath := filepath.Join(projectRoot, "comment-checker")

	// Verify binary exists
	absPath, err := filepath.Abs(binaryPath)
	require.NoError(t, err, "Failed to get absolute path")

	return absPath
}

func Test_CLI_NoComment_ExitZero(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"print(1)"}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	output, err := cmd.CombinedOutput()

	// then
	assert.NoError(t, err, "Expected exit 0 for no comments")
	assert.Contains(t, string(output), "Success")
}

func Test_CLI_WithComment_ExitTwo(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"# comment\nprint(1)"}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	err := cmd.Run()

	// then
	if exitErr, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 2, exitErr.ExitCode(), "Expected exit code 2 for comment detected")
	} else {
		t.Fatalf("Expected ExitError with code 2, got: %v", err)
	}
}

func Test_CLI_NonCodeFile_ExitZeroSkip(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"Write","tool_input":{"file_path":"test.json","content":"{}"}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	output, err := cmd.CombinedOutput()

	// then
	assert.NoError(t, err, "Expected exit 0 for non-code file")
	assert.Contains(t, string(output), "Skipping")
}

func Test_CLI_InvalidJSON_ExitZeroSkip(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `invalid json`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	output, err := cmd.CombinedOutput()

	// then
	assert.NoError(t, err, "Expected exit 0 for invalid JSON")
	assert.Contains(t, string(output), "Skipping")
}

func Test_CLI_EditTool_WithComment_ExitTwo(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"Edit","tool_input":{"file_path":"test.py","old_string":"x","new_string":"# comment\ny"}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	err := cmd.Run()

	// then
	if exitErr, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 2, exitErr.ExitCode(), "Expected exit code 2 for comment in Edit")
	} else {
		t.Fatalf("Expected ExitError with code 2, got: %v", err)
	}
}

func Test_CLI_MultiEdit_WithComment_ExitTwo(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"MultiEdit","tool_input":{"file_path":"test.py","edits":[{"old_string":"a","new_string":"# comment"}]}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	err := cmd.Run()

	// then
	if exitErr, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, 2, exitErr.ExitCode(), "Expected exit code 2 for comment in MultiEdit")
	} else {
		t.Fatalf("Expected ExitError with code 2, got: %v", err)
	}
}

func Test_CLI_BDDComment_ExitZero(t *testing.T) {
	// given
	binaryPath := getBinaryPath(t)
	input := `{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"# given\nprint(1)"}}`

	cmd := exec.Command(binaryPath)
	cmd.Stdin = strings.NewReader(input)

	// when
	output, err := cmd.CombinedOutput()

	// then
	assert.NoError(t, err, "Expected exit 0 for BDD comment")
	assert.Contains(t, string(output), "Success")
}

// ============================================================================
// MULTI-LANGUAGE DETECTION TESTS
// ============================================================================

func Test_Detect_MultiLanguage_Python_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := "# Python comment\nprint('hello')"

	// when
	comments := detector.Detect(code, "test.py", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Python comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Go_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// Go comment
package main`

	// when
	comments := detector.Detect(code, "main.go", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Go comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_TypeScript_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// TypeScript comment
const x: number = 1;`

	// when
	comments := detector.Detect(code, "test.ts", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "TypeScript comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_JavaScript_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// JavaScript comment
const x = 1;`

	// when
	comments := detector.Detect(code, "test.js", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "JavaScript comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Java_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// Java comment
public class Test {}`

	// when
	comments := detector.Detect(code, "Test.java", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Java comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_C_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// C comment
int main() { return 0; }`

	// when
	comments := detector.Detect(code, "main.c", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "C comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Cpp_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// C++ comment
int main() { return 0; }`

	// when
	comments := detector.Detect(code, "main.cpp", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "C++ comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Rust_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// Rust comment
fn main() {}`

	// when
	comments := detector.Detect(code, "main.rs", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Rust comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Ruby_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# Ruby comment
puts "hello"`

	// when
	comments := detector.Detect(code, "test.rb", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Ruby comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Bash_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# Bash comment
echo "hello"`

	// when
	comments := detector.Detect(code, "script.sh", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Bash comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Kotlin_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// Kotlin comment
fun main() {}`

	// when
	comments := detector.Detect(code, "Main.kt", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Kotlin comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

func Test_Detect_MultiLanguage_Swift_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `// Swift comment
print("hello")`

	// when
	comments := detector.Detect(code, "main.swift", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "Swift comment")
	assert.Equal(t, models.CommentTypeLine, comments[0].CommentType)
}

// ============================================================================
// BLOCK COMMENT TESTS
// ============================================================================

func Test_Detect_BlockComment_JavaScript_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `/* block comment */
const x = 1;`

	// when
	comments := detector.Detect(code, "test.js", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "block comment")
	assert.Equal(t, models.CommentTypeBlock, comments[0].CommentType)
}

func Test_Detect_BlockComment_C_Works(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `/* C block comment */
int main() { return 0; }`

	// when
	comments := detector.Detect(code, "main.c", false)

	// then
	assert.Len(t, comments, 1)
	assert.Contains(t, comments[0].Text, "C block comment")
	assert.Equal(t, models.CommentTypeBlock, comments[0].CommentType)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func applyFilterChain(comments []models.CommentInfo) []models.CommentInfo {
	bddFilter := filters.NewBDDFilter()
	directiveFilter := filters.NewDirectiveFilter()
	shebangFilter := filters.NewShebangFilter()

	var filtered []models.CommentInfo
	for _, c := range comments {
		if bddFilter.ShouldSkip(c) {
			continue
		}
		if directiveFilter.ShouldSkip(c) {
			continue
		}
		if shebangFilter.ShouldSkip(c) {
			continue
		}
		filtered = append(filtered, c)
	}

	return filtered
}

func Test_FullPipeline_WithAgentMemo_DetectsAsCodeSmell(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	agentMemoFilter := filters.NewAgentMemoFilter()
	code := `# Changed from old_value to new_value
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.Len(t, filtered, 1)
	assert.True(t, agentMemoFilter.IsAgentMemo(filtered[0]))
}

func Test_FullPipeline_WithAgentMemo_FormatterIncludesWarning(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	code := `# Modified to use new implementation
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)
	message := output.FormatHookMessage(filtered)

	// then
	assert.Contains(t, message, "AGENT MEMO")
	assert.Contains(t, message, "CODE SMELL")
}

func Test_FullPipeline_WithAgentMemo_Korean_DetectsAsCodeSmell(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	agentMemoFilter := filters.NewAgentMemoFilter()
	code := `# 여기서 값이 변경됨
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.Len(t, filtered, 1)
	assert.True(t, agentMemoFilter.IsAgentMemo(filtered[0]))
}

func Test_FullPipeline_RegularComment_NotAgentMemo(t *testing.T) {
	// given
	detector := core.NewCommentDetector()
	agentMemoFilter := filters.NewAgentMemoFilter()
	code := `# Calculate the sum of values
print("hello")`

	// when
	comments := detector.Detect(code, "test.py", false)
	filtered := applyFilterChain(comments)

	// then
	assert.Len(t, filtered, 1)
	assert.False(t, agentMemoFilter.IsAgentMemo(filtered[0]))
}
