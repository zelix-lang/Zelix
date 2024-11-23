/*
    These files are part of the Surf's standard library.
    They're bundled with the main program by the compiler
    which is then converted to machine code.

    -----
    License notice:

    This code is released under the GNU GPL v3 license.
    The code is provided as is without any warranty
    Copyright (c) 2024 Rodrigo R. & all Surf contributors
*/

#include <iostream>
#include <string>
#include <fstream>
#include <unistd.h>
#include <vector>
#include <dirent.h>
#include <filesystem>
#include "../lang/result.h"
#include "../lang/err.h"

Result<bool> write_file(const std::string* path, const std::string* content) {

    // First check if the file exists
    if (!file_exists(path).unwrap()) {
        return Result(false, optional<Err>(Err("File does not exist")));
    }

    // Open the file
    std::ofstream file;
    file.open(*path);

    // Something went wrong
    if (!file.is_open()) {
        return Result(false, optional<Err>(Err("Failed to create file")));
    }

    // Directly write the content to the file
    file << content;

    // Don't forget to close the file!
    file.close();

    return Result(true, optional<Err>());

}

Result<std::string> read_file(const std::string* path) {

    // First check if the file exists
    if (!file_exists(path).unwrap()) {
        return Result(std::string(""), optional<Err>(Err("File does not exist")));
    }

    // Open the file
    std::ifstream file;
    file.open(*path);

    if (!file.is_open()) {
        return Result<std::string>("", optional<Err>(Err("Failed to open file")));
    }

    std::string content;
    std::string line;

    // Read line by line
    while (std::getline(file, line)) {
        content += line + "\n";
    }

    // Don't forget to close the file!
    file.close();

    return Result<std::string>(content, optional<Err>());

}

Result<bool> delete_file(const std::string* path) {
    
    // First check if the file exists
    if (!file_exists(path).unwrap()) {
        return Result(false, optional<Err>(Err("File does not exist")));
    }

    // Remove the file
    if (remove(path->c_str()) != 0) {
        return Result(false, optional<Err>(Err("Failed to delete file")));
    }

    return Result(true, optional<Err>());

}

Result<std::vector<std::string>> walk_dir(const std::string* path) {

    // First check if the file exists
    if (!file_exists(path).unwrap()) {
        return Result(std::vector<std::string>(), optional<Err>(Err("File does not exist")));
    }

    // Check if it's a directory
    if (!is_dir(path).unwrap()) {
        return Result(
            std::vector<std::string>(),
            optional<Err>(Err("Path is not a directory"))
        );
    }

    // Open the directory
    DIR* dir = opendir(path->c_str());

    if (dir == NULL) {
        return Result(
            std::vector<std::string>(),
            optional<Err>(Err("Failed to open directory"))
        );
    }

    std::vector<std::string> files;

    // Read the directory
    struct dirent* entry;
    while ((entry = readdir(dir)) != NULL) {
        files.push_back(entry->d_name);
    }

    // Close the directory
    closedir(dir);

    return Result(files, optional<Err>());

}

Result<bool> file_exists(const std::string* path) {

    // Check if the file exists
    if (std::filesystem::exists(*path)) {
        return Result(true, optional<Err>());
    }

    return Result(false, optional<Err>());

}

Result<bool> is_dir(const std::string* path) {

    // First check if the file exists
    if (!file_exists(path).unwrap()) {
        return Result(false, optional<Err>(Err("File does not exist")));
    }

    // Check if the path is a directory
    if (std::filesystem::is_directory(*path)) {
        return Result(true, optional<Err>());
    }

    return Result(false, optional<Err>());

}

Result<bool> delete_dir(const std::string* path) {

    // Remember that rmdir from unistd.h only works for empty directories
    // To fix that while still maintaning performance, we can use a queue
    // which has a complexity of O(n) instead of O(n^2) from the recursive approach

    std::vector<std::string> queue = { *path };

    while (!queue.empty()) {
        // Get always the first element
        std::string* current_path = &queue[0];
        Result<bool> is_dir_result = is_dir(current_path);

        // Cannot read the directory -> Return the error without crashing
        if (is_dir_result.has_error()) {
            return Result(false, optional<Err>(is_dir_result.get_error()));
        }

        if (is_dir_result.unwrap()) {
            // See if the directory is empty
            Result<std::vector<std::string>> files = walk_dir(current_path);

            // Cannot read the directory -> Return the error without crashing
            if (files.has_error()) {
                return Result(false, optional<Err>(files.get_error()));
            }

            std::vector<std::string>* files_vec = files.unwrap();

            // If the directory is empty, remove it
            if (files_vec->size() == 2) {
                if (rmdir(current_path->c_str()) != 0) {
                    return Result(false, optional<Err>(Err("Failed to delete directory")));
                }
            } else {
                // Add all the files to the queue
                for (const std::string& file : *files_vec) {
                    if (file != "." && file != "..") {
                        queue.push_back(*current_path + "/" + file);
                    }
                }

                // Add the directory to the end of the queue
                // After all files are deleted, the directory should be now empty
                queue.push_back(*current_path);
            }
        } else {
            // Try to remove the file
            Result<bool> delete_result = delete_file(current_path);

            // Cannot delete the file -> Return the error without crashing
            if (delete_result.has_error()) {
                return Result(false, optional<Err>(delete_result.get_error()));
            }
        }

        // Pop the first element
        queue.erase(queue.begin());
    }

    return Result(true, optional<Err>());
}