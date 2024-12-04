use rand::{thread_rng, Rng};

fn generate_random_number(min_size: usize, max_size: usize) -> usize {
    let mut rng = thread_rng();
    rng.gen_range(min_size..max_size)
}

fn generate_string(length: usize) -> String {
    const CHARACTERS: &str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_";
    let characters: Vec<char> = CHARACTERS.chars().collect();

    (0..length)
        .map(|_| characters[
            generate_random_number(0, characters.len())
        ])
        .collect()
}

/// Creates a random string used for prefixing
/// functions, which avoids multiple redefinitions
pub fn create_random_prefix() -> String {
    // Generate a random length for the prefix
    let length = generate_random_number(5, 10);

    generate_string(length)
}