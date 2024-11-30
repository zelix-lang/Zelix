use std::{path::PathBuf, process::exit};

use shared::{logger::{Logger, LoggerImpl}, path::retrieve_path, token::{token::{Token, TokenImpl}, token_type::TokenType}};
use shared::code::import::{Import, Importable};
use super::{sentence_extractor::extract_sentence, standard_locator::locate_standard};

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

pub fn extract_import(tokens: Vec<Token>) -> Import {
    let import_tokens : Vec<Token> = extract_sentence(tokens.clone(), TokenType::Semicolon);

    // A valid import should have be only 2 tokens long
    // @import "file";
    // But we don't receive the @import token so we just have to check
    // if the length is 1

    if import_tokens.len() != 1 {
        let trace;

        if import_tokens.len() == 0 {
            // Tokens always have at least one token
            trace = tokens[0].build_trace();
        } else {
            trace = import_tokens[0].build_trace();
        }

        throw_invalid_import(
            &[
                "Invalid import syntax!",
                "Imports should be in the format: @import \"file\";",
                trace.as_str()
            ],
        );
    }

    let mut from_raw = import_tokens[0].get_value();
    let from : PathBuf;

    // Check for standard libraries
    if from_raw.starts_with("@Surf:standard/") {
        from_raw = from_raw.replacen("@Surf:standard/", "", 1);
        from_raw.push_str(".h");

        // Check if the extension is .hpp instead of .h
        // Some libraries use .hpp as they contain pure C++ code
        // instead of C code
        if !PathBuf::from(from_raw.clone()).exists() {
            from_raw = from_raw.strip_suffix(".h").unwrap().to_string();
            from_raw.push_str(".hpp");
        }

        from = locate_standard(from_raw);
    } else {
        if !from_raw.ends_with(".h") &&
            !from_raw.ends_with(".hpp") &&
            !from_raw.ends_with(".surf")
        {
            from_raw.push_str(".surf");
        }

        from = retrieve_path(PathBuf::from(from_raw.clone()));
    }
    
    Import::new(
        from.clone(),
        import_tokens[0].build_trace()
    )
}