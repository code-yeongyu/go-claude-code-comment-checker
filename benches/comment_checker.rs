use comment_checker::api::check_hook_input;
use comment_checker::detector::CommentDetector;
use comment_checker::filters::AgentMemoFilter;
use comment_checker::models::{CommentInfo, CommentType};
use comment_checker::output::{build_comments_xml, format_hook_message};
use criterion::{Criterion, criterion_group, criterion_main};
use std::hint::black_box;

const PYTHON_SOURCE: &str = "\"\"\"Module docstring.\"\"\"\n\n# regular comment\ndef add(a, b):\n    \"\"\"Function docstring.\"\"\"\n    return a + b\n";

fn benchmark_detector(criterion: &mut Criterion) {
    let detector = CommentDetector::new();
    criterion.bench_function("detect_python_no_docstrings", |bencher| {
        bencher.iter(|| {
            let comments = detector.detect(black_box(PYTHON_SOURCE), black_box("sample.py"), false);
            black_box(comments)
        });
    });
    criterion.bench_function("detect_python_with_docstrings", |bencher| {
        bencher.iter(|| {
            let comments = detector.detect(black_box(PYTHON_SOURCE), black_box("sample.py"), true);
            black_box(comments)
        });
    });
    criterion.bench_function("detect_unsupported_extension", |bencher| {
        bencher.iter(|| {
            let comments =
                detector.detect(black_box(PYTHON_SOURCE), black_box("sample.unknown"), true);
            black_box(comments)
        });
    });
}

fn benchmark_agent_memo(criterion: &mut Criterion) {
    let filter = AgentMemoFilter;
    let early = CommentInfo::new(
        "// Changed from old to new",
        1,
        "sample.py",
        CommentType::Line,
        false,
    );
    let late = CommentInfo::new(
        "// Switched from old to new",
        1,
        "sample.py",
        CommentType::Line,
        false,
    );
    let miss = CommentInfo::new(
        "// calculate checksum before write",
        1,
        "sample.py",
        CommentType::Line,
        false,
    );
    criterion.bench_function("agent_memo_early_hit", |bencher| {
        bencher.iter(|| black_box(filter.is_agent_memo(black_box(&early))));
    });
    criterion.bench_function("agent_memo_late_hit", |bencher| {
        bencher.iter(|| black_box(filter.is_agent_memo(black_box(&late))));
    });
    criterion.bench_function("agent_memo_full_miss", |bencher| {
        bencher.iter(|| black_box(filter.is_agent_memo(black_box(&miss))));
    });
}

fn benchmark_output(criterion: &mut Criterion) {
    let comments: Vec<CommentInfo> = (0..200)
        .map(|index| {
            CommentInfo::new(
                format!("# comment {index}"),
                index + 1,
                "sample.py",
                CommentType::Line,
                false,
            )
        })
        .collect();
    criterion.bench_function("build_comments_xml_many_comments", |bencher| {
        bencher.iter(|| {
            black_box(build_comments_xml(
                black_box(&comments),
                black_box("sample.py"),
            ))
        });
    });
    criterion.bench_function("format_hook_message_many_comments", |bencher| {
        bencher.iter(|| black_box(format_hook_message(black_box(&comments), black_box(""))));
    });
}

fn benchmark_hook_api(criterion: &mut Criterion) {
    let edit_input = r##"{"tool_name":"Edit","tool_input":{"file_path":"sample.py","old_string":"print(1)","new_string":"# new comment\nprint(1)"}}"##;
    let multiedit_input = r##"{"tool_name":"MultiEdit","tool_input":{"file_path":"sample.py","edits":[{"old_string":"print(1)","new_string":"# new comment\nprint(1)"},{"old_string":"print(2)","new_string":"# second comment\nprint(2)"}]}}"##;

    criterion.bench_function("check_hook_input_edit_new_comment", |bencher| {
        bencher.iter(|| {
            black_box(check_hook_input(
                black_box(edit_input),
                black_box("{{comments}}"),
            ))
        });
    });
    criterion.bench_function("check_hook_input_multiedit_new_comments", |bencher| {
        bencher.iter(|| {
            black_box(check_hook_input(
                black_box(multiedit_input),
                black_box("{{comments}}"),
            ))
        });
    });
}

criterion_group!(
    benches,
    benchmark_detector,
    benchmark_agent_memo,
    benchmark_output,
    benchmark_hook_api
);
criterion_main!(benches);
