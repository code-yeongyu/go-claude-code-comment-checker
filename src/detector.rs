use std::cell::RefCell;
use std::collections::HashMap;
use std::sync::{Arc, LazyLock, Mutex};

use tree_sitter::{Parser, Query, QueryCursor, StreamingIterator};

use crate::language::{LanguageSpec, language_for_path};
use crate::models::{CommentInfo, CommentType};

struct Runtime {
    language: tree_sitter::Language,
    comment_query: Query,
}

static RUNTIMES: LazyLock<Mutex<HashMap<&'static str, Arc<Runtime>>>> =
    LazyLock::new(|| Mutex::new(HashMap::new()));

thread_local! {
    static PARSERS: RefCell<HashMap<&'static str, Parser>> = RefCell::new(HashMap::new());
}

#[derive(Debug, Default)]
pub struct CommentDetector;

impl CommentDetector {
    #[must_use]
    pub const fn new() -> Self {
        Self
    }

    #[must_use]
    pub fn detect(
        &self,
        content: &str,
        file_path: &str,
        include_docstrings: bool,
    ) -> Vec<CommentInfo> {
        let Some(spec) = language_for_path(file_path) else {
            return Vec::new();
        };
        let Some(runtime) = runtime_for(spec) else {
            return Vec::new();
        };
        let Some(tree) = parse_with_runtime(&runtime, spec.name, content.as_bytes()) else {
            return Vec::new();
        };

        let mut comments = captures_to_comments(
            &runtime.comment_query,
            &tree,
            content.as_bytes(),
            file_path,
            false,
        );
        if include_docstrings {
            comments.extend(detect_docstrings(spec.name, content, file_path, &comments));
        }
        comments
    }
}

fn runtime_for(spec: LanguageSpec) -> Option<Arc<Runtime>> {
    let Ok(mut runtimes) = RUNTIMES.lock() else {
        return None;
    };
    if let Some(runtime) = runtimes.get(spec.name) {
        return Some(Arc::clone(runtime));
    }

    let language = tree_sitter_language_pack::get_language(spec.name).ok()?;
    let comment_query = Query::new(&language, spec.comment_query).ok()?;
    let runtime = Arc::new(Runtime {
        language,
        comment_query,
    });
    runtimes.insert(spec.name, Arc::clone(&runtime));
    Some(runtime)
}

fn parse_with_runtime(
    runtime: &Runtime,
    language_name: &'static str,
    source: &[u8],
) -> Option<tree_sitter::Tree> {
    PARSERS.with(|parsers| {
        let mut parsers = parsers.borrow_mut();
        if let Some(parser) = parsers.get_mut(language_name) {
            parser.parse(source, None)
        } else {
            let mut parser = Parser::new();
            if parser.set_language(&runtime.language).is_err() {
                return None;
            }
            let tree = parser.parse(source, None);
            parsers.insert(language_name, parser);
            tree
        }
    })
}

fn captures_to_comments(
    query: &Query,
    tree: &tree_sitter::Tree,
    source: &[u8],
    file_path: &str,
    force_docstring: bool,
) -> Vec<CommentInfo> {
    let mut cursor = QueryCursor::new();
    let mut captures = cursor.captures(query, tree.root_node(), source);
    let mut comments = Vec::new();
    captures.advance();
    while let Some((query_match, capture_index)) = captures.get() {
        let Some(capture) = query_match.captures.get(*capture_index) else {
            captures.advance();
            continue;
        };
        let node = capture.node;
        let Ok(text) = node.utf8_text(source) else {
            captures.advance();
            continue;
        };
        let comment_type = if force_docstring {
            CommentType::Docstring
        } else {
            determine_comment_type(text, node.kind())
        };
        let is_docstring = comment_type == CommentType::Docstring;
        comments.push(CommentInfo::new(
            text,
            node.start_position().row + 1,
            file_path,
            comment_type,
            is_docstring,
        ));
        captures.advance();
    }
    comments
}

fn determine_comment_type(text: &str, node_kind: &str) -> CommentType {
    let stripped = text.trim();
    if node_kind == "line_comment" {
        return CommentType::Line;
    }
    if node_kind == "block_comment" || node_kind == "multiline_comment" {
        return CommentType::Block;
    }
    if stripped.starts_with(r#"""""#) || stripped.starts_with("'''") {
        return CommentType::Docstring;
    }
    if stripped.starts_with("//") || stripped.starts_with('#') {
        return CommentType::Line;
    }
    if stripped.starts_with("/*") || stripped.starts_with("<!--") || stripped.starts_with("--") {
        return CommentType::Block;
    }
    CommentType::Line
}

fn detect_docstrings(
    language_name: &str,
    content: &str,
    file_path: &str,
    comments: &[CommentInfo],
) -> Vec<CommentInfo> {
    match language_name {
        "python" => detect_python_docstrings(content, file_path),
        "javascript" | "typescript" | "tsx" | "java" => comments
            .iter()
            .filter(|comment| comment.text.trim_start().starts_with("/**"))
            .map(|comment| {
                CommentInfo::new(
                    comment.text.clone(),
                    comment.line_number,
                    file_path,
                    CommentType::Docstring,
                    true,
                )
            })
            .collect(),
        _ => Vec::new(),
    }
}

fn detect_python_docstrings(content: &str, file_path: &str) -> Vec<CommentInfo> {
    let mut docstrings = Vec::new();
    for (index, line) in content.lines().enumerate() {
        let trimmed = line.trim_start();
        if trimmed.starts_with(r#"""""#) || trimmed.starts_with("'''") {
            docstrings.push(CommentInfo::new(
                trimmed,
                index + 1,
                file_path,
                CommentType::Docstring,
                true,
            ));
        }
    }
    docstrings
}
