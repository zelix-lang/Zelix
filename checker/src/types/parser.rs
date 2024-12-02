use extractor::token_splitter::split_tokens;
use shared::{logger::{Logger, LoggerImpl}, token::{token::{Token, TokenImpl}, token_type::TokenType}};

use super::ParamType;

fn assert_token_type(raw: &Vec<Token>, index: usize, check_for: &TokenType) {
    let token = &raw[index];
    let token_type = token.get_token_type();

    if token_type != *check_for {
        Logger::err(
            "Invalid parametrized type",
            &[
                "This token was unexpected at this time",
                "Define parametrized types like: Result<str>"
            ],
            &[
                token.build_trace().as_str()
            ]
        );
    }
}

pub fn parse_parametrized_type(raw: &Vec<Token>) -> ParamType {
    let mut name = String::new();
    let mut params = Vec::new();

    // In case the tokens only have 1 element
    // it's most likely a type without parameters
    if raw.len() == 1 {
        return ParamType {
            name: raw[0].get_value(),
            params
        };
    }

    // The tokens should have at least 4 elements
    if raw.len() < 4 {
        Logger::err(
            "Invalid parametrized type",
            &[
                "Use parameters in types like: Result<str>"
            ],
            &[
                raw[0].build_trace().as_str()
            ]
        );
    }

    // Validate the token types
    assert_token_type(raw, 0, &TokenType::Unknown);
    assert_token_type(raw, 1, &TokenType::LessThan);
    assert_token_type(raw, raw.len() - 1, &TokenType::GreaterThan);

    name = raw[0].get_value();

    // Split by commas, accounting for nested types
    let mut splitted = split_tokens(
        &raw[2..raw.len() - 1].to_vec(), // Skip the "Result" and enclosing "<" ">"
        &TokenType::Comma,
        &TokenType::LessThan,
        &TokenType::GreaterThan
    );

    // Process each sublist in a loop (queue-based approach)
    while !splitted.is_empty() {
        let first_element = splitted.remove(0);

        if first_element.len() == 1 {
            // Single-token parameter
            params.push(ParamType {
                name: first_element[0].get_value(),
                params: Vec::new(),
            });
        } else {
            // Process nested parameterized types iteratively
            let mut nested_stack = vec![first_element];
            while let Some(current_nested) = nested_stack.pop() {
                if current_nested.len() == 1 {
                    // Nested single token
                    params.push(ParamType {
                        name: current_nested[0].get_value(),
                        params: Vec::new(),
                    });
                } else {
                    // Parse the nested parameterized type
                    let nested_name = current_nested[0].get_value();

                    // Skip "<" and ">"
                    let nested_body = &current_nested[2..current_nested.len() - 1].to_vec();

                    let nested_splitted = split_tokens(
                        nested_body,
                        &TokenType::Comma,
                        &TokenType::LessThan,
                        &TokenType::GreaterThan,
                    );

                    for part in nested_splitted {
                        nested_stack.push(part);
                    }

                    params.push(ParamType {
                        name: nested_name,
                        params: Vec::new(),
                    });
                }
            }
        }
    }

    ParamType { name, params }
}
