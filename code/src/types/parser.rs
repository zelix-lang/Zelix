use std::process::exit;
use logger::{Logger, LoggerImpl};

use crate::{token::{Token, TokenImpl}, token_type::TokenType};

use super::ParamType;

fn assert_token_type(raw: &Vec<Token>, index: usize, check_for: &TokenType) {
    let token = &raw[index];
    let token_type = token.get_token_type();

    if token_type != *check_for {
        Logger::err(
            "Invalid parametrized type",
            &[
                format!(
                    "Expecting a {:?}, but found a {:?}",
                    check_for,
                    token_type
                ).as_str(),
                "This token was unexpected at this time",
                "Define parametrized types like: Result<str>"
            ],
            &[
                token.build_trace().as_str()
            ]
        );

        exit(1);
    }
}

/// Parses a parametrized type into a ParamType struct
/// Example: Result<str> -> ParamType { name: "Result", params: [ParamType { name: "str", params: [] }] }
pub fn parse_parametrized_type(raw: &Vec<Token>) -> ParamType {
    // Handle the simple case directly
    if raw.len() == 1 {
        return ParamType {
            name: raw[0].get_value(),
            params: Vec::new(),
            raw_tokens: raw.to_vec(),
            is_reference: false
        };
    }

    // Validate token structure
    if raw.len() < 4 {
        Logger::err(
            "Invalid parametrized type",
            &["Use parameters in types like: Result<str>"],
            &[raw[0].build_trace().as_str()],
        );
        exit(1);
    }
    
    let is_reference = raw[0].get_token_type() == TokenType::Ampersand;
    let mut raw_clone = raw.clone();

    if is_reference {
        raw_clone.remove(0);
    }

    // Validate required token types in one traversal
    assert_token_type(&raw_clone, 0, &TokenType::Unknown);
    assert_token_type(&raw_clone, 1, &TokenType::LessThan);
    assert_token_type(&raw_clone, &raw_clone.len() - 1, &TokenType::GreaterThan);

    // See if the type is a reference

    // Initialize parameters and extract inner tokens
    let name = &raw_clone[0].get_value();
    let mut params = Vec::new();

    // One-pass parsing for inner tokens
    // This avoids recursion which will cause a stack overflow
    // given a large enough input
    let mut depth = 0;
    let mut start = 2; // Start after '<'
    for (i, token) in raw_clone.iter().enumerate().skip(2).take(&raw_clone.len() - 3) {
        match token.get_token_type() {
            TokenType::LessThan => depth += 1,
            TokenType::GreaterThan => depth -= 1,
            TokenType::Comma if depth == 0 => {
                // Push parameter when depth is 0
                params.push(parse_parametrized_type(&raw_clone[start..i].to_vec()));
                start = i + 1;
            }
            _ => {}
        }
    }

    // Add the last parameter
    if start < raw_clone.len() - 1 {
        params.push(parse_parametrized_type(&raw_clone[start..&raw_clone.len() - 1].to_vec()));
    }

    // Add the ampersand at the end if it's a reference
    // this is done to match C++ syntax
    if is_reference {
        raw_clone.push(raw[0].clone());
    }

    ParamType { name: name.clone(), params, raw_tokens: raw_clone.to_vec(), is_reference }
}
