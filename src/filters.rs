use crate::models::CommentInfo;

const BDD_KEYWORDS: &[&str] = &[
    "given",
    "when",
    "then",
    "arrange",
    "act",
    "assert",
    "when & then",
    "when&then",
];

const DIRECTIVE_PREFIXES: &[&str] = &[
    "type:",
    "noqa",
    "pyright:",
    "ruff:",
    "mypy:",
    "pylint:",
    "flake8:",
    "pyre:",
    "pytype:",
    "eslint-disable",
    "eslint-ignore",
    "prettier-ignore",
    "ts-ignore",
    "ts-expect-error",
    "clippy:",
    "allow",
    "deny",
    "warn",
    "forbid",
];

#[derive(Clone, Copy, Debug, Default)]
pub struct BddFilter;

impl BddFilter {
    #[must_use]
    pub fn should_skip(&self, comment: &CommentInfo) -> bool {
        let normalized = strip_one_prefix(comment.text.trim(), &["#", "//", "--"]);
        BDD_KEYWORDS
            .iter()
            .any(|keyword| normalized.eq_ignore_ascii_case(keyword))
    }
}

#[derive(Clone, Copy, Debug, Default)]
pub struct DirectiveFilter;

impl DirectiveFilter {
    #[must_use]
    pub fn should_skip(&self, comment: &CommentInfo) -> bool {
        let mut normalized = strip_one_prefix(comment.text.trim(), &["#", "//", "/*", "--"]);
        if let Some(rest) = normalized.strip_prefix('@') {
            normalized = rest.trim();
        }
        DIRECTIVE_PREFIXES
            .iter()
            .any(|directive| starts_with_ignore_ascii_case(normalized, directive))
    }
}

#[derive(Clone, Copy, Debug, Default)]
pub struct ShebangFilter;

impl ShebangFilter {
    #[must_use]
    pub fn should_skip(&self, comment: &CommentInfo) -> bool {
        comment.text.trim().starts_with("#!")
    }
}

#[derive(Clone, Copy, Debug, Default)]
pub struct AgentMemoFilter;

impl AgentMemoFilter {
    #[must_use]
    pub fn is_agent_memo(&self, comment: &CommentInfo) -> bool {
        let text = strip_agent_prefixes(comment.text.trim());
        is_english_agent_memo(text)
    }
}

fn strip_one_prefix<'a>(text: &'a str, prefixes: &[&str]) -> &'a str {
    let mut normalized = text.trim();
    for prefix in prefixes {
        if let Some(rest) = normalized.strip_prefix(prefix) {
            normalized = rest.trim();
            break;
        }
    }
    normalized
}

fn strip_agent_prefixes(mut text: &str) -> &str {
    for prefix in ["#", "//", "/*", "--", "*"] {
        if let Some(rest) = text.strip_prefix(prefix) {
            text = rest.trim();
        }
    }
    text
}

fn is_english_agent_memo(text: &str) -> bool {
    let starts_like_memo = match first_ascii_lowercase_byte(text) {
        Some(b'a') => starts_with_any(text, &["added", "add", "after this"]),
        Some(b'b') => starts_with_ignore_ascii_case(text, "before this"),
        Some(b'c') => {
            starts_with_word_pair(text, "changed", &["from", "to"])
                || starts_with_word_pair(text, "change", &["from", "to"])
                || starts_with_word_pair(text, "converted", &["from", "to"])
                || starts_with_word_pair(text, "convert", &["from", "to"])
        }
        Some(b'd') => starts_with_any(text, &["deleted", "delete"]),
        Some(b'h') => starts_with_ignore_ascii_case(text, "here we"),
        Some(b'i') => {
            starts_with_any(text, &["implemented", "implement"])
                || starts_with_word_pair(text, "implementation", &["of", "note"])
        }
        Some(b'm') => {
            starts_with_word_pair(text, "modified", &["from", "to"])
                || starts_with_word_pair(text, "modify", &["from", "to"])
                || starts_with_word_pair(text, "moved", &["from", "to"])
                || starts_with_word_pair(text, "move", &["from", "to"])
                || starts_with_word_pair(text, "migrated", &["from", "to"])
                || starts_with_word_pair(text, "migrate", &["from", "to"])
        }
        Some(b'n') => {
            starts_with_word_pair(text, "now", &["we", "this", "it"])
                || starts_with_ignore_ascii_case(text, "note:")
        }
        Some(b'p') => starts_with_ignore_ascii_case(text, "previously"),
        Some(b'r') => {
            starts_with_any(
                text,
                &["refactor", "replaced", "replace", "removed", "remove"],
            ) || starts_with_word_pair(text, "renamed", &["from", "to"])
                || starts_with_word_pair(text, "rename", &["from", "to"])
        }
        Some(b's') => {
            starts_with_word_pair(text, "switched", &["from", "to"])
                || starts_with_word_pair(text, "switch", &["from", "to"])
        }
        Some(b't') => starts_with_word_pair(
            text,
            "this",
            &[
                "implements",
                "implement",
                "adds",
                "add",
                "removes",
                "remove",
                "changes",
                "change",
                "fixes",
                "fix",
            ],
        ),
        Some(b'u') => {
            starts_with_word_pair(text, "updated", &["from", "to"])
                || starts_with_word_pair(text, "update", &["from", "to"])
        }
        Some(b'w') => starts_with_word_pair(text, "was", &["changed"]),
        _ => false,
    };
    starts_like_memo || has_ascii_arrow(text)
}

fn first_ascii_lowercase_byte(text: &str) -> Option<u8> {
    text.as_bytes().first().map(u8::to_ascii_lowercase)
}

fn starts_with_any(text: &str, prefixes: &[&str]) -> bool {
    prefixes
        .iter()
        .any(|prefix| starts_with_ignore_ascii_case(text, prefix))
}

fn starts_with_word_pair(text: &str, first: &str, seconds: &[&str]) -> bool {
    let Some(rest) = strip_prefix_ignore_ascii_case(text, first) else {
        return false;
    };
    let rest = rest.trim_start();
    seconds
        .iter()
        .any(|second| starts_with_ignore_ascii_case(rest, second))
}

fn has_ascii_arrow(text: &str) -> bool {
    let Some((left, right)) = text.split_once("->") else {
        return false;
    };
    let left = left.trim();
    let right = right.trim();
    !left.is_empty()
        && !right.is_empty()
        && left.bytes().all(|byte| byte.is_ascii_alphabetic())
        && right
            .bytes()
            .next()
            .is_some_and(|byte| byte.is_ascii_alphabetic())
}

fn strip_prefix_ignore_ascii_case<'a>(text: &'a str, prefix: &str) -> Option<&'a str> {
    starts_with_ignore_ascii_case(text, prefix).then(|| &text[prefix.len()..])
}

fn starts_with_ignore_ascii_case(text: &str, prefix: &str) -> bool {
    let text = text.as_bytes();
    let prefix = prefix.as_bytes();
    text.len() >= prefix.len() && text[..prefix.len()].eq_ignore_ascii_case(prefix)
}

#[must_use]
pub fn apply_filters(comments: &[CommentInfo]) -> Vec<CommentInfo> {
    let bdd = BddFilter;
    let directive = DirectiveFilter;
    let shebang = ShebangFilter;
    comments
        .iter()
        .filter(|comment| {
            !bdd.should_skip(comment)
                && !directive.should_skip(comment)
                && !shebang.should_skip(comment)
        })
        .cloned()
        .collect()
}
