mod import_extractor;
mod sentence_extractor;
mod standard_locator;
mod token_splitter;
use std::{path::PathBuf, process::exit};

use import_extractor::extract_import;
use lexer::data_types::is_data_type;
use shared::{logger::{Logger, LoggerImpl}, token::{token::{Token, TokenImpl}, token_type::TokenType}};

use crate::shared::{file_code::{FileCode, FileCodeImpl}, function::{Function, FunctionImpl}, import::Importable, import_group::ImportGroupImpl, param::{Param, ParamImpl}};

pub fn extract_parts(tokens: &Vec<Token>, source: PathBuf) -> FileCode {

    let mut inside_function: bool = false;
    let mut result : FileCode = FileCode::new(source);

    let mut expecting_function_name = false;
    let mut expecting_open_paren = false;
    let mut expecting_params = false;
    let mut expecting_param_type_splitter = false;
    let mut expecting_param_type = false;
    let mut expecting_comma = false;
    let mut expecting_open_curly = false;
    let mut has_function_ended = false;
    let mut expecting_arrow = false;
    let mut expecting_return_type = false;
    let mut expecting_fun_keyword = true;
    let mut is_last_function_public = false;


    // Used to count nested curly braces
    // This is useful because we could know when the function ends
    // if we encounter a closing curly brace and the nested_operations is 0
    let mut nested_operations = 0;

    let mut last_function_name = String::new();
    let mut last_function_return_type: TokenType = TokenType::Unknown;
    let mut last_function_params: Vec<Param> = Vec::new();
    let mut last_function_body: Vec<Token> = Vec::new();
    let mut last_param_name = String::new();

    // Used to skip tokens
    let mut skip_to_index = 0;

    for n in skip_to_index..tokens.len() {
        if skip_to_index > 0 && skip_to_index > n {
            continue;
        }

        let token = &tokens[n];
        let token_type : TokenType = token.get_token_type();

        if token_type == TokenType::Pub {
            is_last_function_public = true;
            expecting_fun_keyword = true;

            continue;
        } else if token_type == TokenType::Import {
            if inside_function {
                Logger::err(
                    "Invalid import",
                    &["Import statement is not allowed inside a function"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            let imports = extract_import(
                tokens.clone()[(n + 1)..].to_vec()
            );

            for import in imports.get_imports() {
                result.add_import(import.clone());
            }

            // +1 because we skipped the import keyword
            skip_to_index = n + imports.get_skipped_tokens() as usize + 1;

        } else if token_type == TokenType::Function {
            if !expecting_fun_keyword {
                Logger::err(
                    "Invalid function",
                    &["Expecting a function keyword"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            if inside_function {
                Logger::err(
                    "Invalid function",
                    &[
                        "Define this function outside of the current function",
                        "Use modules to organize your code"
                    ],
                    &[
                        token.build_trace().as_str(),
                        "You can't define a function inside another function"
                    ]
                );
                
                exit(1);
            }

            expecting_fun_keyword = false;
            inside_function = true;
            has_function_ended = false;
            expecting_function_name = true;
        } else if expecting_function_name {
            // Expecting an unknown token here (not a keyword)
            // Name is going to be validated later by syntax checker

            if token_type != TokenType::Unknown {
                Logger::err(
                    "Invalid function name",
                    &["Function name must be a valid identifier"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_function_name = false;
            last_function_name = token.get_value().to_string();
            expecting_open_paren = true;

        } else if expecting_open_paren {
            if token_type != TokenType::OpenParen {
                Logger::err(
                    "Invalid function declaration",
                    &["Expecting an open parenthesis after the function name"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_open_paren = false;
            expecting_params = true;
        } else if expecting_params || expecting_comma {
            if token_type == TokenType::CloseParen {
                expecting_params = false;
                expecting_arrow = true;
                expecting_open_curly = false;
                
                continue;
            }

            if expecting_comma && token_type != TokenType::Comma {
                Logger::err(
                    "Invalid parameter declaration",
                    &["Expecting a comma after the parameter type"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            if token_type != TokenType::Unknown {
                Logger::err(
                    "Invalid parameter name",
                    &["Parameter name must be a valid identifier"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            last_param_name = token.get_value().to_string();
            expecting_params = false;
            expecting_param_type_splitter = true;
        } else if expecting_arrow {
            if token_type != TokenType::Arrow {
                Logger::err(
                    "Invalid function declaration",
                    &["Expecting an arrow after the function parameters"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_arrow = false;
            expecting_return_type = true;
        } else if expecting_return_type {
            if !is_data_type(token_type.clone()) {
                Logger::err(
                    "Invalid return type",
                    &["Expecting a valid data type"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            last_function_return_type = token_type.clone();
            expecting_return_type = false;
            expecting_open_curly = true;
        } else if expecting_param_type_splitter {
            if token_type != TokenType::Colon {
                Logger::err(
                    "Invalid parameter declaration",
                    &["Expecting a colon after the parameter name"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_param_type_splitter = false;
            expecting_param_type = true;
        } else if expecting_param_type {
            if !is_data_type(token_type.clone()) {
                Logger::err(
                    "Invalid parameter type",
                    &["Expecting a valid data type"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            last_function_params.push(
                Param::new(
                    last_param_name.clone(),
                    token_type.clone()
                )
            );

            last_param_name.clear();
            expecting_param_type = false;
            expecting_comma = true;
        } else if expecting_open_curly {
            if token_type != TokenType::OpenCurly {
                Logger::err(
                    "Invalid function declaration",
                    &["Expecting an open curly brace after the function parameters"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_open_curly = false;
            nested_operations = 1;
            has_function_ended = false;
        } else {
            if token_type == TokenType::OpenCurly {
                nested_operations += 1;
            } else if token_type == TokenType::CloseCurly {
                nested_operations -= 1;

                if nested_operations < 0 {
                    Logger::err(
                        "Invalid token",
                        &["Unexpected closing curly brace"],
                        &[token.build_trace().as_str()]
                    );

                    exit(1);
                }

                if nested_operations == 0 {
                    has_function_ended = true;

                    result.add_function(
                        last_function_name.clone(),
                        Function::new(
                            last_function_params.clone(),
                            last_function_body.clone(),
                            last_function_return_type.clone(),
                            token.clone(),
                            is_last_function_public
                        )
                    );


                    // Reset all flags
                    inside_function = false;
                    expecting_function_name = false;
                    expecting_open_paren = false;
                    expecting_params = false;
                    expecting_param_type_splitter = false;
                    expecting_param_type = false;
                    expecting_comma = false;
                    expecting_open_curly = false;
                    expecting_arrow = false;
                    expecting_return_type = false;
                    is_last_function_public = false;
                    last_function_return_type = TokenType::Unknown;
                    last_function_params.clear();
                    last_function_body.clear();
                    last_function_name.clear();
                    continue;
                }
            }

            last_function_body.push(token.clone());
        }
    }

    if !has_function_ended {
        Logger::err(
            "Invalid function declaration",
            &["Expecting a closing curly brace"],
            &["Function must end with a closing curly brace"]
        );

        exit(1);
    }

    result

}