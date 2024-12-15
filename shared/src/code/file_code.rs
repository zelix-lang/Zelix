use std::{collections::HashMap, path::PathBuf, process::exit};

use logger::{Logger, LoggerImpl};

use super::function::{Function, FunctionImpl};

pub struct FileCode {

    functions: HashMap<String, HashMap<String, Function>>,
    source: PathBuf
    
}

pub trait FileCodeImpl {

    fn new(source: PathBuf) -> Self;

    fn add_function(&mut self, file: String, name: String, function: Function);

    fn get_functions(&self) -> &HashMap<String, HashMap<String, Function>>;
    fn get_source(&self) -> &PathBuf;
    
}

impl FileCodeImpl for FileCode {

    fn new(source: PathBuf) -> Self {
        FileCode {
            functions: HashMap::new(),
            source
        }
    }

    fn add_function(&mut self, name: String, file: String, function: Function) {
        if !self.functions.contains_key(&file) {
            self.functions.insert(file.clone(), HashMap::new());
        }

        let funcs = self.functions.get_mut(&file).unwrap();

        if funcs.contains_key(&name) {
            Logger::err(
                format!("Duplicate function name: {}", name).as_str(),
                &[
                    "Function names must be unique"
                ],
                &[
                    function.get_trace().as_str()
                ]
            );

            exit(1);
        }

        funcs.insert(name, function);
    }

    fn get_functions(&self) -> &HashMap<String, HashMap<String, Function>> {
        &self.functions
    }

    fn get_source(&self) -> &PathBuf {
        &self.source
    }
    
}