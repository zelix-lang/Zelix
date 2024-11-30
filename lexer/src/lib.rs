use std::path::PathBuf;

use import::import_processor::process_imports;
use regex_patterns::NUMBER_REGEX;
use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};
use token_map::{KNOWN_TOKENS, PUNCTUATION_CHARS};
mod regex_patterns;
mod token_map;
mod import;
pub mod data_types;

// Define the Lexer struct, tracking the parser state (string, escape, comments)
pub struct Lexer {
    in_string: bool,
    in_escape: bool,
    in_comment: bool,
    in_block_comment: bool
}

// Trait defining the behavior of a Lexer
pub trait LexerImpl {
    // Method to tokenize the input contents
    fn tokenize(
        &mut self,
        contents: &mut String,
        file: &PathBuf
    ) -> Vec<Token>;

    // Factory method to create a new Lexer instance
    fn new() -> Self;

    // Helper method to calculate and return a Token from the current context
    fn calculate(current_token: &str, file: &String, line: &u32, col: &u32) -> Token;
}

// Implementation of the LexerImpl trait for the Lexer struct
impl LexerImpl for Lexer {
    // Tokenize the given contents into a list of Tokens
    fn tokenize(
        &mut self,
        contents: &mut String,
        file_path: &PathBuf
    ) -> Vec<Token> {
        let file = file_path.to_str().unwrap().to_string();

        // First replace all the imports
        let processed_contents = process_imports(contents, file_path);
        *contents = processed_contents;

        // Initialize variables for tokens, current token, and position tracking
        let mut tokens: Vec<Token> = Vec::new();
        let mut current_token: String = String::new();
        let mut current_line: u32 = 1;
        let mut current_column: u32 = 1;

        // Collect all characters from the contents
        let characters = contents.chars().collect::<Vec<char>>();
        let characters_len = characters.len();

        // Iterate over each character in the input
        for i in 0 .. characters.len() {
            let tokens_len = tokens.len();
            let character = characters[i];

            // Handle newline characters
            if character == '\n' {
                if self.in_comment {
                    self.in_comment = false; // Exit single-line comment
                    current_token.clear();
                }
                
                current_line += 1;
                current_column = 1; // Reset column for the new line
                continue;
            } else {
                current_column += 1; // Increment column for other characters
            }

            // Skip processing if inside a single-line comment
            if self.in_comment {
                continue;
            }
            
            // Handle comment starts ('//' and '/*')
            if !self.in_string && !self.in_block_comment && character == '/' {
                if characters_len < 2 {
                    continue; // Ensure there's a next character
                }
                if characters[i + 1] == '*' {
                    self.in_block_comment = true; // Start block comment
                } else if characters[i + 1] == '/' {
                    self.in_comment = true; // Start single-line comment
                }
                current_token.clear();
            } 
            // Handle block comment end ('*/')
            else if self.in_block_comment {
                if character == '/' && characters[i - 1] == '*' {
                    self.in_block_comment = false; // Exit block comment
                }
                continue;
            } 
            // Handle string literal processing
            else if self.in_string {
                if character == '"' && !self.in_escape {
                    self.in_string = false; // End string literal

                    // Push the completed string literal token
                    tokens.push(
                        Token::new(
                            TokenType::StringLiteral,
                            current_token.clone(),
                            String::from(file.clone()),
                            current_line,
                            current_column
                        )
                    );
                    current_token.clear();
                } else if character == '\\' && !self.in_escape {
                    self.in_escape = true; // Start escape sequence
                    current_token.push(character);
                } else {
                    current_token.push(character);
                    self.in_escape = character == '\\'; // Track escape sequence state
                }
                continue;
            } 
            // Start of a string literal
            else if character == '"' {
                self.in_string = true;
                continue;
            } 
            // Handle whitespace
            else if character.is_whitespace() {
                if !current_token.is_empty() {
                    // Push completed token and clear current token
                    tokens.push(Lexer::calculate(
                        current_token.trim(),
                        &file,
                        &current_line,
                        &current_column
                    ));
                    current_token.clear();
                }
            } 
            // Handle punctuation characters
            else if PUNCTUATION_CHARS.contains(&character) {
                // Allow decimals in numeric literals
                if character == '.' && !current_token.parse::<i128>().is_err() {
                    current_token.push(character);
                    continue;
                }
                if !current_token.is_empty() {
                    // Push the preceding token
                    tokens.push(Lexer::calculate(
                        &current_token,
                        &file,
                        &current_line,
                        &current_column
                    ));
                    current_token.clear();
                }

                // Handle specific token combinations (e.g., brackets or operators)
                if character == ']' {
                    if tokens_len < 2 {
                        continue; // Ensure enough tokens for lookbehind
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
                    {
                        // Combine into a type with array notation
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
                } 
                // Handle operator combinations (e.g., '==', '->', '++')
                else if character == '=' || character == '-' || character == '+' {
                    if tokens_len == 0 {
                        continue; // Ensure there's a preceding token
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
                } 
                // Handle arrow operator ('->')
                else if character == '>' {
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

                // Push the current punctuation token
                tokens.push(
                    Lexer::calculate(
                        &character.to_string(),
                        &file,
                        &current_line,
                        &current_column
                    )
                );
            } 
            // Append character to the current token
            else {
                current_token.push(character);
            }
        }

        // Return the list of tokens
        tokens
    }

    // Create a new Lexer instance with default values
    fn new() -> Self {
        Lexer {
            in_string: false,
            in_escape: false,
            in_comment: false,
            in_block_comment: false
        }
    }

    // Calculate the Token for a given token string and context
    fn calculate(current_token: &str, file: &String, line: &u32, col: &u32) -> Token {
        // Determine the TokenType based on known tokens or patterns
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

        // Create and return the Token
        Token::new(
            token_type.clone(),
            current_token.to_string(),
            file.clone(),
            *line,
            *col
        )
    }
}
