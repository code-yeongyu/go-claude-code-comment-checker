package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yeongyu/comment-checker/pkg/core"
	"github.com/yeongyu/comment-checker/pkg/filters"
	"github.com/yeongyu/comment-checker/pkg/models"
	"github.com/yeongyu/comment-checker/pkg/output"
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

	// Get content to check based on tool type
	content := getContentToCheck(hookInput)
	if content == "" {
		fmt.Fprintln(os.Stderr, "[check-comments] Skipping: No content to check")
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

	// Detect comments
	detector := core.NewCommentDetector()
	comments := detector.Detect(content, filePath, true)

	// No comments found
	if len(comments) == 0 {
		fmt.Fprintln(os.Stderr, "[check-comments] Success: No problematic comments/docstrings found")
		os.Exit(exitPass)
		return
	}

	// Apply filter chain: BDD -> Docstring -> Directive -> Shebang
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
	docstringFilter := filters.NewDocstringFilter()
	directiveFilter := filters.NewDirectiveFilter()
	shebangFilter := filters.NewShebangFilter()

	var filtered []models.CommentInfo
	for _, c := range comments {
		if bddFilter.ShouldSkip(c) {
			continue
		}
		if docstringFilter.ShouldSkip(c) {
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
