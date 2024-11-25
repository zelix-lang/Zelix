use shared::token::{token::{Token, TokenImpl}, token_type::TokenType};
pub mod data_types;

mod globals {
    use std::collections::HashMap;

    use fancy_regex::Regex;
    use shared::token::token_type::TokenType;

    lazy_static::lazy_static! {
        pub static ref NUMBER_REGEX:Regex = Regex::new(r#"^\d+((\.\d+)?)$"#).unwrap();
        pub static ref PUNCTUATION_CHARS: Vec<char> = ";,(){}:+\\->%<.=![]/|*^&".chars()
            .collect();

        pub static ref knwon_tokens : HashMap<String, TokenType> = {
            let mut map = HashMap::new();
            map.insert(String::from("return"), TokenType::Return);
            map.insert(String::from("fun"), TokenType::Function);
            map.insert(String::from("let"), TokenType::Let);
            map.insert(String::from("const"), TokenType::Const);
            map.insert(String::from("while"), TokenType::While);
            map.insert(String::from("for"), TokenType::For);
            map.insert(String::from("break"), TokenType::Break);
            map.insert(String::from("continue"), TokenType::Continue);
            map.insert(String::from("in"), TokenType::In);
            map.insert(String::from("if"), TokenType::If);
            map.insert(String::from("else"), TokenType::Else);
            map.insert(String::from("elseif"), TokenType::ElseIf);
            map.insert(String::from("unsafe"), TokenType::Unsafe);
            map.insert(String::from("="), TokenType::Assign);
            map.insert(String::from("+"), TokenType::Plus);
            map.insert(String::from("-"), TokenType::Minus);
            map.insert(String::from("++"), TokenType::Increment);
            map.insert(String::from("--"), TokenType::Decrement);
            map.insert(String::from("*"), TokenType::Asterisk);
            map.insert(String::from("/"), TokenType::Slash);
            map.insert(String::from("<"), TokenType::LessThan);
            map.insert(String::from(">"), TokenType::GreaterThan);
            map.insert(String::from("+="), TokenType::AssignAdd);
            map.insert(String::from("-="), TokenType::AssignSub);
            map.insert(String::from("*="), TokenType::AssignAsterisk);
            map.insert(String::from("/="), TokenType::AssignSlash);
            map.insert(String::from("=="), TokenType::Equal);
            map.insert(String::from("!="), TokenType::NotEqual);
            map.insert(String::from("<="), TokenType::LessThan);
            map.insert(String::from(">="), TokenType::GreaterThan);
            map.insert(String::from("&"), TokenType::Ampersand);
            map.insert(String::from("|"), TokenType::Bar);
            map.insert(String::from("^"), TokenType::Xor);
            map.insert(String::from("!"), TokenType::Not);
            map.insert(String::from(","), TokenType::Comma);
            map.insert(String::from(";"), TokenType::Semicolon);
            map.insert(String::from("("), TokenType::OpenParen);
            map.insert(String::from(")"), TokenType::CloseParen);
            map.insert(String::from("{"), TokenType::OpenCurly);
            map.insert(String::from("}"), TokenType::CloseCurly);
            map.insert(String::from(":"), TokenType::Colon);
            map.insert(String::from("->"), TokenType::Arrow);
            map.insert(String::from("["), TokenType::OpenBracket);
            map.insert(String::from("]"), TokenType::CloseBracket);
            map.insert(String::from("."), TokenType::Dot);
            map.insert(String::from("%"), TokenType::Slash);
            map.insert(String::from("str"), TokenType::String);
            map.insert(String::from("num"), TokenType::Num);
            map.insert(String::from("nothing"), TokenType::Nothing);
            map.insert(String::from("bool"), TokenType::Bool);
            map.insert(String::from("pub"), TokenType::Pub);
            map.insert(String::from("import"), TokenType::Import);
            map.insert(String::from("from"), TokenType::From);
            map.insert(String::from("as"), TokenType::As);
            map.insert(String::from("str[]"), TokenType::StringArray);
            map.insert(String::from("num[]"), TokenType::NumArray);
            map.insert(String::from("bool[]"), TokenType::BoolArray);
            map.insert(String::from("[discrete]"), TokenType::Discrete);

            map
        };
    }

}

pub struct Lexer {

    in_string: bool,
    in_escape: bool,
    in_comment: bool,
    in_block_comment: bool

}

pub trait LexerImpl {
    fn tokenize(
        &mut self,
        contents: &String,
        file: &String
    ) -> Vec<Token>;

    fn new() -> Self;
    fn calculate(current_token: &str, file: &String, line: &u32, col: &u32) -> Token;
}

impl LexerImpl for Lexer {
    fn tokenize(
        &mut self,
        contents: &String,
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
            } else if globals::PUNCTUATION_CHARS.contains(&character) {
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
        let token_type = match globals::knwon_tokens.get(current_token) {
            Some(token_type) => token_type,
            None => {
                if globals::NUMBER_REGEX.is_match(current_token).unwrap() {
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