# comment-checker

[한국어](README.ko.md)

> 100% vibe coded. zero comments in my code, zero comments in yours.
>
> **this entire repo - code, docs, readme, everything - was written by LLMs.**
> shoutout to claude opus 4.5.

A PostToolUse hook for Claude Code / OpenCode that yells at you when you write unnecessary comments.

Built with Go + tree-sitter. Fast. Opinionated. No mercy.

## why

comments are code smell. if your code needs comments to be understood, your code sucks.

this hook watches every `Write`, `Edit`, `MultiEdit` and screams when it detects comments.

exceptions exist. BDD comments (`# given`, `# when`, `# then`), linter directives (`# noqa`, `// @ts-ignore`), shebangs - these are fine. everything else? delete it. yes, even docstrings.

## install

### homebrew (macos/linux)

```bash
brew tap code-yeongyu/tap
brew install comment-checker
```

### go install

```bash
go install github.com/code-yeongyu/go-claude-code-comment-checker/cmd/comment-checker@latest
```

### manual

grab binary from [releases](https://github.com/code-yeongyu/go-claude-code-comment-checker/releases).

## setup

add to `~/.claude/settings.json` (or `.claude/settings.json` in your project):

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "comment-checker"
          }
        ]
      }
    ]
  }
}
```

done. now claude will think twice before leaving `// TODO: fix later` in your code.

## what it catches

```go
// bad: unnecessary comment
x := 1 + 1  // adds one to one

// bad: todo that will never be done
// TODO: refactor this later

// bad: commented out code
// fmt.Println("debug")
```

## what it allows

```python
# given - BDD comments are fine
def test_something():
    # when
    result = do_thing()
    # then
    assert result == expected
```

```python
# noqa: E501 - linter directives are fine
```

```python
#!/usr/bin/env python - shebangs are fine
```

## 30+ languages

python, go, typescript, javascript, rust, c, c++, java, ruby, php, swift, kotlin, scala, elixir, and more.

if tree-sitter supports it, we support it.

## how it works

1. hook receives JSON from Claude Code
2. extracts content from `Write`/`Edit`/`MultiEdit` tool input
3. detects language from file extension
4. parses AST with tree-sitter
5. finds comment nodes
6. filters out allowed patterns (BDD, directives, shebangs)
7. if anything remains → exit 2 with warning message

## exit codes

| code | meaning |
|------|---------|
| 0 | pass - no comments found or skipped |
| 2 | warning - problematic comments detected |

## philosophy

> "Code is like humor. When you have to explain it, it's bad." - Cory House

write self-documenting code. use meaningful names. extract functions. stop explaining what the code does and make the code explain itself.

## license

MIT. do whatever.
