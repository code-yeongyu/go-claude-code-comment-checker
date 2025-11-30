package core

import (
	"context"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

// CommentDetector detects comments in source code using tree-sitter.
type CommentDetector struct {
	registry *LanguageRegistry
}

// NewCommentDetector creates a new CommentDetector instance.
func NewCommentDetector() *CommentDetector {
	return &CommentDetector{
		registry: NewLanguageRegistry(),
	}
}

// Detect extracts comments from the given source code.
func (d *CommentDetector) Detect(content, filePath string, includeDocstrings bool) []models.CommentInfo {
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	if ext == "" {
		// Handle files like "Dockerfile"
		ext = strings.ToLower(filepath.Base(filePath))
	}

	langName := d.registry.GetLanguageName(ext)
	if langName == "" {
		return nil
	}

	lang := GetLanguage(langName)
	if lang == nil {
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	sourceCode := []byte(content)
	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		return nil
	}
	defer tree.Close()

	queryPattern := QueryTemplates[langName]
	if queryPattern == "" {
		queryPattern = "(comment) @comment"
	}

	query, err := sitter.NewQuery([]byte(queryPattern), lang)
	if err != nil {
		return nil
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(query, tree.RootNode())

	var comments []models.CommentInfo
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			text := node.Content(sourceCode)
			lineNumber := int(node.StartPoint().Row) + 1

			commentType := d.determineCommentType(text, node.Type())
			isDocstring := commentType == models.CommentTypeDocstring

			if isDocstring && !includeDocstrings {
				continue
			}

			comments = append(comments, models.CommentInfo{
				Text:        text,
				LineNumber:  lineNumber,
				FilePath:    filePath,
				CommentType: commentType,
				IsDocstring: isDocstring,
			})
		}
	}

	// Detect docstrings if requested
	if includeDocstrings {
		docstrings := d.detectDocstrings(sourceCode, filePath, lang, langName)
		comments = append(comments, docstrings...)
	}

	return comments
}

// detectDocstrings extracts docstrings using language-specific queries.
func (d *CommentDetector) detectDocstrings(sourceCode []byte, filePath string, lang *sitter.Language, langName string) []models.CommentInfo {
	docQuery, ok := DocstringQueries[langName]
	if !ok {
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		return nil
	}
	defer tree.Close()

	query, err := sitter.NewQuery([]byte(docQuery), lang)
	if err != nil {
		return nil
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(query, tree.RootNode())

	var docstrings []models.CommentInfo
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			text := node.Content(sourceCode)
			lineNumber := int(node.StartPoint().Row) + 1

			docstrings = append(docstrings, models.CommentInfo{
				Text:        text,
				LineNumber:  lineNumber,
				FilePath:    filePath,
				CommentType: models.CommentTypeDocstring,
				IsDocstring: true,
			})
		}
	}

	return docstrings
}

// determineCommentType determines the type of comment based on its text and node type.
func (d *CommentDetector) determineCommentType(text, nodeType string) models.CommentType {
	stripped := strings.TrimSpace(text)

	// Check node type first (for Rust)
	if nodeType == "line_comment" {
		return models.CommentTypeLine
	}
	if nodeType == "block_comment" {
		return models.CommentTypeBlock
	}

	// Check for docstrings
	if strings.HasPrefix(stripped, `"""`) || strings.HasPrefix(stripped, "'''") {
		return models.CommentTypeDocstring
	}

	// Check for line comments
	if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") {
		return models.CommentTypeLine
	}

	// Check for block comments
	if strings.HasPrefix(stripped, "/*") || strings.HasPrefix(stripped, "<!--") || strings.HasPrefix(stripped, "--") {
		return models.CommentTypeBlock
	}

	return models.CommentTypeLine
}
