#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct LanguageSpec {
    pub name: &'static str,
    pub comment_query: &'static str,
    pub docstring_query: Option<&'static str>,
}

#[must_use]
pub fn language_for_path(file_path: &str) -> Option<LanguageSpec> {
    let ext = extension_or_basename(file_path)?;
    language_for_extension(ext)
}

#[must_use]
pub fn language_for_extension(extension: &str) -> Option<LanguageSpec> {
    let ext = extension.trim_start_matches('.');
    let name = language_name_for_alias(ext).or_else(|| {
        ext.bytes()
            .any(|byte| byte.is_ascii_uppercase())
            .then(|| language_name_for_alias(&ext.to_ascii_lowercase()))
            .flatten()
    })?;
    Some(LanguageSpec {
        name,
        comment_query: comment_query(name),
        docstring_query: docstring_query(name),
    })
}

fn language_name_for_alias(extension: &str) -> Option<&'static str> {
    match extension {
        "py" => "python",
        "js" | "jsx" => "javascript",
        "ts" => "typescript",
        "tsx" => "tsx",
        "go" => "go",
        "java" => "java",
        "kt" => "kotlin",
        "scala" => "scala",
        "c" | "h" => "c",
        "cpp" | "cc" | "cxx" | "hpp" => "cpp",
        "rs" => "rust",
        "rb" => "ruby",
        "sh" | "bash" => "bash",
        "cs" => "csharp",
        "swift" => "swift",
        "ex" | "exs" => "elixir",
        "lua" => "lua",
        "php" => "php",
        "ml" | "mli" => "ocaml",
        "sql" => "sql",
        "html" | "htm" => "html",
        "css" => "css",
        "yaml" | "yml" => "yaml",
        "toml" => "toml",
        "hcl" | "tf" => "hcl",
        "dockerfile" => "dockerfile",
        "proto" => "proto",
        "svelte" => "svelte",
        "elm" => "elm",
        "groovy" => "groovy",
        "cue" => "cue",
        _ => return None,
    }
    .into()
}

fn extension_or_basename(file_path: &str) -> Option<&str> {
    let name = file_path.rsplit('/').next().unwrap_or(file_path);
    if let Some((_, ext)) = name.rsplit_once('.') {
        if !ext.is_empty() {
            return Some(ext);
        }
    }
    (!name.is_empty()).then_some(name)
}

fn comment_query(name: &str) -> &'static str {
    match name {
        "rust" | "java" => "(line_comment) @comment\n(block_comment) @comment",
        "kotlin" => "(line_comment) @comment\n(multiline_comment) @comment",
        _ => "(comment) @comment",
    }
}

fn docstring_query(name: &str) -> Option<&'static str> {
    match name {
        "python" => Some(
            r"
		(module . (expression_statement (string) @docstring))
		(class_definition body: (block . (expression_statement (string) @docstring)))
		(function_definition body: (block . (expression_statement (string) @docstring)))
	",
        ),
        "javascript" | "typescript" | "tsx" => Some(
            r#"
		(comment) @jsdoc
		(#match? @jsdoc "^/\*\*")
	"#,
        ),
        "java" => Some(
            r#"
		(comment) @javadoc
		(#match? @javadoc "^/\*\*")
	"#,
        ),
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::{language_for_extension, language_for_path};

    #[test]
    fn language_for_path_supports_dockerfile_basename() {
        // given
        let file_path = "Dockerfile";

        // when
        let language = language_for_path(file_path);

        // then
        assert_eq!(language.map(|spec| spec.name), Some("dockerfile"));
    }

    #[test]
    fn language_for_extension_supports_all_go_aliases() {
        // given
        let aliases = [
            "py",
            "js",
            "jsx",
            "ts",
            "tsx",
            "go",
            "java",
            "kt",
            "scala",
            "c",
            "h",
            "cpp",
            "cc",
            "cxx",
            "hpp",
            "rs",
            "rb",
            "sh",
            "bash",
            "cs",
            "swift",
            "ex",
            "exs",
            "lua",
            "php",
            "ml",
            "mli",
            "sql",
            "html",
            "htm",
            "css",
            "yaml",
            "yml",
            "toml",
            "hcl",
            "tf",
            "dockerfile",
            "proto",
            "svelte",
            "elm",
            "groovy",
            "cue",
        ];

        // when
        let missing: Vec<&str> = aliases
            .iter()
            .copied()
            .filter(|alias| language_for_extension(alias).is_none())
            .collect();

        // then
        assert!(missing.is_empty(), "missing aliases: {missing:?}");
    }

    #[test]
    #[cfg_attr(
        miri,
        ignore = "tree-sitter language pack calls generated C parser FFI"
    )]
    fn language_for_extension_given_supported_aliases_loads_tree_sitter_languages() {
        // given
        let aliases = [
            "py",
            "js",
            "jsx",
            "ts",
            "tsx",
            "go",
            "java",
            "kt",
            "scala",
            "c",
            "h",
            "cpp",
            "cc",
            "cxx",
            "hpp",
            "rs",
            "rb",
            "sh",
            "bash",
            "cs",
            "swift",
            "ex",
            "exs",
            "lua",
            "php",
            "ml",
            "mli",
            "sql",
            "html",
            "htm",
            "css",
            "yaml",
            "yml",
            "toml",
            "hcl",
            "tf",
            "dockerfile",
            "proto",
            "svelte",
            "elm",
            "groovy",
            "cue",
        ];

        // when
        let unloaded: Vec<&str> = aliases
            .iter()
            .copied()
            .filter(|alias| {
                let Some(spec) = language_for_extension(alias) else {
                    return true;
                };
                tree_sitter_language_pack::get_language(spec.name).is_err()
            })
            .collect();

        // then
        assert!(unloaded.is_empty(), "unloaded languages: {unloaded:?}");
    }
}
