use std::collections::HashMap;

use super::function::Function;

#[derive(Debug, Clone)]
pub struct Class {
    methods: HashMap<String, Function>,
    generic_count: usize
}

pub trait ClassImpl {
    fn new(generic_count: usize) -> Self;
    fn add_method(&mut self, name: String, function: Function);
    fn get_methods(&self) -> &HashMap<String, Function>;
    fn get_generic_count(&self) -> usize;
}

impl ClassImpl for Class {
    fn new(generic_count: usize) -> Self {
        Class {
            methods: HashMap::new(),
            generic_count
        }
    }

    fn add_method(&mut self, name: String, function: Function) {
        self.methods.insert(name, function);
    }

    fn get_methods(&self) -> &HashMap<String, Function> {
        &self.methods
    }

    fn get_generic_count(&self) -> usize {
        self.generic_count
    }
}