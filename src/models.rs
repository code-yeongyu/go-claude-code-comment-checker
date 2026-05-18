use std::collections::BTreeMap;

use serde::Serialize;

#[derive(Clone, Copy, Debug, Eq, PartialEq, Serialize)]
#[serde(rename_all = "lowercase")]
pub enum CommentType {
    Line,
    Block,
    Docstring,
}

#[derive(Clone, Debug, Eq, PartialEq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct CommentInfo {
    pub text: String,
    pub line_number: usize,
    pub file_path: String,
    pub comment_type: CommentType,
    pub is_docstring: bool,
    pub metadata: BTreeMap<String, String>,
}

impl CommentInfo {
    #[must_use]
    pub fn new(
        text: impl Into<String>,
        line_number: usize,
        file_path: impl Into<String>,
        comment_type: CommentType,
        is_docstring: bool,
    ) -> Self {
        Self {
            text: text.into(),
            line_number,
            file_path: file_path.into(),
            comment_type,
            is_docstring,
            metadata: BTreeMap::new(),
        }
    }

    #[must_use]
    pub fn normalized_text(&self) -> String {
        self.text.trim().to_lowercase()
    }
}
