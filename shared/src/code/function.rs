use std::collections::HashMap;

use code::{token::Token, types::{parser::parse_parametrized_type, ParamType}};

use super::param::Param;

#[derive(Debug, Clone)]
pub struct Function {
    arguments: HashMap<String, Param>,
    body: Vec<Token>,
    return_type: ParamType,
    trace: String,
    public: bool
}

pub trait FunctionImpl {

    fn new(
        arguments: HashMap<String, Param>, 
        body: Vec<Token>, 
        return_type: Vec<Token>, 
        trace: String,
        public: bool
    ) -> Self;

    fn get_arguments(&self) -> &HashMap<String, Param>;
    fn get_body(&self) -> &Vec<Token>;
    fn get_return_type(&self) -> &ParamType;
    fn get_trace(&self) -> &String;
    fn is_public(&self) -> bool;

}

impl FunctionImpl for Function {

    fn new(
        arguments: HashMap<String, Param>, 
        body: Vec<Token>, 
        return_type: Vec<Token>, 
        trace: String,
        public: bool
    ) -> Self {
        Function {
            arguments,
            body,
            return_type: parse_parametrized_type(&return_type),
            trace,
            public
        }
    }

    fn get_arguments(&self) -> &HashMap<String, Param> {
        &self.arguments
    }

    fn get_body(&self) -> &Vec<Token> {
        &self.body
    }

    fn get_return_type(&self) -> &ParamType {
        &self.return_type
    }

    fn get_trace(&self) -> &String {
        &self.trace
    }

    fn is_public(&self) -> bool {
        self.public
    }
    
}