use clang::TypeKind;

pub struct Function {
    return_type: TypeKind,
    params: Vec<TypeKind>,
}

pub trait FunctionImpl {
    fn new(return_type: &TypeKind) -> Self;

    fn add_param(&mut self, type_: TypeKind);

    fn get_return_type(&self) -> &TypeKind;
    fn get_params(&self) -> &Vec<TypeKind>;
}

impl FunctionImpl for Function {
    fn new(return_type: &TypeKind) -> Self {
        Function {
            return_type: return_type.clone(),
            params: Vec::new(),
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
}