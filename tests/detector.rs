use comment_checker::detector::CommentDetector;
use comment_checker::models::CommentType;

#[test]
fn detect_given_python_line_comment_returns_comment_info() {
    // given
    let detector = CommentDetector::new();
    let code = "# This is a comment\nprint('hello')";

    // when
    let comments = detector.detect(code, "test.py", false);

    // then
    assert_eq!(comments.len(), 1);
    assert_eq!(comments[0].text, "# This is a comment");
    assert_eq!(comments[0].line_number, 1);
    assert_eq!(comments[0].file_path, "test.py");
    assert_eq!(comments[0].comment_type, CommentType::Line);
    assert!(!comments[0].is_docstring);
}

#[test]
fn detect_given_typescript_block_comment_returns_block_comment() {
    // given
    let detector = CommentDetector::new();
    let code = "/* This is a block comment */\nconst x = 1;";

    // when
    let comments = detector.detect(code, "test.ts", false);

    // then
    assert_eq!(comments.len(), 1);
    assert_eq!(comments[0].text, "/* This is a block comment */");
    assert_eq!(comments[0].comment_type, CommentType::Block);
}

#[test]
fn detect_given_python_docstring_with_docstrings_returns_docstring() {
    // given
    let detector = CommentDetector::new();
    let code = "\"\"\"This is a module docstring.\"\"\"\ndef hello():\n    pass";

    // when
    let comments = detector.detect(code, "module.py", true);

    // then
    assert!(comments.iter().any(|comment| {
        comment.is_docstring
            && comment.comment_type == CommentType::Docstring
            && comment.text.contains("module docstring")
    }));
}

#[test]
fn detect_given_unsupported_extension_returns_empty_list() {
    // given
    let detector = CommentDetector::new();

    // when
    let comments = detector.detect("some random content", "test.xyz", false);

    // then
    assert!(comments.is_empty());
}

#[test]
fn detect_given_go_comment_returns_comment_info() {
    // given
    let detector = CommentDetector::new();
    let code = "// This is a Go comment\npackage main\n\nfunc main() {}";

    // when
    let comments = detector.detect(code, "main.go", false);

    // then
    assert_eq!(comments.len(), 1);
    assert_eq!(comments[0].text, "// This is a Go comment");
    assert_eq!(comments[0].line_number, 1);
    assert_eq!(comments[0].file_path, "main.go");
    assert_eq!(comments[0].comment_type, CommentType::Line);
    assert!(!comments[0].is_docstring);
}

#[test]
fn detect_given_multilanguage_smoke_cases_returns_comment() {
    // given
    let detector = CommentDetector::new();
    let cases = [
        ("main.go", "// Go comment\npackage main", "Go comment"),
        (
            "test.ts",
            "// TypeScript comment\nconst x: number = 1;",
            "TypeScript comment",
        ),
        (
            "test.js",
            "// JavaScript comment\nconst x = 1;",
            "JavaScript comment",
        ),
        (
            "Test.java",
            "// Java comment\npublic class Test {}",
            "Java comment",
        ),
        (
            "main.c",
            "// C comment\nint main() { return 0; }",
            "C comment",
        ),
        (
            "main.cpp",
            "// C++ comment\nint main() { return 0; }",
            "C++ comment",
        ),
        ("main.rs", "// Rust comment\nfn main() {}", "Rust comment"),
        ("test.rb", "# Ruby comment\nputs \"hello\"", "Ruby comment"),
        (
            "script.sh",
            "# Bash comment\necho \"hello\"",
            "Bash comment",
        ),
        (
            "Main.kt",
            "// Kotlin comment\nfun main() {}",
            "Kotlin comment",
        ),
        (
            "main.swift",
            "// Swift comment\nprint(\"hello\")",
            "Swift comment",
        ),
    ];

    // when
    let results: Vec<_> = cases
        .iter()
        .map(|(path, source, expected)| {
            let comments = detector.detect(source, path, false);
            (comments, expected)
        })
        .collect();

    // then
    for (comments, expected) in results {
        assert_eq!(comments.len(), 1, "expected one comment for {expected}");
        assert!(comments[0].text.contains(expected));
        assert_eq!(comments[0].comment_type, CommentType::Line);
    }
}

#[test]
fn detect_given_javascript_block_comment_returns_block_comment() {
    // given
    let detector = CommentDetector::new();
    let code = "/* block comment */\nconst x = 1;";

    // when
    let comments = detector.detect(code, "test.js", false);

    // then
    assert_eq!(comments.len(), 1);
    assert!(comments[0].text.contains("block comment"));
    assert_eq!(comments[0].comment_type, CommentType::Block);
}

#[test]
fn detect_given_c_block_comment_returns_block_comment() {
    // given
    let detector = CommentDetector::new();
    let code = "/* C block comment */\nint main() { return 0; }";

    // when
    let comments = detector.detect(code, "main.c", false);

    // then
    assert_eq!(comments.len(), 1);
    assert!(comments[0].text.contains("C block comment"));
    assert_eq!(comments[0].comment_type, CommentType::Block);
}
