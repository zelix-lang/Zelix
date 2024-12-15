pub mod colors {
    use std::collections::HashMap;

    use lazy_static::lazy_static;

    lazy_static! {
        pub static ref color_map : HashMap<String, String> = {
            let mut map = HashMap::new();
            map.insert("black".to_string(), "\x1b[30m".to_string());
            map.insert("red".to_string(), "\x1b[31m".to_string());
            map.insert("green".to_string(), "\x1b[32m".to_string());
            map.insert("yellow".to_string(), "\x1b[33m".to_string());
            map.insert("blue".to_string(), "\x1b[34m".to_string());
            map.insert("magenta".to_string(), "\x1b[35m".to_string());
            map.insert("cyan".to_string(), "\x1b[36m".to_string());
            map.insert("white".to_string(), "\x1b[37m".to_string());
            map.insert("black_bright".to_string(), "\x1b[90m".to_string());
            map.insert("red_bright".to_string(), "\x1b[91m".to_string());
            map.insert("green_bright".to_string(), "\x1b[92m".to_string());
            map.insert("yellow_bright".to_string(), "\x1b[93m".to_string());
            map.insert("blue_bright".to_string(), "\x1b[94m".to_string());
            map.insert("magenta_bright".to_string(), "\x1b[95m".to_string());
            map.insert("cyan_bright".to_string(), "\x1b[96m".to_string());
            map.insert("white_bright".to_string(), "\x1b[97m".to_string());
            map.insert("reset".to_string(), "\x1b[0m".to_string());
            map
        };

        // Store reset by default because it's frequently used
        pub static ref RESET : String = color_map.get("reset").unwrap().clone();
    }
}

pub fn get_color(name: &String) -> &String {
    colors::color_map.get(name).unwrap_or_else(|| &colors::RESET)
}