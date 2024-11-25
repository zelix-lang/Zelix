use std::process::exit;

use shared::{logger::{Logger, LoggerImpl}, token::{token::TokenImpl, token_type::TokenType}};

use crate::shared::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl};

fn throw_value_already_defined(name: &String, trace: &String) {
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
pub fn analyze_scope(source: FileCode) {
    let functions = source.get_functions();
    
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

    if main_function.get_return_type().get_token_type() != TokenType::Nothing {
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