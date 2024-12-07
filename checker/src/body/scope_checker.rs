use std::process::exit;

use code::{token::TokenImpl, token_type::TokenType};
use logger::{Logger, LoggerImpl};
use shared::code::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl};

use super::variable_checker::check_variable;

pub fn throw_value_already_defined(name: &String, trace: &String) {
    Logger::err(
        "Value already defined",
        &[
            "Choose a different name for the value"
        ],
        &[
            trace.as_str(),
            format!("The value {} is already defined", name).as_str()
        ],
    );

    exit(1);
}

// Analyzes the source code to determine undefined variables
pub fn analyze_scope(source: &FileCode) {
    let functions = source.get_functions();
    let headers = source.get_imports();

    for (_, file_functions) in functions {
        // Define a vector that contains the variables allowed in the current scope
        // since there can be multiple scopes in a function
        // we're going to have a Vec<
        for (_, function) in file_functions {
            let body = function.get_body();

        for n in 0..body.len() {
            let token = &body[n];
            let token_type = token.get_token_type();

            match token_type {
                TokenType::Let => {
                    // n + 1 to skip the let token
                    check_variable(
                        body,
                        n + 1,
                        &file_functions,
                        &headers
                    );
                }

                // No need to check in any other case
                _ => {}
            }
        }
        }
    }
}