pub mod parser;

use parser::parse_parametrized_type;
use super::token::Token;

/// Stores data types that have parameters
/// For example: Result<str>
/// ------------------------------------------
/// This is done entirely due to requirements
/// of the static analyzer, to ensure we don't
/// transpile wrong code
/// ------------------------------------------
#[derive(Debug, Clone)]
pub struct ParamType {
    // From Result<str>, the structure would be:
    name: String, // Result
    params: Vec<ParamType>, // [str],
    is_reference: bool, // If the type is a reference (e.g. &str)
    raw_tokens: Vec<Token> // The raw tokens that make up this type
}

pub trait ParamTypeImpl {
    fn new(raw: &Vec<Token>) -> Self;
    fn get_name(&self) -> &String;
    fn get_params(&self) -> &Vec<ParamType>;
    fn get_raw_tokens(&self) -> &Vec<Token>;
    fn is_reference(&self) -> bool;
}

impl ParamTypeImpl for ParamType {
    fn new(raw: &Vec<Token>) -> Self {
        if raw.len() < 1 {
            return ParamType {
                name: String::from("[NATIVE BUILT-IN]"),
                params: Vec::new(),
                is_reference: false,
                raw_tokens: Vec::new(),
            }
        }

        parse_parametrized_type(raw)
    }

    fn get_name(&self) -> &String {
        &self.name
    }

    fn get_params(&self) -> &Vec<ParamType> {
        &self.params
    }

    fn get_raw_tokens(&self) -> &Vec<Token> {
        &self.raw_tokens
    }

    fn is_reference(&self) -> bool {
        self.is_reference
    }
    
}