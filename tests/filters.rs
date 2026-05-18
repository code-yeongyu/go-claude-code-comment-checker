use comment_checker::filters::{
    AgentMemoFilter, BddFilter, DirectiveFilter, ShebangFilter, apply_filters,
};
use comment_checker::models::{CommentInfo, CommentType};

fn line(text: &str) -> CommentInfo {
    CommentInfo::new(text, 1, "test.py", CommentType::Line, false)
}

#[test]
fn bdd_filter_given_keyword_returns_true() {
    // given
    let filter = BddFilter;
    let comment = line("# given");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn bdd_filter_given_when_then_keyword_returns_true() {
    // given
    let filter = BddFilter;
    let comment = line("// when & then");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn bdd_filter_given_arrange_keyword_returns_true() {
    // given
    let filter = BddFilter;
    let comment = CommentInfo::new("-- arrange", 1, "test.sql", CommentType::Line, false);

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn bdd_filter_given_when_ampersand_then_returns_true() {
    // given
    let filter = BddFilter;
    let comment = line("// when&then");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn bdd_filter_given_regular_comment_returns_false() {
    // given
    let filter = BddFilter;
    let comment = line("# regular comment");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(!result);
}

#[test]
fn directive_filter_given_noqa_returns_true() {
    // given
    let filter = DirectiveFilter;
    let comment = line("# noqa: F401");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn directive_filter_given_ts_ignore_returns_true() {
    // given
    let filter = DirectiveFilter;
    let comment = line("// @ts-ignore");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn directive_filter_given_pyright_returns_true() {
    // given
    let filter = DirectiveFilter;
    let comment = line("# pyright: ignore");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn directive_filter_given_eslint_disable_returns_true() {
    // given
    let filter = DirectiveFilter;
    let comment = line("// eslint-disable-next-line");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn directive_filter_given_type_ignore_returns_true() {
    // given
    let filter = DirectiveFilter;
    let comment = line("# type: ignore");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn directive_filter_given_regular_comment_returns_false() {
    // given
    let filter = DirectiveFilter;
    let comment = line("# regular comment");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(!result);
}

#[test]
fn shebang_filter_given_shebang_returns_true() {
    // given
    let filter = ShebangFilter;
    let comment = line("#!/usr/bin/env python");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn shebang_filter_given_shebang_with_space_returns_true() {
    // given
    let filter = ShebangFilter;
    let comment = line("#! /usr/bin/env python");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn shebang_filter_given_bash_shebang_returns_true() {
    // given
    let filter = ShebangFilter;
    let comment = line("#!/bin/bash");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(result);
}

#[test]
fn shebang_filter_given_regular_comment_returns_false() {
    // given
    let filter = ShebangFilter;
    let comment = line("# regular comment");

    // when
    let result = filter.should_skip(&comment);

    // then
    assert!(!result);
}

#[test]
fn agent_memo_filter_given_changed_from_returns_true() {
    // given
    let filter = AgentMemoFilter;
    let comment = line("// Changed from old_value to new_value");

    // when
    let result = filter.is_agent_memo(&comment);

    // then
    assert!(result);
}

#[test]
fn agent_memo_filter_given_late_switched_pattern_returns_true() {
    // given
    let filter = AgentMemoFilter;
    let comment = line("// Switched from old to new");

    // when
    let result = filter.is_agent_memo(&comment);

    // then
    assert!(result);
}

#[test]
fn agent_memo_filter_given_ported_english_patterns_returns_true() {
    // given
    let filter = AgentMemoFilter;
    let comments = [
        line("# Modified to use new implementation"),
        line("// Updated from v1 to v2"),
        line("// Refactored for better performance"),
        line("// Added new validation logic"),
        line("// Removed deprecated function"),
        line("// Implemented new feature"),
        line("// This implements the new API"),
        line("// Here we handle the error case"),
        line("// Now this uses the new format"),
        line("// Previously this was handled differently"),
        line("// Note: this is important"),
        line("// oldValue -> newValue"),
        line("# Converted from callbacks to state machine"),
        line("# Migrated to the new parser"),
        line("# Switched from old code path"),
        line("# Replaced temporary formatter"),
        line("# Deleted stale branch"),
        line("# Before this the parser was recreated"),
        line("# After this the cache is warm"),
    ];

    // when
    let results: Vec<bool> = comments
        .iter()
        .map(|comment| filter.is_agent_memo(comment))
        .collect();

    // then
    assert!(results.iter().all(|result| *result), "results: {results:?}");
}

#[test]
fn agent_memo_filter_given_allowed_non_memo_patterns_returns_false() {
    // given
    let filter = AgentMemoFilter;
    let comments = [
        line("# given"),
        line("# noqa: E501"),
        line("// Calculate the sum of values"),
    ];

    // when
    let results: Vec<bool> = comments
        .iter()
        .map(|comment| filter.is_agent_memo(comment))
        .collect();

    // then
    assert!(
        results.iter().all(|result| !*result),
        "results: {results:?}"
    );
}

#[test]
fn agent_memo_filter_given_regular_comment_returns_false() {
    // given
    let filter = AgentMemoFilter;
    let comment = line("// Calculate checksum before write");

    // when
    let result = filter.is_agent_memo(&comment);

    // then
    assert!(!result);
}

#[test]
fn apply_filters_given_allowed_comments_returns_empty() {
    // given
    let comments = vec![line("# given"), line("// @ts-ignore"), line("#!/bin/bash")];

    // when
    let result = apply_filters(&comments);

    // then
    assert!(result.is_empty());
}
