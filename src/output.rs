use std::collections::HashMap;

use crate::filters::AgentMemoFilter;
use crate::models::CommentInfo;

#[must_use]
pub fn build_comments_xml(comments: &[CommentInfo], file_path: &str) -> String {
    if comments.is_empty() {
        return String::new();
    }

    let mut output = String::with_capacity(24 + file_path.len() + comments.len() * 48);
    push_comments_xml_header(&mut output, file_path);
    let mut line_number = itoa::Buffer::new();
    for comment in comments {
        push_comment_xml(&mut output, comment, &mut line_number);
    }
    output.push_str("</comments>");
    output
}

fn push_comments_xml_header(output: &mut String, file_path: &str) {
    output.push_str("<comments file=\"");
    output.push_str(file_path);
    output.push_str("\">\n");
}

fn push_comment_xml(output: &mut String, comment: &CommentInfo, line_number: &mut itoa::Buffer) {
    output.push_str("\t<comment line-number=\"");
    output.push_str(line_number.format(comment.line_number));
    output.push_str("\">");
    output.push_str(&comment.text);
    output.push_str("</comment>\n");
}

#[must_use]
pub fn format_hook_message(comments: &[CommentInfo], custom_prompt: &str) -> String {
    if comments.is_empty() {
        return String::new();
    }

    let mut group_indexes: HashMap<&str, usize> = HashMap::with_capacity(comments.len());
    let mut groups: Vec<(&str, Vec<&CommentInfo>)> = Vec::new();
    for comment in comments {
        let file_path = comment.file_path.as_str();
        match group_indexes.entry(file_path) {
            std::collections::hash_map::Entry::Occupied(entry) => {
                groups[*entry.get()].1.push(comment);
            }
            std::collections::hash_map::Entry::Vacant(entry) => {
                let index = groups.len();
                entry.insert(index);
                groups.push((file_path, vec![comment]));
            }
        }
    }

    let mut comments_xml = String::with_capacity(comments.len() * 64);
    for (file_path, file_comments) in &groups {
        push_comments_xml_header(&mut comments_xml, file_path);
        let mut line_number = itoa::Buffer::new();
        for comment in file_comments {
            push_comment_xml(&mut comments_xml, comment, &mut line_number);
        }
        comments_xml.push_str("</comments>");
        comments_xml.push('\n');
    }

    if !custom_prompt.is_empty() {
        return custom_prompt.replace("{{comments}}", &comments_xml);
    }

    let agent_memo = AgentMemoFilter;
    let agent_memo_comments: Vec<&CommentInfo> = comments
        .iter()
        .filter(|comment| agent_memo.is_agent_memo(comment))
        .collect();
    let has_agent_memo = !agent_memo_comments.is_empty();

    let mut output = String::new();
    if has_agent_memo {
        output.push_str("🚨 AGENT MEMO COMMENT DETECTED - CODE SMELL ALERT 🚨\n\n");
    } else {
        output.push_str("COMMENT/DOCSTRING DETECTED - IMMEDIATE ACTION REQUIRED\n\n");
    }

    if has_agent_memo {
        output.push_str("⚠️  AGENT MEMO COMMENTS DETECTED - THIS IS A CODE SMELL  ⚠️\n\n");
        output.push_str("You left \"memo-style\" comments that describe WHAT you changed or HOW you implemented something.\n");
        output.push_str(
            "These are typically signs of an AI agent leaving notes for itself or the user.\n\n",
        );
        output.push_str("Examples of agent memo patterns detected:\n");
        output.push_str("  - \"Changed from X to Y\", \"Modified to...\", \"Updated from...\"\n");
        output.push_str("  - \"Added new...\", \"Removed...\", \"Refactored...\"\n");
        output.push_str("  - \"This implements...\", \"Here we...\", \"Now this...\"\n");
        output.push_str("  - \"Note:\", \"Implementation of...\"\n");
        output.push_str("WHY THIS IS BAD:\n");
        output.push_str("  1. Code should be self-documenting - if you need to explain what changed, the code isn't clear enough\n");
        output
            .push_str("  2. These comments become outdated instantly and mislead future readers\n");
        output.push_str(
            "  3. Git history already tracks what changed - comments duplicate this poorly\n",
        );
        output.push_str("  4. It's a sign of rushing without proper refactoring\n\n");
        output.push_str("ACTION REQUIRED:\n");
        output.push_str("  -> REMOVE these memo comments immediately\n");
        output.push_str("  -> If the code needs explanation, refactor it to be clearer instead\n");
        output.push_str("  -> Use meaningful variable/function names that explain intent\n");
        output.push_str("  -> Let git commit messages document the \"what\" and \"why\"\n\n");
        output.push_str("Detected agent memo comments:\n");
        let mut line_number = itoa::Buffer::new();
        for memo in agent_memo_comments {
            output.push_str("  - Line ");
            output.push_str(line_number.format(memo.line_number));
            output.push_str(": ");
            output.push_str(memo.text.trim());
            output.push('\n');
        }
        output.push_str("\n---\n\n");
    }

    output.push_str(
        "Your recent changes contain comments or docstrings, which triggered this hook.\n",
    );
    output.push_str("You need to take immediate action. You must follow the conditions below.\n");
    output.push_str(
        "(Listed in priority order - you must always act according to this priority order)\n\n",
    );
    output.push_str("CRITICAL WARNING: This hook message MUST NEVER be ignored, even if you receive it multiple times.\n");
    output.push_str("You MUST provide corresponding explanation or action for EACH occurrence of this message.\n");
    output.push_str(
        "Ignoring this message or failing to respond appropriately is strictly prohibited.\n\n",
    );
    output.push_str("PRIORITY-BASED ACTION GUIDELINES:\n\n");
    output.push_str("1. This is a comment/docstring that already existed before\n");
    output.push_str("\t-> Explain to the user that this is an existing comment/docstring and proceed (justify it)\n\n");
    output.push_str("2. This is a newly written comment: but it's in given, when, then format\n");
    output.push_str("\t-> Tell the user it's a BDD comment and proceed (justify it)\n");
    output.push_str("\t-> Note: This applies to comments only, not docstrings\n\n");
    output.push_str(
        "3. This is a newly written comment/docstring: but it's a necessary comment/docstring\n",
    );
    output.push_str("\t-> Tell the user why this comment/docstring is absolutely necessary and proceed (justify it)\n");
    output.push_str("\t-> Examples of necessary comments: complex algorithms, security-related, performance optimization, regex, mathematical formulas\n");
    output.push_str("\t-> Examples of necessary docstrings: public API documentation, complex module/class interfaces\n");
    output.push_str("\t-> IMPORTANT: Most docstrings are unnecessary if the code is self-explanatory. Only keep truly essential ones.\n\n");
    output.push_str(
        "4. This is a newly written comment/docstring: but it's an unnecessary comment/docstring\n",
    );
    output.push_str("\t-> Apologize to the user and remove the comment/docstring.\n");
    output.push_str(
        "\t-> Make the code itself clearer so it can be understood without comments/docstrings.\n",
    );
    output.push_str("\t-> For verbose docstrings: refactor code to be self-documenting instead of adding lengthy explanations.\n\n");
    output.push_str("MANDATORY REQUIREMENT: You must acknowledge this hook message and take one of the above actions.\n");
    output.push_str("Review in the above priority order and take the corresponding action EVERY TIME this appears.\n\n");
    output.push_str("REMINDER: These rules apply to ALL your future code, not just this specific edit. Always be deliberate and cautious when writing comments - only add them when absolutely necessary.\n\n");
    output.push_str("Detected comments/docstrings:\n");
    output.push_str(&comments_xml);
    output
}
