use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};

use super::variable_transpiler::transpile_variable;

pub fn transpile_body(tokens: &Vec<Token>, transpiled_code: &mut String) {
    // Used to skip tokens
    let mut skip_to_index = 0;

    for n in 0..tokens.len() {
        if n < skip_to_index {
            continue;
        }
        
        let token = &tokens[n];
        let token_type = token.get_token_type();
        let is_string = token_type == TokenType::StringLiteral;

        if is_string {
            transpiled_code.push_str("\"");
        } else if token_type == TokenType::Let {

            // Add 2 to skip the let token and the semicolon
            skip_to_index = transpile_variable(tokens, n, transpiled_code);
            continue;
        }

        transpiled_code.push_str(&token.get_value());

        if is_string {
            transpiled_code.push_str("\"");
        }
    }
}