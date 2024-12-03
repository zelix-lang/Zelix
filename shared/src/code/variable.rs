use code::types::ParamType;

pub struct Variable {
    var_type: ParamType,
    is_reference: bool,
}

pub trait VariableImpl {
    fn new(var_type: ParamType, is_reference: bool) -> Self;
    fn get_var_type(&self) -> &ParamType;
    fn is_reference(&self) -> bool;
}

impl VariableImpl for Variable {

    fn new(var_type: ParamType, is_reference: bool) -> Self {
        Variable {
            var_type,
            is_reference
        }
    }

    fn get_var_type(&self) -> &ParamType {
        &self.var_type
    }

    fn is_reference(&self) -> bool {
        self.is_reference
    }

}