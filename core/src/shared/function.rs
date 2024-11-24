use shared::token::{token::Token, token_type::TokenType};

use super::param::Param;

#[derive(Debug, Clone)]
pub struct Function {
    arguments: Vec<Param>,
    body: Vec<Token>,
    return_type: TokenType
}

pub trait FunctionImpl {

    fn new(arguments: Vec<Param>, body: Vec<Token>, return_type: TokenType) -> Self;

    fn get_arguments(&self) -> &Vec<Param>;
    fn get_body(&self) -> &Vec<Token>;
    fn get_return_type(&self) -> &TokenType;

}

impl FunctionImpl for Function {

    fn new(arguments: Vec<Param>, body: Vec<Token>, return_type: TokenType) -> Self {
        Function {
            arguments,
            body,
            return_type
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
    
}