use cranelift::prelude::*;
use cranelift_module::{Module, Linkage};
use cranelift_object::{ObjectBuilder, ObjectModule};
use lexer::token::Token;
use logos::Lexer;
use shared::result::try_unwrap;
use target_lexicon::Triple;
use std::fs::File;
use std::io::Write;
use std::path::PathBuf;
use cranelift_codegen::settings::Flags;

pub fn compile(tokens: Lexer<'_, Token>, out_dir: PathBuf) {
    // New triple for the target
    let triple = Triple::host();

    let builder = try_unwrap(
        ObjectBuilder::new(
            try_unwrap(
                cranelift_codegen::isa::lookup(triple.clone()).unwrap().finish(
                    Flags::new(cranelift_codegen::settings::builder()),
                ),
                "Failed to create target machine",
            ),
            "surf_lang".to_string(),
            cranelift_module::default_libcall_names(),
        ),
        "Failed to create object builder",
    );

    let mut module = ObjectModule::new(builder);

}