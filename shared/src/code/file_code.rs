use std::{collections::HashMap, path::PathBuf, process::exit};

use clang::Index;
use c_parser::{header::Header, parse_header_file, wrapper::wrap_header};
use logger::{Logger, LoggerImpl};

use super::{function::{Function, FunctionImpl}, header_reader::read_ast, import::{Import, Importable}};

pub struct FileCode {

    functions: HashMap<String, Function>,
    imports: Vec<Header>,
    seen_imports: Vec<Import>,
    source: PathBuf
    
}

pub trait FileCodeImpl {

    fn new(source: PathBuf) -> Self;

    fn add_function(&mut self, name: String, function: Function);
    fn add_import(&mut self, import: Import, index: &Index);

    fn get_functions(&self) -> &HashMap<String, Function>;
    fn get_imports(&self) -> &Vec<Header>;
    fn get_seen_imports(&self) -> &Vec<Import>;
    fn get_source(&self) -> &PathBuf;
    
}

impl FileCodeImpl for FileCode {

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

    fn add_import(&mut self, import: Import, index: &Index) {
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
        let translation_unit: clang::TranslationUnit<'_> = parse_header_file(
            &import.get_from().to_str().unwrap().to_string(), 
            index
        );

        let ast = read_ast(translation_unit.get_entity());
        let wrapped_ast = wrap_header(ast);
        self.imports.push(wrapped_ast);
    }

    fn get_functions(&self) -> &HashMap<String, Function> {
        &self.functions
    }

    fn get_imports(&self) -> &Vec<Header> {
        &self.imports
    }

    fn get_seen_imports(&self) -> &Vec<Import> {
        &self.seen_imports
    }

    fn get_source(&self) -> &PathBuf {
        &self.source
    }
    
}