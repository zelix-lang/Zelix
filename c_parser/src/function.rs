use clang::TypeKind;

#[derive(Debug, Clone)]
pub struct Function {
    return_type: TypeKind,
    params: Vec<TypeKind>,
    generic_count: usize
}

pub trait FunctionImpl {
    fn new(return_type: &TypeKind, generic_count: usize) -> Self;

    fn add_param(&mut self, type_: TypeKind);

    fn get_return_type(&self) -> &TypeKind;
    fn get_params(&self) -> &Vec<TypeKind>;
    fn get_generic_count(&self) -> usize;
}

impl FunctionImpl for Function {
    fn new(return_type: &TypeKind, generic_count: usize) -> Self {
        Function {
            return_type: return_type.clone(),
            params: Vec::new(),
            generic_count
        }
    }

    fn add_param(&mut self, type_: TypeKind) {
        self.params.push(type_);
    }

    fn get_return_type(&self) -> &TypeKind {
        &self.return_type
    }

    fn get_params(&self) -> &Vec<TypeKind> {
        &self.params
    }

    fn get_generic_count(&self) -> usize {
        self.generic_count
    }
}