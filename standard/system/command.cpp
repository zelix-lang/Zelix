#include <cstdlib>
#include <iostream>
#include "../lang/result.h"

// Executes a shell command and returns the output
void exec(const char* cmd) {
    system(cmd);
}