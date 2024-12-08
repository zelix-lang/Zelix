pub struct Variable {
    // Used to check lifetime of the variable
    is_reference_to_param: bool
}

pub trait VariableImpl {
    fn new(is_reference_to_param: bool) -> Self;
    fn is_reference_to_param(&self) -> bool;
}

impl VariableImpl for Variable {
    fn new(is_reference_to_param: bool) -> Self {
        Variable {
            is_reference_to_param,
        }
    }

    fn is_reference_to_param(&self) -> bool {
        self.is_reference_to_param
    }

}