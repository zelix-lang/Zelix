#ifndef SYSTEM_ENV_H
#define SYSTEM_ENV_H

#include <string>
#include "../lang/result.h"

Result<std::string> get_env(const std::string* key);
Result<bool> set_env(const std::string* key, const std::string* value);

#endif