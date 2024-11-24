use std::collections::HashMap;

use super::{function::Function, import::Import};

#[derive(Debug, Clone)]
pub struct FileCode {

    functions: HashMap<String, Function>,
    imports: Vec<Import>
    
}

pub trait FileCodeImpl {

    fn new() -> Self;

    fn add_function(&mut self, name: String, function: Function);
    fn add_import(&mut self, import: Import);

    fn get_functions(&self) -> &HashMap<String, Function>;
    fn get_imports(&self) -> &Vec<Import>;
    
}

impl FileCodeImpl for FileCode {

    fn new() -> Self {
        FileCode {
            functions: HashMap::new(),
            imports: Vec::new()
        }
    }

    fn add_function(&mut self, name: String, function: Function) {
        self.functions.insert(name, function);
    }

    fn add_import(&mut self, import: Import) {
        self.imports.push(import);
    }

    fn get_functions(&self) -> &HashMap<String, Function> {
        &self.functions
    }

    fn get_imports(&self) -> &Vec<Import> {
        &self.imports
    }
    
}