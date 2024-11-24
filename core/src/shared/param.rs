use shared::token::token_type::TokenType;

#[derive(Debug, Clone)]
pub struct Param {

    name: String,
    data_type: TokenType

}

pub trait ParamImpl {

    fn new(name: String, data_type: TokenType) -> Self;

    fn get_name(&self) -> &String;
    fn get_data_type(&self) -> &TokenType;

}

impl ParamImpl for Param {

    fn new(name: String, data_type: TokenType) -> Self {
        Param {
            name,
            data_type
        }
    }
    
    fn get_name(&self) -> &String {
        &self.name
    }

    fn get_data_type(&self) -> &TokenType {
        &self.data_type
    }

}