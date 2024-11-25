use crate::path::discard_cwd;

use super::token_type::TokenType;

#[derive(Debug, Clone, PartialEq)]
pub struct Token {
    token_type: TokenType,
    value: String,
    file: String,
    line: u32,
    column: u32,
}

pub trait TokenImpl {
    fn new(token_type: TokenType, value: String, file: String, line: u32, column: u32) -> Token;
    fn get_token_type(&self) -> TokenType;
    fn get_value(&self) -> String;
    fn get_file(&self) -> String;
    fn get_line(&self) -> u32;
    fn get_column(&self) -> u32;
    fn build_trace(&self) -> String;
}

impl TokenImpl for Token {
    fn new(token_type: TokenType, value: String, file: String, line: u32, column: u32) -> Token {
        Token {
            token_type,
            value,
            file,
            line,
            column,
        }
    }

    fn get_token_type(&self) -> TokenType {
        self.token_type.clone()
    }

    fn get_value(&self) -> String {
        self.value.clone()
    }

    fn get_file(&self) -> String {
        self.file.clone()
    }

    fn get_line(&self) -> u32 {
        self.line
    }

    fn get_column(&self) -> u32 {
        self.column
    }

    fn build_trace(&self) -> String {
        format!(
            "At {}:{}:{}",
            discard_cwd(self.get_file()),
            self.get_line(),
            self.get_column()
        )
    }
}