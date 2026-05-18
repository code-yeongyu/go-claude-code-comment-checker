use comment_checker::api::{check_hook_input, check_hook_input_json, detect_comments_json};
use serde_json::Value;

#[test]
fn check_hook_input_given_clean_write_returns_pass() {
    // given
    let input =
        r#"{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"print(1)"}}"#;

    // when
    let result = check_hook_input(input, "");

    // then
    assert_eq!(result.exit_code, 0);
    assert!(result.message.contains("Success"));
}

#[test]
fn check_hook_input_json_given_comment_returns_camel_case_block() -> serde_json::Result<()> {
    // given
    let input = r##"{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"# comment\nprint(1)"}}"##;

    // when
    let json: Value = serde_json::from_str(&check_hook_input_json(input, ""))?;

    // then
    assert_eq!(json.get("exitCode").and_then(Value::as_i64), Some(2));
    assert!(
        json.get("message")
            .and_then(Value::as_str)
            .is_some_and(|message| message.contains("<comments file=\"test.py\">"))
    );
    Ok(())
}

#[test]
fn detect_comments_json_given_python_comment_returns_comment_payload() -> serde_json::Result<()> {
    // given
    let content = "# comment\nprint(1)";

    // when
    let json: Value = serde_json::from_str(&detect_comments_json(content, "test.py", true))?;

    // then
    assert_eq!(json.as_array().map(Vec::len), Some(1));
    assert_eq!(
        json.as_array()
            .and_then(|comments| comments.first())
            .and_then(|comment| comment.get("lineNumber"))
            .and_then(Value::as_u64),
        Some(1),
    );
    Ok(())
}
