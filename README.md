# comment-checker

Multi-language comment detection for Claude Code hooks. Go implementation with tree-sitter.

## Installation

```bash
go build ./cmd/comment-checker
```

Or install globally:

```bash
go install github.com/yeongyu/comment-checker/cmd/comment-checker@latest
```

## Usage

Pipe JSON input from Claude Code hooks:

```bash
echo '{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"print(1)"}}' | ./comment-checker
```

### Input Format

```json
{
  "tool_name": "Write",
  "tool_input": {
    "file_path": "path/to/file.py",
    "content": "# comment\nprint(1)"
  }
}
```

Supported `tool_name` values:
- `Write`: Uses `content` field
- `Edit`: Uses `new_string` field
- `MultiEdit`: Combines all `edits[].new_string` fields

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Pass (no comments) or Skip (non-code file, invalid input) |
| 2 | Block (problematic comments found) |

## Supported Languages

30+ languages via go-tree-sitter:

- **Popular**: Python, Go, TypeScript, JavaScript, Java, Rust, C, C++, C#
- **Scripting**: Ruby, PHP, Bash, Lua, Perl
- **Functional**: Elixir, Scala, OCaml, Elm
- **Config**: YAML, TOML, HCL, Dockerfile
- **Other**: Swift, Kotlin, Groovy, SQL, HTML, CSS, Markdown

## Filters

Comments are automatically filtered if they match:

| Filter | Example |
|--------|---------|
| BDD | `# given`, `# when`, `# then` |
| Docstring | `"""docstring"""`, `/** JSDoc */` |
| Directive | `# noqa`, `// @ts-ignore`, `// eslint-disable` |
| Shebang | `#!/usr/bin/env python` |

## Library Usage

```go
import (
    "github.com/yeongyu/comment-checker/pkg/core"
    "github.com/yeongyu/comment-checker/pkg/filters"
    "github.com/yeongyu/comment-checker/pkg/output"
)

// Detect comments
detector := core.NewCommentDetector()
comments := detector.Detect(content, "file.py", false)

// Apply filters
bddFilter := filters.NewBDDFilter()
docstringFilter := filters.NewDocstringFilter()
directiveFilter := filters.NewDirectiveFilter()
shebangFilter := filters.NewShebangFilter()

var filtered []models.CommentInfo
for _, c := range comments {
    if bddFilter.ShouldSkip(c) || docstringFilter.ShouldSkip(c) ||
       directiveFilter.ShouldSkip(c) || shebangFilter.ShouldSkip(c) {
        continue
    }
    filtered = append(filtered, c)
}

// Format output
message := output.FormatHookMessage(filtered)
```

## Development

```bash
# Build
go build ./...

# Test
go test ./... -v

# Lint
go vet ./...
```

## License

MIT
