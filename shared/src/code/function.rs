use crate::token::token::Token;

use super::param::Param;

#[derive(Debug, Clone)]
pub struct Function {
    arguments: Vec<Param>,
    body: Vec<Token>,
    return_type: Token,
    trace: String,
    public: bool
}

pub trait FunctionImpl {

    fn new(
        arguments: 
        Vec<Param>, 
        body: Vec<Token>, 
        return_type: Token, 
        trace: String,
        public: bool
    ) -> Self;

    fn get_arguments(&self) -> &Vec<Param>;
    fn get_body(&self) -> &Vec<Token>;
    fn get_return_type(&self) -> &Token;
    fn get_trace(&self) -> &String;
    fn is_public(&self) -> bool;

}

impl FunctionImpl for Function {

    fn new(
        arguments: 
        Vec<Param>, 
        body: Vec<Token>, 
        return_type: Token, 
        trace: String,
        public: bool
    ) -> Self {
        Function {
            arguments,
            body,
            return_type,
            trace,
            public
        }
    }

    fn get_arguments(&self) -> &Vec<Param> {
        &self.arguments
    }

    fn get_body(&self) -> &Vec<Token> {
        &self.body
    }

    fn get_return_type(&self) -> &Token {
        &self.return_type
    }

    fn get_trace(&self) -> &String {
        &self.trace
    }

    fn is_public(&self) -> bool {
        self.public
    }
    
}