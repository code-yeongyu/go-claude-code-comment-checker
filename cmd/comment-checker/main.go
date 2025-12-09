package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/core"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/filters"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/output"
	"github.com/spf13/cobra"
)

// ToolInput represents the tool_input field from JSON input.
type ToolInput struct {
	FilePath  string `json:"file_path"`
	Content   string `json:"content"`
	NewString string `json:"new_string"`
	OldString string `json:"old_string"`
	Edits     []struct {
		OldString string `json:"old_string"`
		NewString string `json:"new_string"`
	} `json:"edits"`
}

// HookInput represents the JSON input from Claude Code hooks.
type HookInput struct {
	SessionID      string    `json:"session_id"`
	ToolName       string    `json:"tool_name"`
	TranscriptPath string    `json:"transcript_path"`
	Cwd            string    `json:"cwd"`
	HookEventName  string    `json:"hook_event_name"`
	ToolInput      ToolInput `json:"tool_input"`
	ToolResponse   any       `json:"tool_response"`
}

const (
	exitPass  = 0
	exitBlock = 2
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "comment-checker",
		Short: "Check for problematic comments in source code",
		Long:  "A hook for Claude Code that detects and warns about comments and docstrings in source code.",
		Run:   run,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: Command execution failed")
		os.Exit(exitPass)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Read JSON from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: Failed to read stdin")
		os.Exit(exitPass)
		return
	}

	// Handle empty input
	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No input provided")
		os.Exit(exitPass)
		return
	}

	// Parse JSON
	var hookInput HookInput
	if err := json.Unmarshal(input, &hookInput); err != nil {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: Invalid input format")
		os.Exit(exitPass)
		return
	}

	// Get file path
	filePath := hookInput.ToolInput.FilePath
	if filePath == "" {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No file path provided")
		os.Exit(exitPass)
		return
	}

	// Check if file is a code file (supported extension)
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	if ext == "" {
		// Handle files like "Dockerfile"
		ext = strings.ToLower(filepath.Base(filePath))
	}

	registry := core.NewLanguageRegistry()
	if !registry.IsSupported(ext) {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: Non-code file")
		os.Exit(exitPass)
		return
	}

	// Detect comments based on tool type
	detector := core.NewCommentDetector()
	var comments []models.CommentInfo

	switch hookInput.ToolName {
	case "Edit":
		// For Edit: only detect NEW comments (in new_string but not in old_string)
		if hookInput.ToolInput.NewString == "" {
			fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No content to check")
			os.Exit(exitPass)
			return
		}
		comments = detectNewCommentsForEdit(
			detector,
			hookInput.ToolInput.OldString,
			hookInput.ToolInput.NewString,
			filePath,
		)
	case "MultiEdit":
		// For MultiEdit: aggregate new comments from all edits
		if len(hookInput.ToolInput.Edits) == 0 {
			fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No content to check")
			os.Exit(exitPass)
			return
		}
		for _, edit := range hookInput.ToolInput.Edits {
			if edit.NewString == "" {
				continue
			}
			editComments := detectNewCommentsForEdit(
				detector,
				edit.OldString,
				edit.NewString,
				filePath,
			)
			comments = append(comments, editComments...)
		}
	default:
		// For Write and others: check entire content
		content := getContentToCheck(hookInput)
		if content == "" {
			fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No content to check")
			os.Exit(exitPass)
			return
		}
		comments = detector.Detect(content, filePath, true)
	}

	// No comments found
	if len(comments) == 0 {
		fmt.Fprintln(os.Stderr, "[check-comments] Success: No problematic comments/docstrings found")
		os.Exit(exitPass)
		return
	}

	// Apply filter chain: BDD -> Directive -> Shebang
	filtered := applyFilters(comments)

	// No problematic comments after filtering
	if len(filtered) == 0 {
		fmt.Fprintln(os.Stderr, "[check-comments] Success: No problematic comments/docstrings found")
		os.Exit(exitPass)
		return
	}

	// Problematic comments found - output warning and exit with code 2
	message := output.FormatHookMessage(filtered)
	fmt.Fprint(os.Stderr, message)
	os.Exit(exitBlock)
}

// getContentToCheck extracts the content to check based on tool type.
func getContentToCheck(input HookInput) string {
	switch input.ToolName {
	case "Write":
		return input.ToolInput.Content
	case "Edit":
		return input.ToolInput.NewString
	case "MultiEdit":
		// Combine all new_string values from edits
		var parts []string
		for _, edit := range input.ToolInput.Edits {
			if edit.NewString != "" {
				parts = append(parts, edit.NewString)
			}
		}
		return strings.Join(parts, "\n")
	default:
		// Unknown tool type, try content first, then new_string
		if input.ToolInput.Content != "" {
			return input.ToolInput.Content
		}
		return input.ToolInput.NewString
	}
}

// applyFilters applies all filters in order and returns remaining comments.
func applyFilters(comments []models.CommentInfo) []models.CommentInfo {
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

// buildCommentTextSet creates a set of normalized comment texts for comparison.
func buildCommentTextSet(comments []models.CommentInfo) map[string]struct{} {
	set := make(map[string]struct{}, len(comments))
	for _, c := range comments {
		set[c.NormalizedText()] = struct{}{}
	}
	return set
}

// filterNewComments returns comments that exist in newComments but not in oldComments.
func filterNewComments(oldComments, newComments []models.CommentInfo) []models.CommentInfo {
	if len(oldComments) == 0 {
		return newComments
	}

	oldSet := buildCommentTextSet(oldComments)

	var newOnly []models.CommentInfo
	for _, c := range newComments {
		if _, exists := oldSet[c.NormalizedText()]; !exists {
			newOnly = append(newOnly, c)
		}
	}
	return newOnly
}

// detectNewCommentsForEdit detects comments that are newly added in Edit operation.
func detectNewCommentsForEdit(detector *core.CommentDetector, oldString, newString, filePath string) []models.CommentInfo {
	oldComments := detector.Detect(oldString, filePath, true)
	newComments := detector.Detect(newString, filePath, true)

	return filterNewComments(oldComments, newComments)
}
