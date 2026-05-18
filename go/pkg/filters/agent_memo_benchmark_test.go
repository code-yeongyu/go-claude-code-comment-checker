package filters

import (
	"testing"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

func BenchmarkAgentMemoEarlyHit(b *testing.B) {
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Changed from old to new"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !filter.IsAgentMemo(comment) {
			b.Fatal("expected memo")
		}
	}
}

func BenchmarkAgentMemoLateHit(b *testing.B) {
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// Switched from old to new"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !filter.IsAgentMemo(comment) {
			b.Fatal("expected memo")
		}
	}
}

func BenchmarkAgentMemoFullMiss(b *testing.B) {
	filter := NewAgentMemoFilter()
	comment := models.CommentInfo{Text: "// calculate checksum before write"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if filter.IsAgentMemo(comment) {
			b.Fatal("expected non-memo")
		}
	}
}
