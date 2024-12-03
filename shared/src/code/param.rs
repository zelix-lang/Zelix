use code::token_type::TokenType;

#[derive(Debug, Clone)]
pub struct Param {

    name: String,
    data_type: TokenType,
    trace: String

}

pub trait ParamImpl {

    fn new(name: String, data_type: TokenType, trace: String) -> Self;

    fn get_name(&self) -> &String;
    fn get_data_type(&self) -> &TokenType;
    fn get_trace(&self) -> &String;

}

impl ParamImpl for Param {

    fn new(name: String, data_type: TokenType, trace: String) -> Self {
        Param {
            name,
            data_type,
            trace
        }
    }
    
    fn get_name(&self) -> &String {
        &self.name
    }

    fn get_data_type(&self) -> &TokenType {
        &self.data_type
    }

    fn get_trace(&self) -> &String {
        &self.trace
    }

}