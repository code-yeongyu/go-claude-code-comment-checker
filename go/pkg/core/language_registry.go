package core

import (
	"strings"

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

var ExtensionToLanguage = map[string]string{
	"py": "python",

	"js": "javascript", "jsx": "javascript",
	"ts": "typescript", "tsx": "tsx",

	"go": "golang",

	"java": "java", "kt": "kotlin", "scala": "scala",

	"c": "c", "h": "c",
	"cpp": "cpp", "cc": "cpp", "cxx": "cpp", "hpp": "cpp",

	"rs": "rust",

	"rb": "ruby",

	"sh": "bash", "bash": "bash",

	"cs": "csharp",

	"swift": "swift",

	"ex": "elixir", "exs": "elixir",

	"lua": "lua",

	"php": "php",

	"ml": "ocaml", "mli": "ocaml",

	"sql": "sql",

	"html": "html", "htm": "html",
	"css": "css",

	"yaml": "yaml", "yml": "yaml",
	"toml": "toml",
	"hcl":  "hcl", "tf": "hcl",

	"dockerfile": "dockerfile",
	"proto":      "protobuf",
	"svelte":     "svelte",
	"elm":        "elm",
	"groovy":     "groovy",
	"cue":        "cue",
}

type LanguageRegistry struct{}

func NewLanguageRegistry() *LanguageRegistry {
	return &LanguageRegistry{}
}

func (r *LanguageRegistry) GetLanguageName(extension string) string {
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))
	return ExtensionToLanguage[ext]
}

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

func (r *LanguageRegistry) IsSupported(extension string) bool {
	return r.GetLanguageName(extension) != ""
}
