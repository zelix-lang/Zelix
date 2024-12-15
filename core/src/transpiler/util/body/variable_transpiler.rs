use c_parser::{class::ClassImpl, function::FunctionImpl, header::{Header, HeaderImpl}};
use code::{token::{Token, TokenImpl}, token_type::TokenType};
use lexer::data_types::is_data_type;

use extractor::{sentence_extractor::extract_sentence, token_splitter::extract_tokens_before};

use crate::transpiler::util::type_transpiler::transpile_type;

pub fn transpile_variable(
    tokens: &Vec<Token>,
    start_at: usize,
    transpiled_code: &mut String,
    imports: &Vec<Header>
) -> usize {
    let sentence: Vec<Token> = extract_sentence(
        // Also skip the let token
        tokens.clone()[(start_at + 1)..].to_vec(),
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
        // +1 for the name
        // +1 for the colon
        &sentence[2..].to_vec(),
        &TokenType::Assign
    );

    // +2 that we extracted before
    // +1 for the equals sign
    let var_value: &[Token] = &sentence[(var_type_tokens.len() + 3)..];

    transpile_type(&var_type_tokens, transpiled_code);

    transpiled_code.push_str(var_name);
    transpiled_code.push_str(" = ");

    // Get the first token to check if it needs a data type
    let first_token = &var_value[0];
    let mut first_processed = false; 

    if first_token.get_token_type() == TokenType::Unknown {
        let first_value = first_token.get_value();

        for import in imports {
            for (name, class) in import.get_classes() {
                if *name == first_value {
                    if class.get_generic_count() > 0 {
                        first_processed = true;
                        // This will result in code like:
                        // Result<String> my_var = new Result<String>(...);
                        // from source code:
                        // let my_var : Result<String> = Result(...);
                        transpiled_code.push_str(&first_value);
                        // Make sure the compiler knows the type
                        transpile_type(
                            &var_type_tokens[1..].to_vec(), 
                            transpiled_code
                        );
                    }

                    break;
                }
            }

            for (name, function) in import.get_functions() {
                if *name == first_value {
                    if function.get_generic_count() > 0 {
                        first_processed = true;
                        // This will result in code like:
                        // Result<String> my_var = new Result<String>(...);
                        // from source code:
                        // let my_var : Result<String> = Result(...);
                        transpiled_code.push_str(&first_value);
                        // Make sure the compiler knows the type
                        transpile_type(
                            &var_type_tokens[1..].to_vec(),
                            transpiled_code
                        );
                    }

                    break;
                }
            }
        }
    }

    for n in 0..var_value.len() {
        if first_processed && n == 0 {
            continue;
        }

        let token = &var_value[n];
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
    start_at + sentence.len() + 1
}