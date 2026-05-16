package output

import (
	"strconv"
	"strings"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

func BuildCommentsXML(comments []models.CommentInfo, filePath string) string {
	if len(comments) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.Grow(24 + len(filePath) + len(comments)*48)
	sb.WriteString("<comments file=\"")
	sb.WriteString(filePath)
	sb.WriteString("\">\n")
	for _, comment := range comments {
		sb.WriteString("\t<comment line-number=\"")
		sb.WriteString(strconv.Itoa(comment.LineNumber))
		sb.WriteString("\">")
		sb.WriteString(comment.Text)
		sb.WriteString("</comment>\n")
	}
	sb.WriteString("</comments>")

	return sb.String()
}
