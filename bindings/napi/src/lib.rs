use napi_derive::napi;

#[napi(js_name = "version")]
pub fn version() -> String {
    env!("CARGO_PKG_VERSION").to_owned()
}

// CLIPPY-ALLOW: N-API string parameters must be owned JavaScript primitives.
#[allow(clippy::needless_pass_by_value)]
#[napi(js_name = "checkHookJson")]
pub fn check_hook_json(input: String, custom_prompt: Option<String>) -> String {
    let prompt = custom_prompt.unwrap_or_default();
    comment_checker::api::check_hook_input_json(&input, &prompt)
}

// CLIPPY-ALLOW: N-API string parameters must be owned JavaScript primitives.
#[allow(clippy::needless_pass_by_value)]
#[napi(js_name = "detectCommentsJson")]
pub fn detect_comments_json(
    content: String,
    file_path: String,
    include_docstrings: bool,
) -> String {
    comment_checker::api::detect_comments_json(&content, &file_path, include_docstrings)
}
