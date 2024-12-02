use std::process::exit;

use shared::{logger::{Logger, LoggerImpl}, token::{token::{Token, TokenImpl}, token_type::TokenType}};

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

    // Validate required token types in one traversal
    assert_token_type(raw, 0, &TokenType::Unknown);
    assert_token_type(raw, 1, &TokenType::LessThan);
    assert_token_type(raw, raw.len() - 1, &TokenType::GreaterThan);

    // Initialize parameters and extract inner tokens
    let name = raw[0].get_value();
    let mut params = Vec::new();

    // One-pass parsing for inner tokens
    // This avoids recursion which will cause a stack overflow
    // given a large enough input
    let mut depth = 0;
    let mut start = 2; // Start after '<'
    for (i, token) in raw.iter().enumerate().skip(2).take(raw.len() - 3) {
        match token.get_token_type() {
            TokenType::LessThan => depth += 1,
            TokenType::GreaterThan => depth -= 1,
            TokenType::Comma if depth == 0 => {
                // Push parameter when depth is 0
                params.push(parse_parametrized_type(&raw[start..i].to_vec()));
                start = i + 1;
            }
            _ => {}
        }
    }

    // Add the last parameter
    if start < raw.len() - 1 {
        params.push(parse_parametrized_type(&raw[start..raw.len() - 1].to_vec()));
    }

    ParamType { name, params }
}
