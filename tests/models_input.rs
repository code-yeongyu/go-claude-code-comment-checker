use std::fs;

use comment_checker::input::{read_file, read_string};
use comment_checker::models::{CommentInfo, CommentType};

#[test]
fn normalized_text_given_whitespace_returns_stripped_lowercase() {
    // given
    let comment = CommentInfo::new("  # HELLO  ", 1, "test.py", CommentType::Line, false);

    // when
    let normalized = comment.normalized_text();

    // then
    assert_eq!(normalized, "# hello");
}

#[test]
fn read_string_given_content_returns_same_content() {
    // given
    let content = "print('hello')";

    // when
    let result = read_string(content);

    // then
    assert_eq!(result, content);
}

#[test]
fn read_file_given_missing_file_returns_empty_string() {
    // given
    let path = "/tmp/comment-checker-missing-file.py";

    // when
    let result = read_file(path);

    // then
    assert_eq!(result, "");
}

#[test]
fn read_file_given_latin1_file_falls_back_to_latin1() -> std::io::Result<()> {
    // given
    let dir = tempfile::tempdir()?;
    let path = dir.path().join("latin1.py");
    fs::write(&path, [0x63, 0x61, 0x66, 0xe9])?;

    // when
    let result = read_file(&path);

    // then
    assert_eq!(result, "cafe\u{301}".replace("e\u{301}", "é"));
    Ok(())
}
