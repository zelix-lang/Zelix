use std::collections::HashMap;

use c_parser::header::Header;
use code::token::{Token, TokenImpl};
use shared::code::{function::Function, param::Param};

use super::variable::Variable;

pub fn check_is_reference_to_param(
    imports: &Vec<Header>,
    scopes: &Vec<HashMap<String, Variable>>,
    parameters: &HashMap<String, Param>,
    functions: &HashMap<String, Function>,    
    value: &Vec<Token>
) -> bool {
    let value_len = value.len();

    if value_len == 1 {
        let value_token = &value[0];
        let value_name = value_token.get_value();
    }

    true
}