use code::{token::{Token, TokenImpl}, token_type::TokenType};

use crate::sentence_extractor::extract_sentence;

pub fn extract_function_returns(function_body: &Vec<Token>) -> Vec<Vec<Token>> {
    function_body
        .iter()
        .enumerate()
        .filter(|(_, token)| token.get_token_type() == TokenType::Return)
        .map(|(i, _)| {
            // Extract the return value
            extract_sentence(
                function_body[i + 1..].to_vec(),
                TokenType::Semicolon
            )
        })
        .collect()
}