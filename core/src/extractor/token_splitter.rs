use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};

pub fn split_tokens(
    tokens: &Vec<Token>,
    splitter: TokenType
) -> Vec<Vec<Token>> {
    let mut result: Vec<Vec<Token>> = Vec::new();
    let mut current: Vec<Token> = Vec::new();

    for token in tokens {
        if token.get_token_type() == splitter {
            if current.len() < 1 {
                continue;
            }

            result.push(current);
            current = Vec::new();
        } else {
            current.push(token.clone());
        }
    }

    if current.len() > 0 {
        result.push(current);
    }

    result
}