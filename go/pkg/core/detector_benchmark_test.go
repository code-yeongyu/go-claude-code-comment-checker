package core

import (
	"strings"
	"testing"
)

var benchmarkPythonSource = strings.Join([]string{
	`"""Module docstring."""`,
	"",
	"# regular comment",
	"def add(a, b):",
	`    """Function docstring."""`,
	"    return a + b",
	"",
}, "\n")

func BenchmarkCommentDetectorDetectPythonNoDocstrings(b *testing.B) {
	detector := NewCommentDetector()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comments := detector.Detect(benchmarkPythonSource, "sample.py", false)
		if len(comments) != 1 {
			b.Fatalf("expected 1 comment, got %d", len(comments))
		}
	}
}

func BenchmarkCommentDetectorDetectPythonWithDocstrings(b *testing.B) {
	detector := NewCommentDetector()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comments := detector.Detect(benchmarkPythonSource, "sample.py", true)
		if len(comments) != 3 {
			b.Fatalf("expected 3 comments, got %d", len(comments))
		}
	}
}

func BenchmarkCommentDetectorUnsupportedExtension(b *testing.B) {
	detector := NewCommentDetector()
	source := benchmarkPythonSource + "\n"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comments := detector.Detect(source, "sample.unknown", true)
		if comments != nil {
			b.Fatalf("expected nil comments, got %d", len(comments))
		}
	}
}
