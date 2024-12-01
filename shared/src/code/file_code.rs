use std::{collections::HashMap, path::PathBuf, process::exit};

use clang::Entity;
use parser::{create_c_instance, create_index, parse_header_file};

use crate::logger::{Logger, LoggerImpl};

use super::{function::{Function, FunctionImpl}, header_reader::read_ast, import::{Import, Importable}};

#[derive(Debug, Clone)]
pub struct FileCode<'a> {

    functions: HashMap<String, Function>,
    imports: Vec<Entity<'a>>,
    seen_imports: Vec<Import>,
    source: PathBuf
    
}

pub trait FileCodeImpl {

    fn new(source: PathBuf) -> Self;

    fn add_function(&mut self, name: String, function: Function);
    fn add_import(&mut self, import: Import);

    fn get_functions(&self) -> &HashMap<String, Function>;
    fn get_imports(&self) -> &Vec<Entity>;
    fn get_seen_imports(&self) -> &Vec<Import>;
    fn get_source(&self) -> &PathBuf;
    
}

impl FileCodeImpl for FileCode<'_> {

    fn new(source: PathBuf) -> Self {
        FileCode {
            functions: HashMap::new(),
            imports: Vec::new(),
            seen_imports: Vec::new(),
            source
        }
    }

    fn add_function(&mut self, name: String, function: Function) {
        if self.functions.contains_key(&name) {
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

        self.functions.insert(name, function);
    }

    fn add_import(&mut self, import: Import) {
        // Don't include duplicate imports
        if self.seen_imports.contains(&import) {
            return;
        }

        self.seen_imports.push(import.clone());

        // Since imports to Surf files are always rendered before lexing
        // the only possible way this method is called is for imports
        // pointing to the standard library, which we have said before
        // to have a .hpp, .h file extension always, so we don't need to
        // check for the file extension here

        // Gather the information from the .hpp or .h file
        let c_instance = create_c_instance();
        let index = create_index(&c_instance);
        let translation_unit = parse_header_file(
            &import.get_from().to_str().unwrap().to_string(), 
            &index
        );
        let entity = translation_unit.get_entity();
        read_ast(&entity);
    }

    fn get_functions(&self) -> &HashMap<String, Function> {
        &self.functions
    }

    fn get_imports(&self) -> &Vec<Entity> {
        &self.imports
    }

    fn get_seen_imports(&self) -> &Vec<Import> {
        &self.seen_imports
    }

    fn get_source(&self) -> &PathBuf {
        &self.source
    }
    
}