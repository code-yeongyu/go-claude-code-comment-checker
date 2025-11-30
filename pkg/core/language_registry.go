package core

import (
	"strings"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/smacker/go-tree-sitter/css"
	"github.com/smacker/go-tree-sitter/cue"
	"github.com/smacker/go-tree-sitter/dockerfile"
	"github.com/smacker/go-tree-sitter/elixir"
	"github.com/smacker/go-tree-sitter/elm"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/groovy"
	"github.com/smacker/go-tree-sitter/hcl"
	"github.com/smacker/go-tree-sitter/html"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/kotlin"
	"github.com/smacker/go-tree-sitter/lua"
	"github.com/smacker/go-tree-sitter/ocaml"
	"github.com/smacker/go-tree-sitter/php"
	"github.com/smacker/go-tree-sitter/protobuf"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/smacker/go-tree-sitter/scala"
	"github.com/smacker/go-tree-sitter/sql"
	"github.com/smacker/go-tree-sitter/svelte"
	"github.com/smacker/go-tree-sitter/swift"
	"github.com/smacker/go-tree-sitter/toml"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/smacker/go-tree-sitter/yaml"
)

// ExtensionToLanguage maps file extensions to tree-sitter language names.
var ExtensionToLanguage = map[string]string{
	// Python
	"py": "python",
	// JavaScript/TypeScript
	"js": "javascript", "jsx": "javascript",
	"ts": "typescript", "tsx": "tsx",
	// Go
	"go": "golang",
	// Java/Kotlin/Scala
	"java": "java", "kt": "kotlin", "scala": "scala",
	// C/C++
	"c": "c", "h": "c",
	"cpp": "cpp", "cc": "cpp", "cxx": "cpp", "hpp": "cpp",
	// Rust
	"rs": "rust",
	// Ruby
	"rb": "ruby",
	// Shell
	"sh": "bash", "bash": "bash",
	// C#
	"cs": "csharp",
	// Swift
	"swift": "swift",
	// Elixir
	"ex": "elixir", "exs": "elixir",
	// Lua
	"lua": "lua",
	// PHP
	"php": "php",
	// OCaml
	"ml": "ocaml", "mli": "ocaml",
	// SQL
	"sql": "sql",
	// Web
	"html": "html", "htm": "html",
	"css": "css",
	// Config
	"yaml": "yaml", "yml": "yaml",
	"toml": "toml",
	"hcl":  "hcl", "tf": "hcl",
	// Others
	"dockerfile": "dockerfile",
	"proto":      "protobuf",
	"svelte":     "svelte",
	"elm":        "elm",
	"groovy":     "groovy",
	"cue":        "cue",
}

// LanguageRegistry provides thread-safe access to tree-sitter parsers.
type LanguageRegistry struct {
	mu      sync.RWMutex
	parsers map[string]*sitter.Parser
}

// NewLanguageRegistry creates a new LanguageRegistry instance.
func NewLanguageRegistry() *LanguageRegistry {
	return &LanguageRegistry{
		parsers: make(map[string]*sitter.Parser),
	}
}

// GetLanguageName returns the tree-sitter language name for a file extension.
func (r *LanguageRegistry) GetLanguageName(extension string) string {
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))
	return ExtensionToLanguage[ext]
}

// GetLanguage returns the tree-sitter Language for the given language name.
func GetLanguage(name string) *sitter.Language {
	switch name {
	case "python":
		return python.GetLanguage()
	case "javascript":
		return javascript.GetLanguage()
	case "typescript":
		return typescript.GetLanguage()
	case "tsx":
		return tsx.GetLanguage()
	case "golang":
		return golang.GetLanguage()
	case "java":
		return java.GetLanguage()
	case "c":
		return c.GetLanguage()
	case "cpp":
		return cpp.GetLanguage()
	case "rust":
		return rust.GetLanguage()
	case "ruby":
		return ruby.GetLanguage()
	case "bash":
		return bash.GetLanguage()
	case "csharp":
		return csharp.GetLanguage()
	case "kotlin":
		return kotlin.GetLanguage()
	case "swift":
		return swift.GetLanguage()
	case "elixir":
		return elixir.GetLanguage()
	case "lua":
		return lua.GetLanguage()
	case "php":
		return php.GetLanguage()
	case "scala":
		return scala.GetLanguage()
	case "ocaml":
		return ocaml.GetLanguage()
	case "sql":
		return sql.GetLanguage()
	case "html":
		return html.GetLanguage()
	case "css":
		return css.GetLanguage()
	case "yaml":
		return yaml.GetLanguage()
	case "toml":
		return toml.GetLanguage()
	case "dockerfile":
		return dockerfile.GetLanguage()
	case "protobuf":
		return protobuf.GetLanguage()
	case "hcl":
		return hcl.GetLanguage()
	case "svelte":
		return svelte.GetLanguage()
	case "elm":
		return elm.GetLanguage()
	case "groovy":
		return groovy.GetLanguage()
	case "cue":
		return cue.GetLanguage()
	default:
		return nil
	}
}

// GetParser returns a parser for the given extension.
// Parsers are cached for reuse.
func (r *LanguageRegistry) GetParser(extension string) *sitter.Parser {
	langName := r.GetLanguageName(extension)
	if langName == "" {
		return nil
	}

	r.mu.RLock()
	if parser, ok := r.parsers[langName]; ok {
		r.mu.RUnlock()
		return parser
	}
	r.mu.RUnlock()

	lang := GetLanguage(langName)
	if lang == nil {
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	r.mu.Lock()
	r.parsers[langName] = parser
	r.mu.Unlock()

	return parser
}

// IsSupported returns true if the extension is supported.
func (r *LanguageRegistry) IsSupported(extension string) bool {
	return r.GetLanguageName(extension) != ""
}
