use clang::{Entity, EntityKind, TypeKind};

use super::{class::{Class, ClassImpl}, function::{Function, FunctionImpl}, Header, HeaderImpl};

fn find_data_type(entity: &Entity) -> TypeKind {
    // We're ensured that at this point the entity doesn't need its children ASTs
    // to be processed, so we just focus on the given entity
    let return_type_option = entity.get_result_type();

    match return_type_option {
        Some(return_type) => return_type.get_kind(),
        // No data type, return void
        None => TypeKind::Void
    }
}

fn find_parameters(entity: &Entity) -> Vec<TypeKind> {
    let mut result = Vec::new();
    let arguments_optional = entity.get_arguments();

    // No arguments, return empty
    if arguments_optional.is_none() {
        return result;
    }

    // Find the data type of each parameter and push it to the result
    for param in arguments_optional.unwrap() {
        result.push(
            find_data_type(&param)
        );
    }

    result
}

fn process_function(entity: &Entity)-> Function {
    // Wrap the function
    let return_type = find_data_type(entity);
    let mut func = Function::new(&return_type);

    for param in find_parameters(entity) {
        func.add_param(param);
    }

    func
}

pub fn wrap_header(ast: Vec<Entity>) -> Header {
    let mut result = Header::new();

    for entity in ast {
        let kind = entity.get_kind();

        match kind {
            EntityKind::FunctionDecl => {
                // Add the function to the header
                result.add_function(
                    entity.get_name().unwrap(),
                     process_function(&entity)
                );
            }

            EntityKind::ClassDecl => {
                // Wrap the class
                let class = Class::new();
                // Add the class to the header
                result.add_class(
                    entity.get_name().unwrap(),
                    class
                );
            }

            EntityKind::Method => {
                // First find the class
                let class_name = entity.get_semantic_parent()
                        .and_then(|p| p.get_name())
                        .unwrap_or_default();

                // Find the class
                let class = result.find_class(&class_name)
                    .unwrap_or_else(|| 
                        panic!(
                            "Class {} not found. This is likely a misconfiguration of the standard library", 
                            class_name
                        )
                    );

                // Wrap the function
                let function = process_function(&entity);

                // Add the method to the class
                class.add_method(
                    entity.get_name().unwrap(),
                    function
                );
            }

            _ => {}
        }
    }

    result
}