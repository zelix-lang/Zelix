use std::process::exit;

use shared::{logger::{Logger, LoggerImpl}, token::{token::{Token, TokenImpl}, token_type::TokenType}};

// Useful to extract a sentence from a list of tokens
// for example:
// import "println"; -> ["import", "println"]
pub fn extract_sentence(
    tokens: Vec<Token>,
    delimiter: TokenType
) -> Vec<Token> {

    // Store all the tokens in the sentence
    let mut sentence: Vec<Token> = Vec::new();
    let mut has_met_delimiter = false;

    for token in tokens {
        if token.get_token_type() == delimiter {
            has_met_delimiter = true;
            break;
        }
        
        sentence.push(token);
    }

    if !has_met_delimiter {
        Logger::err(
            "Invalid sentence!",
            &[
                "Your syntax is invalid!"
            ],
            &[
                format!(
                    "Expected to find a {:?}, but it was never found!",
                    delimiter.clone()
                ).as_str()
            ]
        );

        exit(1);
    }

    sentence

}