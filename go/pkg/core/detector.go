package core

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/code-yeongyu/go-claude-code-comment-checker/pkg/models"
)

type CommentDetector struct {
	registry         *LanguageRegistry
	queryMu          sync.Mutex
	commentQueries   map[string]*sitter.Query
	docstringQueries map[string]*sitter.Query
}

func NewCommentDetector() *CommentDetector {
	return &CommentDetector{
		registry:         NewLanguageRegistry(),
		commentQueries:   make(map[string]*sitter.Query),
		docstringQueries: make(map[string]*sitter.Query),
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

	query := d.commentQueryFor(langName, lang)
	if query == nil {
		return nil
	}

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
	query := d.docstringQueryFor(langName, lang)
	if query == nil {
		return nil
	}

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

func (d *CommentDetector) commentQueryFor(langName string, lang *sitter.Language) *sitter.Query {
	queryPattern := QueryTemplates[langName]
	if queryPattern == "" {
		queryPattern = "(comment) @comment"
	}

	return d.cachedQuery(d.commentQueries, langName, queryPattern, lang)
}

func (d *CommentDetector) docstringQueryFor(langName string, lang *sitter.Language) *sitter.Query {
	docQuery, ok := DocstringQueries[langName]
	if !ok {
		return nil
	}

	return d.cachedQuery(d.docstringQueries, langName, docQuery, lang)
}

func (d *CommentDetector) cachedQuery(
	cache map[string]*sitter.Query,
	langName string,
	queryPattern string,
	lang *sitter.Language,
) *sitter.Query {
	d.queryMu.Lock()
	defer d.queryMu.Unlock()

	query := cache[langName]
	if query != nil {
		return query
	}

	compiledQuery, err := sitter.NewQuery([]byte(queryPattern), lang)
	if err != nil {
		return nil
	}

	cache[langName] = compiledQuery
	return compiledQuery
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
