use std::path::Path;

#[must_use]
pub fn read_file(path: impl AsRef<Path>) -> String {
    let Ok(data) = std::fs::read(path) else {
        return String::new();
    };
    String::from_utf8(data).unwrap_or_else(|error| latin1_to_string(error.as_bytes()))
}

#[must_use]
pub fn read_string(content: &str) -> String {
    content.to_owned()
}

fn latin1_to_string(bytes: &[u8]) -> String {
    let mut output = String::with_capacity(bytes.len());
    for byte in bytes {
        output.push(char::from(*byte));
    }
    output
}
