use code::types::ParamType;

pub struct Variable {
    var_type: ParamType,
    // Used to check lifetime of the variable
    is_reference_to_heap: bool,
    is_reference_to_param: bool,
    // Used to validate pointer operations
    is_referene: bool,
}

pub trait VariableImpl {
    fn new(var_type: ParamType, is_reference_to_heap: bool, is_reference_to_param: bool, is_referene: bool) -> Self;
    fn get_var_type(&self) -> &ParamType;
    fn is_reference_to_heap(&self) -> bool;
    fn is_reference_to_param(&self) -> bool;
    fn is_referene(&self) -> bool;
}

impl VariableImpl for Variable {
    fn new(var_type: ParamType, is_reference_to_heap: bool, is_reference_to_param: bool, is_referene: bool) -> Self {
        Variable {
            var_type,
            is_reference_to_heap,
            is_reference_to_param,
            is_referene,
        }
    }

    fn get_var_type(&self) -> &ParamType {
        &self.var_type
    }

    fn is_reference_to_heap(&self) -> bool {
        self.is_reference_to_heap
    }

    fn is_reference_to_param(&self) -> bool {
        self.is_reference_to_param
    }

    fn is_referene(&self) -> bool {
        self.is_referene
    }
}