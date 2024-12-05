pub mod import_extractor;
pub mod sentence_extractor;
pub mod token_splitter;
pub mod return_extractor;
mod standard_locator;

use std::{path::PathBuf, process::exit};

use code::token::{Token, TokenImpl};
use code::token_type::TokenType;
use import_extractor::extract_import;
use c_parser::{create_c_instance, create_index};
use logger::{Logger, LoggerImpl};
use shared::code::import::{Import, Importable};

use shared::code::{file_code::{FileCode, FileCodeImpl}, function::{Function, FunctionImpl}, param::{Param, ParamImpl}};
use standard_locator::locate_standard;
use token_splitter::extract_tokens_before;

pub fn extract_parts(tokens: &Vec<Token>, source: PathBuf) -> FileCode {

    if tokens.is_empty() {
        Logger::err(
            "Refused to parse an empty file",
            &[
                "Write some code inside this file"
            ],
            &[
                format!(
                    "At {}",
                    source.to_str().unwrap().to_string()
                ).as_str()
            ]
        );

        exit(1);
    }

    let mut inside_function: bool = false;
    let mut result : FileCode = FileCode::new(source);

    let clang = create_c_instance();
    let index = create_index(&clang);

    // Add all the lang standard functions to the imports
    result.add_import(
        Import::new(
            locate_standard("lang/panic.h".to_string()),
            tokens[0].build_trace()
        ),
        &index
    );

    result.add_import(
        Import::new(
            locate_standard("lang/err.hpp".to_string()),
            tokens[0].build_trace()
        ),
        &index
    );

    result.add_import(
        Import::new(
        locate_standard("lang/result.hpp".to_string()),
            tokens[0].build_trace()
        ),
        &index
    );

    let mut expecting_function_name = false;
    let mut expecting_open_paren = false;
    let mut expecting_params = false;
    let mut expecting_param_type_splitter = false;
    let mut expecting_param_type = false;
    let mut expecting_open_curly = false;
    let mut has_function_ended = false;
    let mut expecting_arrow = false;
    let mut expecting_return_type = false;
    let mut is_last_function_public = false;
    let mut is_last_param_reference = false;

    // Used to count nested curly braces
    // This is useful because we could know when the function ends
    // if we encounter a closing curly brace and the nested_operations is 0
    let mut nested_operations = 0;

    let mut last_function_name = String::new();
    let mut last_function_return_type: Vec<Token> = Vec::new();
    let mut last_function_params: Vec<Param> = Vec::new();
    let mut last_function_body: Vec<Token> = Vec::new();
    let mut last_param_type_tokens: Vec<Token> = Vec::new();
    let mut last_param_name = String::new();

    // Used to skip tokens
    let mut skip_to_index = 0;

    // Used to determine nested levels in parameters
    let mut nested_level: isize = 0;

    for n in skip_to_index..tokens.len() {
        if skip_to_index > 0 && skip_to_index > n {
            continue;
        }

        let token = &tokens[n];
        let token_type : TokenType = token.get_token_type();

        if token_type == TokenType::Pub {
            is_last_function_public = true;

            if n + 1 >= tokens.len() {
                Logger::err(
                    "Invalid public declaration",
                    &["Expecting a function after the public keyword"],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            continue;
        } else if
            token_type == TokenType::Import
            && inside_function
        {
            Logger::err(
                "Invalid import",
                &["Import statement is not allowed inside a function"],
                &[token.build_trace().as_str()]
            );

            exit(1);
        } else if !inside_function {
            if token_type == TokenType::Import {
                let import = extract_import(
                    tokens.clone()[(n + 1)..].to_vec()
                );
                
                result.add_import(
                    import.clone(),
                    &index
                );
    
                // +1 because we skipped the import keyword
                // +1 for the semicolon
                // +1 for the string literal
                // total tokens skipped = 3
                skip_to_index = n + 3;
                continue;
            }

            if token_type == TokenType::Function {
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
    
                inside_function = true;
                has_function_ended = false;
                expecting_function_name = true;
                continue;
            }

            Logger::err(
                "Invalid token",
                &[
                    "You can't have code outside of a function"
                ],
                &[
                    "Unexpected token",
                    token.build_trace().as_str()
                ]
            );

            exit(1);

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
        } else if expecting_params {
            if token_type == TokenType::CloseParen {
                expecting_params = false;
                expecting_arrow = true;
                expecting_param_type_splitter = false;
                expecting_open_curly = false;
                
                continue;
            }

            if token_type != TokenType::Unknown {
                Logger::err(
                    "Invalid parameter name",
                    &[
                        "Parameter name must be a valid identifier",
                        format!(
                            "Got {:?}",
                            token_type
                        ).as_str()
                    ],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            last_param_name = token.get_value().to_string();
            expecting_params = false;
            expecting_param_type_splitter = true;
        } else if expecting_arrow {
            if token_type == TokenType::OpenCurly {
                // Implicitly add a void return type
                expecting_arrow = false;
                expecting_return_type = false;
                expecting_open_curly = false;
                nested_operations = 1;

                last_function_return_type = vec![
                    Token::new(
                        TokenType::Nothing,
                        "nothing".to_string(),
                        token.get_file(),
                        token.get_line(),
                        token.get_column()
                    )
                ];

                continue;
            } else if token_type != TokenType::Arrow {
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
            // Not going to validate the return type here
            // the static analyzer will do that

            let return_type = extract_tokens_before(
                &tokens[(n)..].to_vec(),
                &TokenType::OpenCurly
            );

            skip_to_index = n + return_type.len();
            last_function_return_type = return_type.clone();
            expecting_return_type = false;
            expecting_open_curly = true;
        } else if expecting_param_type_splitter {
            if token_type != TokenType::Colon {
                Logger::err(
                    "Invalid parameter declaration",
                    &[
                        "Expecting a colon after the parameter name",
                        format!(
                            "Got {:?}",
                            token_type
                        ).as_str()
                    ],
                    &[token.build_trace().as_str()]
                );

                exit(1);
            }

            expecting_param_type_splitter = false;
            expecting_param_type = true;
        } else if expecting_param_type {
            if token_type == TokenType::Ampersand {
                if is_last_param_reference {
                    Logger::err(
                        "Multiple reference",
                        &[
                            "You can't have multiple references for a single parameter"
                        ],
                        &[
                            token.build_trace().as_str()
                        ]
                    );

                    exit(1);
                }

                is_last_param_reference = true;
                // Wait for the parameter name
                continue;
            }

            if nested_level == 0 {
                if token_type == TokenType::Comma || token_type == TokenType::CloseParen {
                    if !last_param_type_tokens.is_empty() {
                        last_function_params.push(
                            Param::new(
                                last_param_name.clone(),
                                last_param_type_tokens.clone(),
                                token.build_trace(),
                                is_last_param_reference.clone()
                            )
                        );
                    }

                    last_param_name.clear();
                    last_param_type_tokens.clear();
                    is_last_param_reference = false;
                    expecting_param_type = false;

                    if token_type == TokenType::CloseParen {
                        expecting_params = false;
                        expecting_arrow = true;
                        expecting_open_curly = false;
                            
                        continue;
                    }

                    // Expect for another parameter
                    expecting_params = true;

                    continue;
                }
            }

            if token_type == TokenType::LessThan {
                nested_level += 1;
            } else if token_type == TokenType::GreaterThan {
                nested_level -= 1;

                if nested_level < 0 {
                    Logger::err(
                        "Invalid token",
                        &["Unexpected closing angle bracket"],
                        &[token.build_trace().as_str()]
                    );

                    exit(1);
                }
            }

            last_param_type_tokens.push(token.clone());

            // We don't check if the parameter is a data type here
            // that's done by the static analyzer
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
                        token.get_file(),
                        Function::new(
                            last_function_params.clone(),
                            last_function_body.clone(),
                            last_function_return_type.clone(),
                            token.build_trace(),
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
                    expecting_open_curly = false;
                    expecting_arrow = false;
                    expecting_return_type = false;
                    is_last_function_public = false;
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
            &[
                "Function must end with a closing curly brace",
                tokens[tokens.len() - 1].build_trace().as_str()
            ]
        );

        exit(1);
    }

    result

}