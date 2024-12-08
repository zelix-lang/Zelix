use std::{collections::HashMap, process::exit};

use c_parser::{class::ClassImpl, header::Header};
use code::{token::{Token, TokenImpl}, token_type::TokenType, types::{parser::parse_parametrized_type, ParamTypeImpl}};
use extractor::token_splitter::extract_tokens_before;
use fancy_regex::Regex;
use lazy_static::lazy_static;
use lexer::data_types::is_data_type;
use shared::code::{function::Function, param::Param, value_name::value_name::{CPP_KEYWORDS, VALUE_NAME_REGEX}};
use logger::{Logger, LoggerImpl};
use util::result::try_unwrap;

use crate::header::header_checker::{check_header_value_definition, find_imported_classes};

use super::{scope_checker::throw_value_already_defined, variable::{Variable, VariableImpl}, variable_value_checker::check_is_reference_to_param};

lazy_static! {
    // Used to print warnings for cammel case variable names
    // Surf encourages snake case variable names!
    pub static ref CAMMEL_CASE_REGEX: Regex = 
        Regex::new(r"^[a-zA-Z][A-Z0-9]{1,}$").unwrap();
}

fn check_variable_name(var_name: &String, trace: &String) {
    if try_unwrap(
        CAMMEL_CASE_REGEX.is_match(var_name),
        "Failed to validate a variable name"
    ) {
        Logger::warn(
            "Consider using snake case for variable names",
            &[
                format!(
                    "Consider converting {} to snake case",
                    var_name
                ).as_str(),
                trace.as_str()
            ],
        );
    }
}

// Returns: Variable type, Variable name, Variable value
pub fn check_and_parse_variable(
    tokens: &Vec<Token>,
    start: usize,
    // Used to check if a value is already defined
    functions: &HashMap<String, Function>,
    imports: &Vec<Header>,
    scopes: &Vec<HashMap<String, Variable>>,
    parameters: &HashMap<String, Param>
) -> (Variable, Token) {
    // Variable definitions should be already validated by now
    // Example definition:
    // let my_var : str = "Hello, world!";
    // Number of tokens: 7

    // The first token should be the variable name (We don't receive the let token)
    let variable_tokens = &extract_tokens_before(&tokens[start..].to_vec(), &TokenType::Semicolon);
    if variable_tokens.len() < 4 {
        Logger::err(
            "Invalid variable definition",
            &[
                "Variable definitions must have at least 4 tokens"
            ],
            &[
                tokens[start].build_trace().as_str()
            ]
        );

        exit(1);
    }

    let var_name = &variable_tokens[0];
    let var_name_value = var_name.get_value();
    let colon = &variable_tokens[1];
    let value_tokens = variable_tokens[3..].to_vec();

    let var_type_tokens = extract_tokens_before(
        // + 1 for the name
        // + 1 for the colon
        &variable_tokens[2..].to_vec(),
        &TokenType::Assign
    );

    let var_type_param_type = parse_parametrized_type(
        &var_type_tokens
    );

    let var_type_tokens_len = var_type_tokens.len();

    // The same 2 we extracted before
    let equals = &variable_tokens[var_type_tokens_len + 2];

    check_variable_name(&var_name.get_value(), &var_name.build_trace());

    // The var name should be an unknown token (not a keyword)
    if
        var_name.get_token_type() != TokenType::Unknown
        || !try_unwrap(
            VALUE_NAME_REGEX.is_match(var_name.get_value().as_str()),
            "Failed to validate a variable name"
        )
        || CPP_KEYWORDS.contains(&var_name.get_value().as_str())
    {
        Logger::err(
            "Invalid variable name",
            &[
                "Variable names must be unknown tokens"
            ],
            &[
                var_name.build_trace().as_str()
            ]
        );

        exit(1);
    }

    // The colon should be a colon
    if colon.get_token_type() != TokenType::Colon {
        Logger::err(
            "Invalid variable definition",
            &[
                "Expected a colon after the variable name"
            ],
            &[
                colon.build_trace().as_str()
            ]
        );

        exit(1);
    }

    // The equals should be an equals
    if equals.get_token_type() != TokenType::Assign {
        println!("{:?}", equals);

        Logger::err(
            "Invalid variable definition",
            &[
                "Expected an equals sign after the variable type"
            ],
            &[
                equals.build_trace().as_str()
            ]
        );

        exit(1);
    }

    // Check if the var name is already defined
    if
        functions.contains_key(var_name_value.as_str()) ||
        check_header_value_definition(&var_name_value, imports)
    {
        throw_value_already_defined(
            &var_name_value,
            &var_name.build_trace()
        );
    }

    let generic_params = var_type_param_type.get_params();

    let raw_tokens = var_type_param_type.get_raw_tokens();
    let var_type = &raw_tokens[0];

    if !is_data_type(var_type.get_token_type()) {
        // Check if the variable name is in cammel case
        let class_optional = find_imported_classes(
            &var_type.get_value(),
             imports
        );
    
        if class_optional.is_none() {
            Logger::err(
                "Invalid data type",
                &[
                    "The data type is not recognized"
                ],
                &[
                    var_type.build_trace().as_str()
                ]
            );
    
            exit(1);
        }

        let generic_count = class_optional.unwrap().get_generic_count();
        let provided_count = generic_params.len();

        if generic_count != provided_count {
            Logger::err(
                "Mismatched generic parameters",
                &[
                    format!(
                        "The data type takes {} parameters, but {} were provided",
                        generic_count,
                        provided_count
                    ).as_str()
                ],
                &[
                    var_type.build_trace().as_str()
                ]
            );
    
            exit(1);
        }

    }

    // If the variable itself isn't a reference
    // whatever value it has won't be a reference if returned
    // so before checking lifetime, we'll check if the variable
    // is a reference
    let value_tokens_len = value_tokens.len();

    if value_tokens_len == 0 {
        Logger::err(
            "Invalid variable definition",
            &[
                "Expected a value after the equals sign"
            ],
            &[
                equals.build_trace().as_str()
            ]
        );

        exit(1);
    }

    let is_reference_to_param = check_is_reference_to_param(
        scopes, 
        parameters,
        &value_tokens
    );

    if
        !is_reference_to_param && var_type_param_type.is_reference()
        || is_reference_to_param && !var_type_param_type.is_reference()
    {
        Logger::err(
            "Invalid reference",
            &[
                "The variable is a reference, but the value is not"
            ],
            &[
                equals.build_trace().as_str()
            ]
        );

        exit(1);
    }

    // TODO: Validate that the variable's value matches the type

    let parsed_variable = Variable::new(
        is_reference_to_param
    );

    (
        parsed_variable,
        var_name.clone()
    )
    
}