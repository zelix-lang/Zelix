use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};

use crate::extractor::sentence_extractor::extract_sentence;

use super::type_transpiler::transpile_type;

pub fn transpile_variable(tokens: &Vec<Token>, n: usize, transpiled_code: &mut String) -> usize {
    let sentence: Vec<Token> = extract_sentence(
        // Also skip the let token
        tokens.clone()[(n + 1)..].to_vec(),
        TokenType::Semicolon
    );

    // Variable definitions should be already validated by now
    // Example definition:
    // let my_var : str = "Hello, world!";
    // Number of tokens: 7
    // We don't include the let token and the semicolon
    // so we're left with 5 tokens
    // we're going to skip the first 4 so we get the value

    let var_type = &sentence[2];
    let var_name = &sentence[0].get_value();
    let var_value: &[Token] = &sentence[4..];

    transpile_type(var_type, transpiled_code);
    transpiled_code.push_str(var_name);
    transpiled_code.push_str(" = ");
    
    for token in var_value {
        transpiled_code.push_str(&token.get_value());
    }

    // Add 1 to skip the let token, we still need the semicolon
    n + sentence.len() + 1
}