use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};
use token_map::{KNOWN_TOKENS, NUMBER_REGEX, PUNCTUATION_CHARS};
mod import_processor;
pub mod data_types;
mod token_map;

pub struct Lexer {

    in_string: bool,
    in_escape: bool,
    in_comment: bool,
    in_block_comment: bool

}

pub trait LexerImpl {
    fn tokenize(
        &mut self,
        contents: &mut String,
        file: &String
    ) -> Vec<Token>;

    fn new() -> Self;
    fn calculate(current_token: &str, file: &String, line: &u32, col: &u32) -> Token;
}

impl LexerImpl for Lexer {
    fn tokenize(
        &mut self,
        contents: &mut String,
        file: &String
    ) -> Vec<Token> {
        let mut tokens: Vec<Token> = Vec::new();
        let mut current_token: String = String::new();
        let mut current_line: u32 = 1;
        let mut current_column: u32 = 1;
        let characters = contents.chars().collect::<Vec<char>>();
        let characters_len = characters.len();

        for i in 0 .. characters.len() {
            let tokens_len = tokens.len();
            let character = characters[i];

            if character == '\n' {
                if self.in_comment {
                    self.in_comment = false;
                    current_token.clear();
                }

                current_line += 1;
                current_column = 1;

                continue;
            } else {
                current_column += 1;
            }

            if self.in_comment {
                continue;
            }
            
            if !self.in_string && !self.in_block_comment && character == '/' {
                if characters_len < 2 {
                    continue;
                }

                if characters[i + 1] == '*' {
                    self.in_block_comment = true;
                } else if characters[i + 1] == '/' {
                    self.in_comment = true;
                }

                current_token.clear();

            } else if self.in_block_comment {
                if character == '/' && characters[i - 1] == '*' {
                    self.in_block_comment = false;
                }

                continue;
            } else if self.in_string {
                if character == '"' && !self.in_escape {
                    self.in_string = false;

                    tokens.push(
                        Token::new(
                            TokenType::StringLiteral,
                            current_token.clone(),
                            String::from(file),
                            current_line,
                            current_column
                        )
                    );

                    current_token.clear();
                } else if character == '\\' && !self.in_escape {
                    self.in_escape = true;
                    current_token.push(character);
                } else {
                    current_token.push(character);
                    self.in_escape = character == '\\';
                }

                continue;
            } else if character == '"' {
                self.in_string = true;
                continue;
            } else if character.is_whitespace() {
                if !current_token.is_empty() {
                    tokens.push(Lexer::calculate(
                        current_token.trim(),
                        &file,
                        &current_line,
                        &current_column
                    ));

                    current_token.clear();
                }
            } else if PUNCTUATION_CHARS.contains(&character) {
                if character == '.' && !current_token.parse::<i128>().is_err() {
                    current_token.push(character);
                    continue;
                }

                if !current_token.is_empty() {
                    tokens.push(Lexer::calculate(
                        &current_token,
                        &file,
                        &current_line,
                        &current_column
                    ));

                    current_token.clear();
                }

                if character == ']' {
                    if tokens_len < 2 {
                        continue;
                    }

                    let last_token = &tokens[tokens_len - 1].clone();
                    let last_last_token = &tokens[tokens_len - 2].clone();

                    if
                        (data_types::is_data_type(last_token.get_token_type()) &&
                        last_last_token.get_token_type() == TokenType::OpenBracket)
                        || (
                            last_last_token.get_token_type() == TokenType::OpenBracket
                            && last_token.get_value() == "discrete"
                        )
                        // num[], string[], bool[]
                    {
                        tokens.pop();
                        tokens.pop();

                        tokens.push(
                            Lexer::calculate(
                                &format!("{}[]", last_token.get_value().trim()),
                                &file,
                                &current_line,
                                &current_column
                            )
                        );

                        continue;
                    }
                } else if
                    character == '=' ||
                    character == '-' ||
                    character == '+'
                {
                    if tokens_len == 0 {
                        continue;
                    }

                    let last_token = &tokens[tokens_len - 1].clone();

                    if
                        (
                            character == '=' &&
                            (
                                last_token.get_token_type() == TokenType::Assign
                                || last_token.get_token_type() == TokenType::LessThan
                                || last_token.get_token_type() == TokenType::GreaterThan
                                || last_token.get_token_type() == TokenType::Not
                                || last_token.get_token_type() == TokenType::Plus
                                || last_token.get_token_type() == TokenType::Minus
                                || last_token.get_token_type() == TokenType::Asterisk
                                || last_token.get_token_type() == TokenType::Slash
                            )
                        )
                        || (
                            (character == '-' && last_token.get_token_type() == TokenType::Minus) ||
                            (character == '+' && last_token.get_token_type() == TokenType::Plus)
                        )
                    {
                        tokens.pop();
                        tokens.push(
                            Lexer::calculate(
                                &format!("{}{}", last_token.get_value().trim(), character),
                                &file,
                                &current_line,
                                &current_column
                            )
                        );
                        continue;
                    }
                } else if
                    character == '>'
                {
                    if tokens_len == 0 || tokens[tokens_len - 1].get_token_type() != TokenType::Minus {
                        continue;
                    }

                    tokens.pop();
                    tokens.push(
                        Lexer::calculate(
                            &String::from("->"),
                            &file,
                            &current_line,
                            &current_column
                        )
                    );

                    continue;
                }

                tokens.push(
                    Lexer::calculate(
                        &character.to_string(),
                        &file,
                        &current_line,
                        &current_column
                    )
                );
            } else {
                current_token.push(character);
            }
        }

        tokens
    }

    fn new() -> Self {
        Lexer {
            in_string: false,
            in_escape: false,
            in_comment: false,
            in_block_comment: false
        }
    }

    fn calculate(current_token: &str, file: &String, line: &u32, col: &u32) -> Token {
        let token_type = match KNOWN_TOKENS.get(current_token) {
            Some(token_type) => token_type,
            None => {
                if NUMBER_REGEX.is_match(current_token).unwrap() {
                    &TokenType::NumLiteral
                } else if current_token == "true" || current_token == "false" {
                    &TokenType::BoolLiteral
                } else {
                    &TokenType::Unknown
                }
            }
        };

        Token::new(
            token_type.clone(),
            current_token.to_string(),
            file.clone(),
            *line,
            *col
        )
    }

}