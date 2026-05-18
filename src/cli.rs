use clap::Parser;

use crate::api::{EXIT_PASS, check_hook_input};

#[must_use]
pub fn run_from_env() -> i32 {
    let args = Args::parse();
    run(&args)
}

#[must_use]
pub fn run(args: &Args) -> i32 {
    let mut input = String::new();
    if std::io::Read::read_to_string(&mut std::io::stdin(), &mut input).is_err() {
        eprintln!("[check-comments] Skipping: Failed to read stdin");
        return EXIT_PASS;
    }
    run_with_input(&input, args.prompt.as_deref().unwrap_or_default())
}

#[derive(Debug, Parser)]
#[command(
    name = "comment-checker",
    about = "Check for problematic comments in source code",
    long_about = "A hook for Claude Code that detects and warns about comments and docstrings in source code."
)]
pub struct Args {
    #[arg(long)]
    pub prompt: Option<String>,
    #[arg()]
    pub command: Option<String>,
}

#[must_use]
pub fn run_with_input(input: &str, custom_prompt: &str) -> i32 {
    let result = check_hook_input(input, custom_prompt);
    eprint!("{}", result.message);
    result.exit_code
}
