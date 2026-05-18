use comment_checker::models::{CommentInfo, CommentType};
use comment_checker::output::{build_comments_xml, format_hook_message};

fn comment(text: &str, line_number: usize, file_path: &str) -> CommentInfo {
    CommentInfo::new(text, line_number, file_path, CommentType::Line, false)
}

#[test]
fn build_comments_xml_given_single_comment_returns_exact_shape() {
    // given
    let comments = vec![comment("# TODO: fix this", 10, "src/app.py")];

    // when
    let result = build_comments_xml(&comments, "src/app.py");

    // then
    assert_eq!(
        result,
        "<comments file=\"src/app.py\">\n\t<comment line-number=\"10\"># TODO: fix this</comment>\n</comments>",
    );
}

#[test]
fn build_comments_xml_given_empty_comments_returns_empty_string() {
    // given
    let comments = Vec::new();

    // when
    let result = build_comments_xml(&comments, "src/app.py");

    // then
    assert_eq!(result, "");
}

#[test]
fn format_hook_message_given_single_comment_contains_default_warning() {
    // given
    let comments = vec![comment("# TODO: fix this", 10, "src/app.py")];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED"));
    assert!(result.contains("<comments file=\"src/app.py\">"));
    assert!(result.contains("<comment line-number=\"10\"># TODO: fix this</comment>"));
    assert!(result.contains("PRIORITY-BASED ACTION GUIDELINES:"));
}

#[test]
fn format_hook_message_given_multiple_comments_single_file_groups_xml() {
    // given
    let comments = vec![
        comment("# First comment", 5, "src/main.py"),
        comment("# Second comment", 15, "src/main.py"),
    ];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("<comments file=\"src/main.py\">"));
    assert!(result.contains("<comment line-number=\"5\"># First comment</comment>"));
    assert!(result.contains("<comment line-number=\"15\"># Second comment</comment>"));
    assert_eq!(result.matches("<comments file=\"src/main.py\">").count(), 1);
    assert_eq!(result.matches("</comments>").count(), 1);
}

#[test]
fn format_hook_message_given_comments_multiple_files_returns_separate_xml_blocks() {
    // given
    let comments = vec![
        comment("# Comment in file 1", 10, "src/file1.py"),
        comment("# Comment in file 2", 20, "src/file2.py"),
    ];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("<comments file=\"src/file1.py\">"));
    assert!(result.contains("<comments file=\"src/file2.py\">"));
    assert!(result.contains("<comment line-number=\"10\"># Comment in file 1</comment>"));
    assert!(result.contains("<comment line-number=\"20\"># Comment in file 2</comment>"));
    assert_eq!(result.matches("</comments>").count(), 2);
}

#[test]
fn format_hook_message_given_empty_list_returns_empty_string() {
    // given
    let comments = Vec::new();

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert_eq!(result, "");
}

#[test]
fn format_hook_message_given_docstring_comment_returns_formatted_message() {
    // given
    let comments = vec![CommentInfo::new(
        "\"\"\"Module docstring.\"\"\"",
        1,
        "src/utils.py",
        CommentType::Docstring,
        true,
    )];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED"));
    assert!(result.contains("<comments file=\"src/utils.py\">"));
    assert!(result.contains(r#"<comment line-number="1">"""Module docstring."""</comment>"#));
    assert!(result.contains("MANDATORY REQUIREMENT:"));
}

#[test]
fn format_hook_message_given_comment_uses_tabs_for_xml_indentation() {
    // given
    let comments = vec![comment("# Comment", 5, "src/test.py")];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("\t<comment line-number=\"5\"># Comment</comment>"));
}

#[test]
fn format_hook_message_given_custom_prompt_replaces_placeholder() {
    // given
    let comments = vec![comment("# Test comment", 10, "src/app.py")];
    let custom_prompt = "CUSTOM WARNING\n{{comments}}\nPlease fix.";

    // when
    let result = format_hook_message(&comments, custom_prompt);

    // then
    assert!(result.contains("CUSTOM WARNING"));
    assert!(result.contains("<comments file=\"src/app.py\">"));
    assert!(result.contains("Please fix."));
    assert!(!result.contains("COMMENT/DOCSTRING DETECTED"));
}

#[test]
fn format_hook_message_given_custom_prompt_without_placeholder_returns_custom_only() {
    // given
    let comments = vec![comment("# Test", 1, "test.py")];
    let custom_prompt = "Simple warning without placeholder.";

    // when
    let result = format_hook_message(&comments, custom_prompt);

    // then
    assert_eq!(result, "Simple warning without placeholder.");
}

#[test]
fn format_hook_message_given_agent_memo_contains_memo_warning() {
    // given
    let comments = vec![comment(
        "# Modified to use new implementation",
        1,
        "src/app.py",
    )];

    // when
    let result = format_hook_message(&comments, "");

    // then
    assert!(result.contains("AGENT MEMO"));
    assert!(result.contains("CODE SMELL"));
}
