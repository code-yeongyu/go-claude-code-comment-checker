use wasm_bindgen::prelude::wasm_bindgen;

#[wasm_bindgen(js_name = checkHookJson)]
pub fn check_hook_json(input: &str, custom_prompt: Option<String>) -> String {
    let prompt = custom_prompt.unwrap_or_default();
    comment_checker::api::check_hook_input_json(input, &prompt)
}

#[wasm_bindgen(js_name = detectCommentsJson)]
pub fn detect_comments_json(content: &str, file_path: &str, include_docstrings: bool) -> String {
    comment_checker::api::detect_comments_json(content, file_path, include_docstrings)
}
