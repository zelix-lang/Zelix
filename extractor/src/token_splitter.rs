use std::process::exit;

use code::{token::{Token, TokenImpl}, token_type::TokenType};
use logger::{Logger, LoggerImpl};

/// Splits tokens into different vectors using the provided identifier
pub fn split_tokens(
    tokens: &Vec<Token>,
    identifier: &TokenType,
    nested_open_identifier: &TokenType,
    nested_close_identifier: &TokenType
) -> Vec<Vec<Token>> {
    let mut result: Vec<Vec<Token>> = Vec::new();
    let mut current: Vec<Token> = Vec::new();
    let mut nested_level : isize = 0;

    for token in tokens {
        if token.get_token_type() == *nested_open_identifier {
            nested_level += 1;
        } else if token.get_token_type() == *nested_close_identifier {
            nested_level -= 1;

            if nested_level < 0 {
                Logger::err(
                    "Invalid parametrized type",
                    &[
                        "Unexpected closing token"
                    ],
                    &[
                        token.build_trace().as_str()
                    ]
                );

                exit(1);
            }
        }

        if token.get_token_type() == *identifier {
            if nested_level > 0 {
                current.push(token.clone());
                continue;
            }

            result.push(current.clone());

            // Instead of reassigning, we clear the vector
            // which is more efficient
            current.clear();
        } else {
            current.push(token.clone());
        }
    }

    // Push the remaining leftovers
    result.push(current);
    result
}

/// Extracts all tokens before the provided delimiter
pub fn extract_tokens_before(
    tokens: &Vec<Token>,
    delimiter: &TokenType
) -> Vec<Token> {
    let mut result: Vec<Token> = Vec::new();
    let mut has_met_delimiter = false;

    for token in tokens {
        if token.get_token_type() == *delimiter {
            has_met_delimiter = true;
            break;
        }

        result.push(token.clone());
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

    result
}