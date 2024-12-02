mod parser;
use parser::parse_parametrized_type;
use shared::token::token::Token;

/// Stores data types that have parameters
/// For example: Result<str>
/// ------------------------------------------
/// This is done entirely due to requirements
/// of the static analyzer, to ensure we don't
/// transpile wrong code
/// ------------------------------------------
pub struct ParamType {
    // From Result<str>, the structure would be:
    pub name: String, // Result
    pub params: Vec<ParamType> // [str]
}

trait ParamTypeImpl {
    fn new(raw: &Vec<Token>) -> Self;
    fn get_name(&self) -> &String;
    fn get_params(&self) -> &Vec<ParamType>;
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
}