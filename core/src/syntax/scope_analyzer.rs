use std::process::exit;

use shared::{logger::{Logger, LoggerImpl}, token::token_type::TokenType};

use crate::shared::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl, import::Importable, param::ParamImpl};

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
    let imports = source.get_imports()
        .iter()
        .map(|import| import.get_name())
        .collect::<Vec<String>>();

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

    if main_function.get_return_type() != &TokenType::Nothing {
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

    for (name, function) in functions {
        // Check the arguments' names to see if any of them are
        // already defined through an import
        for arg in function.get_arguments() {
            if imports.contains(arg.get_name()) {
                throw_value_already_defined(arg.get_name(), &arg.get_trace());
            }
        }

        // Check if the function name is already defined through an import
        if imports.contains(name) {
            throw_value_already_defined(name, function.get_trace());
        }

    }
}