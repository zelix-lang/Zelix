use std::collections::HashMap;
use shared::token::token_type::TokenType;

lazy_static::lazy_static! {
    // A vector of punctuation characters to be used in tokenization or parsing
    pub static ref PUNCTUATION_CHARS: Vec<char> = ";,(){}:+\\->%<.=![]/|*^&".chars()
        .collect();

    // HashMap to store known tokens and their corresponding TokenType
    pub static ref KNOWN_TOKENS: HashMap<String, TokenType> = {
        let mut map = HashMap::new();

        // Keywords
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

        // Operators and symbols
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
        map.insert(String::from("<="), TokenType::LessThanOrEqual);
        map.insert(String::from(">="), TokenType::GreaterThanOrEqual);
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
        map.insert(String::from("%"), TokenType::Percent);

        // Data types
        map.insert(String::from("str"), TokenType::String);
        map.insert(String::from("num"), TokenType::Num);
        map.insert(String::from("nothing"), TokenType::Nothing);
        map.insert(String::from("bool"), TokenType::Bool);

        // Access modifiers
        map.insert(String::from("pub"), TokenType::Pub);

        // Array types
        map.insert(String::from("str[]"), TokenType::StringArray);
        map.insert(String::from("num[]"), TokenType::NumArray);
        map.insert(String::from("bool[]"), TokenType::BoolArray);

        // Special annotations
        map.insert(String::from("[discrete]"), TokenType::Discrete);

        map
    };
}
