use shared::token::{token::Token, token_type::TokenType};

use super::param::Param;

#[derive(Debug, Clone)]
pub struct Function {
    arguments: Vec<Param>,
    body: Vec<Token>,
    return_type: TokenType,
    trace: Token,
    public: bool
}

pub trait FunctionImpl {

    fn new(
        arguments: 
        Vec<Param>, 
        body: Vec<Token>, 
        return_type: TokenType, 
        trace: Token,
        public: bool
    ) -> Self;

    fn get_arguments(&self) -> &Vec<Param>;
    fn get_body(&self) -> &Vec<Token>;
    fn get_return_type(&self) -> &TokenType;
    fn get_trace(&self) -> &Token;
    fn is_public(&self) -> bool;

}

impl FunctionImpl for Function {

    fn new(
        arguments: 
        Vec<Param>, 
        body: Vec<Token>, 
        return_type: TokenType, 
        trace: Token,
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

    fn get_return_type(&self) -> &TokenType {
        &self.return_type
    }

    fn get_trace(&self) -> &Token {
        &self.trace
    }

    fn is_public(&self) -> bool {
        self.public
    }
    
}