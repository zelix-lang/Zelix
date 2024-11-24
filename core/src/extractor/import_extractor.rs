use std::{path::PathBuf, process::exit};

use shared::{logger::{Logger, LoggerImpl}, path::retrieve_path, token::{token::{Token, TokenImpl}, token_type::TokenType}};
use crate::shared::import::{Import, Importable};
use super::{sentence_extractor::extract_sentence, standard_locator::locate_standard, token_splitter::split_tokens};

fn throw_invalid_import(details: &[&str]) {
    Logger::err(
        "Invalid import!",
        &[
            "Your import is invalid!"
        ],
        details
    );

    exit(1);
}

pub fn extract_import(tokens: Vec<Token>) -> Vec<Import> {
    let import_tokens : Vec<Token> = extract_sentence(tokens.clone(), TokenType::Semicolon);
    let mut imports : Vec<Import> = Vec::new();

    // A valid import should have a "From" keyword
    // example: "import primtln from ..."
    let before_from = extract_sentence(import_tokens.clone(), TokenType::From);
    // No need to extract again, just get the rest of the tokens
    let after_from = &import_tokens[before_from.len() + 1..];

    // Everything after "From" should be exactly 1 string literal
    if after_from.len() != 1 || after_from[0].get_token_type() != TokenType::StringLiteral {
        let trace: Token;

        if after_from.len() == 0 {
            trace = import_tokens[import_tokens.len() - 1].clone();
        } else {
            trace = after_from[0].clone();
        }

        throw_invalid_import(    
            &[
                "Expected to find a string literal after 'From' keyword!",
                trace.build_trace().as_str()
            ]
        );
    }

    let mut from_raw = after_from[0].get_value();
    let from : PathBuf;

    // Check for standard libraries
    if from_raw.starts_with("@Surf:standard/") {
        from_raw = from_raw.replacen("@Surf:standard/", "", 1);
        from_raw.push_str(".h");

        from = locate_standard(from_raw);
    } else {
        if !from_raw.ends_with(".h") && !from_raw.ends_with(".surf") {
            from_raw.push_str(".surf");
        }

        from = retrieve_path(PathBuf::from(from_raw.clone()));
    }

    let split = split_tokens(&before_from, TokenType::Comma);

    for token in split {
        // This should be an unknown token (not a keyword)
        // and since we're splitting, the length of the vec
        // should be 1

        if token.len() != 1 || token[0].get_token_type() != TokenType::Unknown {
            throw_invalid_import(
                &[
                    "Expected to find an unknown token!",
                    token[0].build_trace().as_str()
                ]
            );
        }

        imports.push(
            Import::new(
                token[0].get_value().clone(),
                from.clone()
            )
        );
    }

    imports
}