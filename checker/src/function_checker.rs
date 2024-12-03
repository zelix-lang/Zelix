use std::{collections::HashMap, path::PathBuf, process::exit};

use logger::{Logger, LoggerImpl};
use shared::{path::discard_cwd, token::{token::TokenImpl, token_type::TokenType}};

use shared::code::{function::{Function, FunctionImpl}, value_name::value_name::VALUE_NAME_REGEX};

pub fn analyze_functions(
    // Pass by reference to avoid moving the value or cloning it 
    functions: &HashMap<String, Function>,
    source: &PathBuf
) {

    // First check if the main function is even defined
    if !functions.contains_key("main") {
        Logger::err(
            "Main function not defined",
            &[
                "The main function must be defined in the code"
            ],
            &[
                format!(
                    "At {}",
                    discard_cwd(source.to_str().unwrap().to_string())
                ).as_str()
            ]
        );

        exit(1);
    }

    // The return type of the main function must be Nothing
    let main_function = functions.get("main").unwrap();
    if main_function.get_return_type().get(0).unwrap().get_token_type() != TokenType::Nothing {
        Logger::err(
            "Invalid return type for main function",
            &[
                "The return type of the main function must be Nothing"
            ],
            &[
                main_function.get_trace().as_str()
            ]
        );

        exit(1);
    }

    for (name, function) in functions.iter() {
        if !VALUE_NAME_REGEX.is_match(name.as_str()).unwrap_or(false) {
            Logger::err(
                format!("Invalid function name: {}", name).as_str(),
                &[
                    "Function names must start with a letter or an underscore",
                ],
                &[
                    function.get_trace().as_str()
                ]
            );

            exit(1);
        }

        // The lexer should have already checked if the function name is repeated
        // so we don't need to check it here
    }

}