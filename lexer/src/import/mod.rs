use std::{path::PathBuf, process::exit};

use code::{token::{Token, TokenImpl}, token_type::TokenType};
use logger::{Logger, LoggerImpl};

/// Checks if a list of tokens contains imports
/// that don't point to the standard library
pub fn find_code_imports(tokens: &Vec<Token>) -> Option<(usize, &Token)> {
    tokens.iter().enumerate().find(|(index, token)| {
        if token.get_token_type() != TokenType::Import {
            return false;
        }

        let has_next = index + 1 < tokens.len();

        if !has_next {
            return false;
        }

        let next_token = &tokens[index + 1];

        if next_token.get_token_type() != TokenType::StringLiteral {
            return false;
        }

        let value = next_token.get_value();
        // Exclude standard library imports
        !value.starts_with("@Surf:standard/")
    })
}

pub fn build_chain_trace(chain: &Vec<PathBuf>) -> Vec<String> {
    let mut trace = Vec::new();

    for n in 0..chain.len() {
        let path = &chain[n];
        let path_str = path.to_str().unwrap();

        let mut line = String::from(" ".repeat(n));
        line.push_str("-> ");
        line.push_str(path_str);
        trace.push(line);
    }

    trace
}

pub fn check_import_is_valid(import: &PathBuf) {
    if !import.exists() || !import.is_file() {
        Logger::err(
            "Invalid import",
            &[
                "Check if the import exists and is a file"
            ],
            &[
                format!(
                    "{:?} does not exist or is not a file",
                    import
                ).as_str()
            ],
        );

        exit(1);
    }
}