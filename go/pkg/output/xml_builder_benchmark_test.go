package output

import (
	"fmt"
	"testing"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

func benchmarkComments(count int) []models.CommentInfo {
	comments := make([]models.CommentInfo, count)
	for i := range comments {
		comments[i] = models.CommentInfo{
			Text:        fmt.Sprintf("# comment %d", i),
			LineNumber:  i + 1,
			FilePath:    "sample.py",
			CommentType: models.CommentTypeLine,
		}
	}
	return comments
}

func BenchmarkBuildCommentsXMLManyComments(b *testing.B) {
	comments := benchmarkComments(200)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		xml := BuildCommentsXML(comments, "sample.py")
		if len(xml) == 0 {
			b.Fatal("expected XML")
		}
	}
}

func BenchmarkFormatHookMessageManyComments(b *testing.B) {
	comments := benchmarkComments(200)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		message := FormatHookMessage(comments, "")
		if len(message) == 0 {
			b.Fatal("expected message")
		}
	}
}
