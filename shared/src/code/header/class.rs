use std::collections::HashMap;

use super::function::Function;

pub struct Class {
    methods: HashMap<String, Function>
}

pub trait ClassImpl {
    fn new() -> Self;

    fn add_method(&mut self, name: String, function: Function);

    fn get_methods(&self) -> &HashMap<String, Function>;
}

impl ClassImpl for Class {
    fn new() -> Self {
        Class {
            methods: HashMap::new()
        }
    }

    fn add_method(&mut self, name: String, function: Function) {
        self.methods.insert(name, function);
    }

    fn get_methods(&self) -> &HashMap<String, Function> {
        &self.methods
    }
}