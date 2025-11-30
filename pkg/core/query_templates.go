package core

// QueryTemplates maps language names to their tree-sitter comment query patterns.
var QueryTemplates = map[string]string{
	"python":     "(comment) @comment",
	"javascript": "(comment) @comment",
	"typescript": "(comment) @comment",
	"tsx":        "(comment) @comment",
	"golang":     "(comment) @comment",
	"rust": `
		(line_comment) @comment
		(block_comment) @comment
	`,
	"swift": "(comment) @comment",
	"kotlin": `
		(line_comment) @comment
		(multiline_comment) @comment
	`,
	"java": `
		(line_comment) @comment
		(block_comment) @comment
	`,
	"elixir":   "(comment) @comment",
	"c":        "(comment) @comment",
	"cpp":      "(comment) @comment",
	"csharp":   "(comment) @comment",
	"ruby":     "(comment) @comment",
	"php":      "(comment) @comment",
	"bash":     "(comment) @comment",
	"lua":      "(comment) @comment",
	"ocaml":    "(comment) @comment",
	"sql":      "(comment) @comment",
	"html":     "(comment) @comment",
	"css":      "(comment) @comment",
	"yaml":     "(comment) @comment",
	"toml":     "(comment) @comment",
	"hcl":      "(comment) @comment",
	"svelte":   "(comment) @comment",
	"elm":      "(comment) @comment",
	"groovy":   "(comment) @comment",
	"cue":      "(comment) @comment",
	"scala":    "(comment) @comment",
	"protobuf": "(comment) @comment",
}

// DocstringQueries maps language names to their tree-sitter docstring query patterns.
var DocstringQueries = map[string]string{
	"python": `
		(module . (expression_statement (string) @docstring))
		(class_definition body: (block . (expression_statement (string) @docstring)))
		(function_definition body: (block . (expression_statement (string) @docstring)))
	`,
	"javascript": `
		(comment) @jsdoc
		(#match? @jsdoc "^/\\*\\*")
	`,
	"typescript": `
		(comment) @jsdoc
		(#match? @jsdoc "^/\\*\\*")
	`,
	"tsx": `
		(comment) @jsdoc
		(#match? @jsdoc "^/\\*\\*")
	`,
	"java": `
		(comment) @javadoc
		(#match? @javadoc "^/\\*\\*")
	`,
}
