use std::{collections::HashMap, process::exit};

use code::{token::TokenImpl, token_type::TokenType};
use logger::{Logger, LoggerImpl};
use shared::code::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl};

use super::{variable::Variable, variable_checker::check_and_parse_variable};

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
        // we're going to have a Vec<HashMap<String, Variable>>
        // where each Vec represents a scope and each HashMap represents the variables in that scope
        let mut scope_variables: Vec<HashMap<String, Variable>> = Vec::new();
        // Push a new HashMap for the current scope (function body)
        // The variable checker should check for redefinition of imported classes and functions
        // so no need to re-check here
        scope_variables.push(HashMap::new());

        for (_, function) in file_functions {
            let body = function.get_body();

        for n in 0..body.len() {
            let token = &body[n];
            let token_type = token.get_token_type();

            match token_type {
                TokenType::Let => {
                    // n + 1 to skip the let token
                    let (parse_variable, var_name) =
                        check_and_parse_variable(
                            body,
                            n + 1,
                            &file_functions,
                            &headers
                        );

                    // Check if the variable is already defined
                    for scope in &scope_variables {
                        if scope.contains_key(&var_name.get_value()) {
                            throw_value_already_defined(
                                &var_name.get_value(),
                                &var_name.build_trace()
                            );
                        }
                    }

                    // Add the variable to the current scope
                    scope_variables.last_mut().unwrap().insert(
                        var_name.get_value(),
                        parse_variable
                    );
                }

                TokenType::OpenCurly => {
                    // Push a new HashMap for the new scope
                    scope_variables.push(HashMap::new());
                }

                TokenType::CloseCurly => {
                    // Pop the last scope
                    scope_variables.pop();
                }

                // No need to check in any other case
                _ => {}
            }
        }
        }
    }
}