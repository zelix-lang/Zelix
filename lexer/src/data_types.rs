use shared::token::token_type::TokenType;

mod globals{
    use lazy_static::lazy_static;
    use shared::token::token_type::TokenType;

    lazy_static! {
        pub static ref DATA_TYPES: Vec<TokenType> = vec![
            TokenType::StringArray,
            TokenType::NumArray,
            TokenType::BoolArray,
            TokenType::String,
            TokenType::Num,
            TokenType::Bool,
            TokenType::Discrete,
            TokenType::Nothing
        ];
    }
}

pub fn is_data_type(token_type: TokenType) -> bool {
    globals::DATA_TYPES.contains(&token_type)
}