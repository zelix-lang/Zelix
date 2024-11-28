use std::path::PathBuf;

#[derive(Debug, Clone, PartialEq)]
pub struct Import {

    /// Where to import from
    from: PathBuf,

    /// The trace of the import chain
    trace: String

}

pub trait Importable {

    fn new(from: PathBuf, trace: String) -> Self;
    fn get_from(&self) -> PathBuf;
    fn get_trace(&self) -> String;
    
}

impl Importable for Import {

    fn new(from: PathBuf, trace: String) -> Import {
        Import {
            from,
            trace
        }
    }

    fn get_from(&self) -> PathBuf {
        self.from.clone()
    }

    fn get_trace(&self) -> String {
        self.trace.clone()
    }

}