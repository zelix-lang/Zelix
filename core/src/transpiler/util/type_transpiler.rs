use code::{token::{Token, TokenImpl}, token_type::TokenType};

pub fn transpile_type(
    tokens: &Vec<Token>,
    transpiled_code: &mut String
) {

    for token in tokens {
        let token_type = token.get_token_type();
        
        if token_type == TokenType::Nothing {
            transpiled_code.push_str("void ");
        } else if token_type == TokenType::Num {
            transpiled_code.push_str("double ");
        } else if token_type == TokenType::String {
            transpiled_code.push_str("std::string ");
        } else if token_type == TokenType::Bool {
            transpiled_code.push_str("bool ");
        } else {
            // Static analyzer should catch errors in case this is undefined
            transpiled_code.push_str(token.get_value().as_str());
            transpiled_code.push_str(" ");
        }
    }

}