package output

import (
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

// FormatHookMessage formats comment detection results for Claude Code hooks.
// Groups comments by file path and builds complete error message with
// instructions and XML blocks for each file.
// Returns formatted hook error message, or empty string if no comments provided.
func FormatHookMessage(comments []models.CommentInfo) string {
	if len(comments) == 0 {
		return ""
	}

	// Group comments by file path
	byFile := make(map[string][]models.CommentInfo)
	fileOrder := make([]string, 0)
	for _, comment := range comments {
		if _, exists := byFile[comment.FilePath]; !exists {
			fileOrder = append(fileOrder, comment.FilePath)
		}
		byFile[comment.FilePath] = append(byFile[comment.FilePath], comment)
	}

	var sb strings.Builder

	// Header
	sb.WriteString("COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED\n\n")

	// Guidelines
	sb.WriteString("Your recent changes contain comments or docstrings, which triggered this hook.\n")
	sb.WriteString("You need to take immediate action. You must follow the conditions below.\n")
	sb.WriteString("(Listed in priority order - you must always act according to this priority order)\n\n")

	sb.WriteString("CRITICAL WARNING: This hook message MUST NEVER be ignored, even if you receive it multiple times.\n")
	sb.WriteString("You MUST provide corresponding explanation or action for EACH occurrence of this message.\n")
	sb.WriteString("Ignoring this message or failing to respond appropriately is strictly prohibited.\n\n")

	sb.WriteString("PRIORITY-BASED ACTION GUIDELINES:\n\n")

	sb.WriteString("1. This is a comment/docstring that already existed before\n")
	sb.WriteString("\t-> Explain to the user that this is an existing comment/docstring and proceed (justify it)\n\n")

	sb.WriteString("2. This is a newly written comment: but it's in given, when, then format\n")
	sb.WriteString("\t-> Tell the user it's a BDD comment and proceed (justify it)\n")
	sb.WriteString("\t-> Note: This applies to comments only, not docstrings\n\n")

	sb.WriteString("3. This is a newly written comment/docstring: but it's a necessary comment/docstring\n")
	sb.WriteString("\t-> Tell the user why this comment/docstring is absolutely necessary and proceed (justify it)\n")
	sb.WriteString("\t-> Examples of necessary comments: complex algorithms, security-related, performance optimization, regex, mathematical formulas\n")
	sb.WriteString("\t-> Examples of necessary docstrings: public API documentation, complex module/class interfaces\n")
	sb.WriteString("\t-> IMPORTANT: Most docstrings are unnecessary if the code is self-explanatory. Only keep truly essential ones.\n\n")

	sb.WriteString("4. This is a newly written comment/docstring: but it's an unnecessary comment/docstring\n")
	sb.WriteString("\t-> Apologize to the user and remove the comment/docstring.\n")
	sb.WriteString("\t-> Make the code itself clearer so it can be understood without comments/docstrings.\n")
	sb.WriteString("\t-> For verbose docstrings: refactor code to be self-documenting instead of adding lengthy explanations.\n\n")

	sb.WriteString("MANDATORY REQUIREMENT: You must acknowledge this hook message and take one of the above actions.\n")
	sb.WriteString("Review in the above priority order and take the corresponding action EVERY TIME this appears.\n\n")

	sb.WriteString("Detected comments/docstrings:\n")

	// XML blocks for each file (preserve order)
	for _, filePath := range fileOrder {
		fileComments := byFile[filePath]
		sb.WriteString(BuildCommentsXML(fileComments, filePath))
		sb.WriteString("\n")
	}

	return sb.String()
}
