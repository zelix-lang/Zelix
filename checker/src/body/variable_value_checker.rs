use std::collections::HashMap;

use code::{token::{Token, TokenImpl}, token_type::TokenType};
use shared::code::param::{Param, ParamImpl};

use super::variable::{Variable, VariableImpl};

pub fn check_is_reference_to_param(
    scopes: &Vec<HashMap<String, Variable>>,
    parameters: &HashMap<String, Param>,
    value: &Vec<Token>
) -> bool {
    let value_token = &value[0];

    if value_token.get_token_type() != TokenType::Unknown {
        return false;
    }
    
    let value_name = value_token.get_value();

    if parameters.contains_key(&value_name) {
        return parameters.get(&value_name).unwrap().is_reference();
    }

    for scope in scopes{ 
        for (name, variable) in scope {
            if *name == value_name {
                // The variable should be a reference to a parameter
                // Otherwise it will outlive the function
                // which will cause a segmentation fault
                return variable.is_reference_to_param();
            }
        }
    }

    false
}