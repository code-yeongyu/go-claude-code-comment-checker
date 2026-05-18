use std::io::Write;
use std::process::{Command, Stdio};

fn run_checker(input: &str, args: &[&str]) -> std::io::Result<std::process::Output> {
    let mut child = Command::new(env!("CARGO_BIN_EXE_comment-checker"))
        .args(args)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()?;
    let Some(stdin) = child.stdin.as_mut() else {
        return Err(std::io::Error::other("checker stdin unavailable"));
    };
    stdin.write_all(input.as_bytes())?;
    child.wait_with_output()
}

#[test]
fn cli_given_no_comment_exits_zero() -> std::io::Result<()> {
    // given
    let input =
        r#"{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"print(1)"}}"#;

    // when
    let output = run_checker(input, &[])?;

    // then
    assert!(output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Success"));
    Ok(())
}

#[test]
fn cli_given_comment_exits_two() -> std::io::Result<()> {
    // given
    let input = r##"{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"# comment\nprint(1)"}}"##;

    // when
    let output = run_checker(input, &[])?;

    // then
    assert_eq!(output.status.code(), Some(2));
    assert!(String::from_utf8_lossy(&output.stderr).contains("COMMENT/DOCSTRING DETECTED"));
    Ok(())
}

#[test]
fn cli_given_check_alias_and_bdd_comment_exits_zero() -> std::io::Result<()> {
    // given
    let input = r##"{"tool_name":"Write","tool_input":{"file_path":"test.py","content":"# given\nprint(1)"}}"##;

    // when
    let output = run_checker(input, &["check"])?;

    // then
    assert!(output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Success"));
    Ok(())
}

#[test]
fn cli_given_invalid_json_exits_zero_with_skip() -> std::io::Result<()> {
    // given
    let input = "invalid json";

    // when
    let output = run_checker(input, &[])?;

    // then
    assert!(output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Skipping: Invalid input format"));
    Ok(())
}

#[test]
fn cli_given_edit_existing_comment_only_exits_zero() -> std::io::Result<()> {
    // given
    let input = r##"{"tool_name":"Edit","tool_input":{"file_path":"test.py","old_string":"# comment\nx = 1","new_string":"# comment\nx = 2"}}"##;

    // when
    let output = run_checker(input, &[])?;

    // then
    assert!(output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Success"));
    Ok(())
}

#[test]
fn cli_given_multiedit_new_comment_exits_two() -> std::io::Result<()> {
    // given
    let input = r##"{"tool_name":"MultiEdit","tool_input":{"file_path":"test.py","edits":[{"old_string":"x","new_string":"# comment\ny"}]}}"##;

    // when
    let output = run_checker(input, &[])?;

    // then
    assert_eq!(output.status.code(), Some(2));
    Ok(())
}
