use clang::Entity;

use crate::env::STANDARD_LIBRARY_LOCATION;

// Reads the AST from the parsed initial root given by clang
pub fn read_ast(entity: Entity) -> Vec<Entity> {
    let mut result : Vec<Entity> = Vec::new();
    result.push(entity.clone());

    // We don't know the total depth level, so we use usize
    let mut idx : usize = 0;

    // Avoid recursion and use a loop
    loop {
        if idx >= result.len() {
            break;
        }

        let element = result[idx];
        let children = element.get_children();

        // Nothing more to read, break the loop
        if children.is_empty() && idx == result.len() - 1 {
            break;
        }

        // Add the children to the result
        for child in children {
            let location_optional = child.get_location();

            // No location, skip
            if location_optional.is_none() {
                continue;
            }

            let location = location_optional.unwrap().get_file_location();
            let file_optional = location.file;

            // No file, skip
            if file_optional.is_none() {
                continue;
            }

            // Skip if the file is not in the standard library
            if !file_optional.unwrap().get_path().starts_with(STANDARD_LIBRARY_LOCATION.clone()) {
                continue;
            }

            result.push(child);
        }

        idx += 1;
    }

    result
}