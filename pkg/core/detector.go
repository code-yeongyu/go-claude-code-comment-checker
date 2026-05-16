package core

import (
	"context"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

type CommentDetector struct {
	registry *LanguageRegistry
}

func NewCommentDetector() *CommentDetector {
	return &CommentDetector{
		registry: NewLanguageRegistry(),
	}
}

func (d *CommentDetector) Detect(content, filePath string, includeDocstrings bool) []models.CommentInfo {
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	if ext == "" {
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

	if includeDocstrings {
		docstrings := d.detectDocstrings(sourceCode, filePath, lang, langName, tree)
		comments = append(comments, docstrings...)
	}

	return comments
}

func (d *CommentDetector) detectDocstrings(sourceCode []byte, filePath string, lang *sitter.Language, langName string, tree *sitter.Tree) []models.CommentInfo {
	docQuery, ok := DocstringQueries[langName]
	if !ok {
		return nil
	}

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

func (d *CommentDetector) determineCommentType(text, nodeType string) models.CommentType {
	stripped := strings.TrimSpace(text)

	if nodeType == "line_comment" {
		return models.CommentTypeLine
	}
	if nodeType == "block_comment" {
		return models.CommentTypeBlock
	}

	if strings.HasPrefix(stripped, `"""`) || strings.HasPrefix(stripped, "'''") {
		return models.CommentTypeDocstring
	}

	if strings.HasPrefix(stripped, "//") || strings.HasPrefix(stripped, "#") {
		return models.CommentTypeLine
	}

	if strings.HasPrefix(stripped, "/*") || strings.HasPrefix(stripped, "<!--") || strings.HasPrefix(stripped, "--") {
		return models.CommentTypeBlock
	}

	return models.CommentTypeLine
}
