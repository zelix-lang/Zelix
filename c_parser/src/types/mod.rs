pub mod parser;
use parser::parse_parametrized_type;
use shared::token::token::Token;

/// Stores data types that have parameters
/// For example: Result<str>
/// ------------------------------------------
/// This is done entirely due to requirements
/// of the static analyzer, to ensure we don't
/// transpile wrong code
/// ------------------------------------------
#[derive(Debug)]
pub struct ParamType {
    // From Result<str>, the structure would be:
    name: String, // Result
    params: Vec<ParamType>, // [str],
    raw_tokens: Vec<Token> // The raw tokens that make up this type
}

pub trait ParamTypeImpl {
    fn new(raw: &Vec<Token>) -> Self;
    fn get_name(&self) -> &String;
    fn get_params(&self) -> &Vec<ParamType>;
    fn get_raw_tokens(&self) -> &Vec<Token>;
}

impl ParamTypeImpl for ParamType {
    fn new(raw: &Vec<Token>) -> Self {
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
}