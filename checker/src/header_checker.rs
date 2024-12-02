use shared::code::header::{class::Class, Header, HeaderImpl};

pub fn find_imported_classes(check_for: &String, headers: &Vec<Header>) -> Option<Class> {
    for header in headers {
        let classes = header.get_classes();
        if classes.contains_key(check_for) {
            return Some(
                classes.get(check_for)
                    .unwrap()
                    .clone()
            );
        }
    }

    None
}

pub fn check_header_value_definition(check_for: &String, headers: &Vec<Header>) -> bool {
    for header in headers {
        let functions = header.get_functions();
        if functions.contains_key(check_for) {
            return true;
        }

        let classes = header.get_classes();
        if classes.contains_key(check_for) {
            return true;
        }
    }

    false
}