use std::collections::HashMap;

use class::Class;
use function::Function;

pub mod class;
pub mod function;
pub mod wrapper;

// Surf's custom wrapper for CLang's AST
pub struct Header {
    classes: HashMap<String, Class>,
    functions: HashMap<String, Function>
}

pub trait HeaderImpl {
    fn new() -> Self;

    fn add_class(&mut self, name: String, class: Class);
    fn add_function(&mut self, name: String, function: Function);

    fn get_classes(&self) -> &HashMap<String, Class>;
    fn get_functions(&self) -> &HashMap<String, Function>;
    fn find_class(&mut self, name: &str) -> Option<&mut Class>;
}

impl HeaderImpl for Header {
    fn new() -> Self {
        Header {
            classes: HashMap::new(),
            functions: HashMap::new()
        }
    }

    fn add_class(&mut self, name: String, class: Class) {
        self.classes.insert(name, class);
    }

    fn add_function(&mut self, name: String, function: Function) {
        self.functions.insert(name, function);
    }

    fn get_classes(&self) -> &HashMap<String, Class> {
        &self.classes
    }

    fn get_functions(&self) -> &HashMap<String, Function> {
        &self.functions
    }

    fn find_class(&mut self, name: &str) -> Option<&mut Class> {
        self.classes.get_mut(name)
    }
}