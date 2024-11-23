use cranelift::prelude::*;
use cranelift_module::{Module, Linkage};
use cranelift_object::{ObjectBuilder, ObjectModule};
use target_lexicon::Triple;
use std::fs::File;
use std::io::Write;
use cranelift_codegen::settings::Flags;

fn main() {
    // Create the target triple (e.g., native platform)
    let triple = Triple::host();

    let builder = ObjectBuilder::new(
        cranelift_codegen::isa::lookup(triple.clone()).unwrap().finish(
            Flags::new(cranelift_codegen::settings::builder()),
        ).expect("Failed to create target machine"),
        "basic_lang".to_string(),
    )
    cranelift_module::default_libcall_names(),
    .unwrap();
    let mut module = ObjectModule::new(builder);

    // Define a simple function
    let mut ctx = module.make_context();
    let mut builder_ctx = FunctionBuilderContext::new();
    
    ctx.func.signature.params.push(AbiParam::new(types::I64));
    ctx.func.signature.params.push(AbiParam::new(types::I64));
    ctx.func.signature.returns.push(AbiParam::new(types::I64));

    {
        let mut builder = FunctionBuilder::new(&mut ctx.func, &mut builder_ctx);
        let entry_block = builder.create_block();
        
        builder.append_block_params_for_function_params(entry_block);
        builder.switch_to_block(entry_block);
        builder.seal_block(entry_block);

        let params: Vec<Value> = builder.block_params(entry_block).to_vec();
        
        let sum = {
            let ins = builder.ins();
            ins.iadd(params[0], params[1])
        };
        builder.ins().return_(&[sum]);
    }

    let func_id = module
        .declare_function("add", Linkage::Export, &ctx.func.signature)
        .unwrap();

    module.define_function(func_id, &mut ctx).unwrap();
    module.clear_context(&mut ctx);

    // Create a main function that calls the add function
    {
        let mut ctx = module.make_context();
        let mut builder_ctx = FunctionBuilderContext::new();
        
        ctx.func.signature.returns.push(AbiParam::new(types::I64));

        {
            let mut builder = FunctionBuilder::new(&mut ctx.func, &mut builder_ctx);
            let entry_block = builder.create_block();
            
            builder.switch_to_block(entry_block);
            builder.seal_block(entry_block);

            // main function has no parameters, we use constants
            let param1 = builder.ins().iconst(types::I64, 25);
            let param2 = builder.ins().iconst(types::I64, 25);

            // call the add function
            let add_func = module.declare_func_in_func(func_id, &mut builder.func);
            let call = builder.ins().call(add_func, &[param1, param2]);
            let sum = builder.inst_results(call)[0];
            
            builder.ins().return_(&[sum]);
        }

        let func_id = module
            .declare_function("main", Linkage::Export, &ctx.func.signature)
            .unwrap();

        module.define_function(func_id, &mut ctx).unwrap();
        module.clear_context(&mut ctx);
    }

    // Emit object file
    let obj = module.finish();
    let mut file = File::create("output.o")
        .expect("Failed to create output file");
    file.write_all(&obj.emit()
                        .expect("Failed to emit object file")
    ).expect("Failed to write object file");

    println!("Object file generated: output.o");
    // execute clang to generate the executable
    let output = std::process::Command::new("clang")
        .args(&["output.o", "-o", "output"])
        .output()
        .expect("Failed to execute clang");

    println!("{}", String::from_utf8_lossy(&output.stdout));
}