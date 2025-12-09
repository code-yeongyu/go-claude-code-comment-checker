# comment-checker: 주석 감지 훅

> 이 저장소의 코드, 문서, README 모두 LLM으로 작성되었습니다. Claude Opus 4.5에게 감사를 전합니다.

Claude Code / OpenCode가 코드를 작성할 때마다 불필요한 주석을 남기곤 합니다.
"여기서 변수를 선언합니다", "1에 1을 더합니다" 같은 주석들입니다.

comment-checker는 이런 주석을 감지하고 경고하는 PostToolUse 훅입니다.
Go와 tree-sitter 기반으로 30개 이상의 언어를 지원합니다.

## 왜 필요한가요

주석이 필요하다는 건 코드가 충분히 명확하지 않다는 신호입니다.
좋은 코드는 그 자체로 의도를 설명합니다.

물론 모든 주석이 나쁜 건 아닙니다.
BDD 테스트 주석, 린터 지시문, shebang은 허용합니다.
하지만 "i를 1 증가시킵니다" 같은 주석이나 docstring은 삭제해야 합니다.

## 설치

### Homebrew (macOS/Linux)

```sh
brew tap code-yeongyu/tap
brew install comment-checker
```

### Go Install

```sh
go install github.com/code-yeongyu/go-claude-code-comment-checker/cmd/comment-checker@latest
```

### 바이너리 직접 다운로드

[Releases](https://github.com/code-yeongyu/go-claude-code-comment-checker/releases) 페이지에서 플랫폼에 맞는 바이너리를 받아 PATH에 추가하세요.

---

## 설정

`~/.claude/settings.json` (또는 프로젝트의 `.claude/settings.json`)에 다음을 추가합니다:

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

이제 Claude가 `Write`, `Edit`, `MultiEdit` 도구를 사용할 때마다 주석을 검사합니다.

---

## 허용되는 주석

### BDD 테스트 주석

테스트 코드의 given-when-then 구조는 허용합니다.

```python
def test_user_login():
    # given
    user = create_user()
    
    # when
    result = login(user)
    
    # then
    assert result.success
```

### 린터/타입체커 지시문

도구 제어를 위한 주석은 허용합니다.

```python
# noqa: E501
# type: ignore
# pylint: disable=line-too-long
```

```typescript
// @ts-ignore
// eslint-disable-next-line
```

### Shebang

스크립트 실행을 위한 shebang은 허용합니다.

```python
#!/usr/bin/env python3
```

---

## 경고 대상

### 설명 주석

코드를 보면 알 수 있는 내용입니다.

```go
x := 1 + 1  // 1에 1을 더함
```

### TODO 주석

TODO는 보통 방치됩니다. 지금 하거나, 이슈로 등록하세요.

```python
# TODO: 나중에 리팩토링
```

### 주석 처리된 코드

버전 관리 시스템이 있으니 삭제해도 됩니다.

```javascript
// console.log(debugInfo);
```

---

## 지원 언어

tree-sitter 덕분에 대부분의 주요 언어를 지원합니다.

- **웹**: TypeScript, JavaScript, HTML, CSS, Vue, Svelte
- **백엔드**: Go, Python, Rust, Java, C#, Ruby, PHP
- **시스템**: C, C++
- **설정**: YAML, TOML, JSON
- **기타**: SQL, Shell, Kotlin, Swift, Scala, Elixir 등 30개 이상

---

## 종료 코드

| 코드 | 의미 |
|------|------|
| 0 | 주석 없음 또는 허용된 주석만 있음 |
| 2 | 불필요한 주석 감지됨 |

---

## 라이선스

MIT
