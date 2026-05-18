use comment_checker::detector::CommentDetector;
use comment_checker::filters::{AgentMemoFilter, apply_filters};
use comment_checker::output::format_hook_message;

#[test]
fn full_pipeline_given_no_comments_returns_empty() {
    // given
    let detector = CommentDetector::new();
    let code = "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}";

    // when
    let comments = detector.detect(code, "main.go", false);

    // then
    assert!(comments.is_empty());
}

#[test]
fn full_pipeline_given_bdd_comment_filters_out() {
    // given
    let detector = CommentDetector::new();
    let code = "# given\ndef test_something():\n    pass";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert!(!comments.is_empty());
    assert!(filtered.is_empty());
}

#[test]
fn full_pipeline_given_regular_comment_returns_formatted_message() {
    // given
    let detector = CommentDetector::new();
    let code = "# This is a regular comment\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);
    let message = format_hook_message(&filtered, "");

    // then
    assert_eq!(filtered.len(), 1);
    assert!(message.contains("COMMENT/DOCSTRING DETECTED"));
    assert!(message.contains("test.py"));
    assert!(message.contains("This is a regular comment"));
}

#[test]
fn full_pipeline_given_docstring_keeps_as_code_smell() {
    // given
    let detector = CommentDetector::new();
    let code = "\"\"\"This is a docstring\"\"\"\ndef hello():\n    pass";

    // when
    let comments = detector.detect(code, "test.py", true);
    let filtered = apply_filters(&comments);

    // then
    assert!(!comments.is_empty());
    assert!(!filtered.is_empty());
}

#[test]
fn full_pipeline_given_directive_filters_out() {
    // given
    let detector = CommentDetector::new();
    let code = "# noqa: E501\nprint(\"very long line\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert!(!comments.is_empty());
    assert!(filtered.is_empty());
}

#[test]
fn full_pipeline_given_shebang_filters_out() {
    // given
    let detector = CommentDetector::new();
    let code = "#!/usr/bin/env python\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert!(!comments.is_empty());
    assert!(filtered.is_empty());
}

#[test]
fn full_pipeline_given_agent_memo_detects_code_smell() {
    // given
    let detector = CommentDetector::new();
    let agent_memo_filter = AgentMemoFilter;
    let code = "# Changed from old_value to new_value\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert_eq!(filtered.len(), 1);
    assert!(agent_memo_filter.is_agent_memo(&filtered[0]));
}

#[test]
fn full_pipeline_given_agent_memo_formatter_includes_warning() {
    // given
    let detector = CommentDetector::new();
    let code = "# Modified to use new implementation\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);
    let message = format_hook_message(&filtered, "");

    // then
    assert!(message.contains("AGENT MEMO"));
    assert!(message.contains("CODE SMELL"));
}

#[test]
fn full_pipeline_given_switched_from_detects_code_smell() {
    // given
    let detector = CommentDetector::new();
    let agent_memo_filter = AgentMemoFilter;
    let code = "# Switched from old to new\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert_eq!(filtered.len(), 1);
    assert!(agent_memo_filter.is_agent_memo(&filtered[0]));
}

#[test]
fn full_pipeline_given_regular_comment_is_not_agent_memo() {
    // given
    let detector = CommentDetector::new();
    let agent_memo_filter = AgentMemoFilter;
    let code = "# Calculate the sum of values\nprint(\"hello\")";

    // when
    let comments = detector.detect(code, "test.py", false);
    let filtered = apply_filters(&comments);

    // then
    assert_eq!(filtered.len(), 1);
    assert!(!agent_memo_filter.is_agent_memo(&filtered[0]));
}
