package output

import (
	"fmt"
	"strings"

	"github.com/yeongyu/comment-checker/pkg/models"
)

// BuildCommentsXML builds <comments> XML block for a given file and its comments.
// Returns XML formatted string with comments, or empty string if no comments provided.
func BuildCommentsXML(comments []models.CommentInfo, filePath string) string {
	if len(comments) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<comments file=\"%s\">\n", filePath))
	for _, comment := range comments {
		sb.WriteString(fmt.Sprintf("\t<comment line-number=\"%d\">%s</comment>\n", comment.LineNumber, comment.Text))
	}
	sb.WriteString("</comments>")

	return sb.String()
}
