use code::{token::{Token, TokenImpl}, token_type::TokenType};
use lexer::data_types::is_data_type;

use extractor::{sentence_extractor::extract_sentence, token_splitter::extract_tokens_before};

use crate::transpiler::util::type_transpiler::transpile_type;

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

    let var_name = &sentence[0].get_value();
    let var_type_tokens = extract_tokens_before(
        &tokens[3..].to_vec(),
        &TokenType::Assign
    );

    // +2 for the space between the parametrized type
    // +1 for the equals sign
    let var_value: &[Token] = &sentence[(var_type_tokens.len() + 3)..];

    // Make the compiler automatically infer the type
    transpiled_code.push_str("auto ");

    transpiled_code.push_str(var_name);
    transpiled_code.push_str(" = ");

    for token in var_value {
        let token_type = token.get_token_type();
        
        if is_data_type(token_type.clone()) {
            transpile_type(&vec![token.clone()], transpiled_code);
            continue;
        }

        let is_string = token_type == TokenType::StringLiteral;

        if is_string {
            transpiled_code.push_str("\"");
        }

        transpiled_code.push_str(&token.get_value());

        if is_string {
            transpiled_code.push_str("\"");
        }

    }

    // Add 1 to skip the let token, we still need the semicolon
    n + sentence.len() + 1
}