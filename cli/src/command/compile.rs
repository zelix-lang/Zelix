use std::path::PathBuf;

use crate::command::lexe_base::lexe_base;

pub fn compile_command(path: Option<PathBuf>) {

    println!("{:?}", lexe_base(path));

}