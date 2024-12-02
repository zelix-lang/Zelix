use std::{collections::HashMap, process::exit};

use shared::{logger::{Logger, LoggerImpl}, token::{token::TokenImpl, token_type::TokenType}};

use shared::code::function::{Function, FunctionImpl};

pub fn check_main_function(functions: &HashMap<String, Function>) {
    // Reached this point, the main function is always defined
    // So it's safe to get and unwrap it
    let main_function = functions.get("main").unwrap();
    if main_function.get_arguments().len() > 0 {
        Logger::err(
            "Main function cannot have arguments",
            &[
                "The main function cannot have arguments"
            ],
            &[
                main_function.get_trace().as_str()
            ]
        );

        exit(1);
    }

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
}