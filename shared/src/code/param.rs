use code::token::Token;

#[derive(Debug, Clone)]
pub struct Param {

    name: String,
    data_type: Vec<Token>,
    trace: String,
    // Used to know if a param is a reference
    // (e.g. &str)
    is_reference: bool

}

pub trait ParamImpl {

    fn new(name: String, data_type: Vec<Token>, trace: String, is_reference: bool) -> Self;

    fn get_name(&self) -> &String;
    fn get_data_type(&self) -> &Vec<Token>;
    fn get_trace(&self) -> &String;
    fn is_reference(&self) -> bool;

}

impl ParamImpl for Param {

    fn new(name: String, data_type: Vec<Token>, trace: String, is_reference: bool) -> Self {
        Param {
            name,
            data_type,
            trace,
            is_reference
        }
    }
    
    fn get_name(&self) -> &String {
        &self.name
    }

    fn get_data_type(&self) -> &Vec<Token> {
        &self.data_type
    }

    fn get_trace(&self) -> &String {
        &self.trace
    }

    fn is_reference(&self) -> bool {
        self.is_reference
    }

}