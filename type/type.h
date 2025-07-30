/*
    The Fluent Programming Language
    -----------------------------------------------------
    This code is released under the GNU GPL v3 license.
    For more information, please visit:
    https://www.gnu.org/licenses/gpl-3.0.html
    -----------------------------------------------------
    Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
    This program comes with ABSOLUTELY NO WARRANTY.
    For details type `fluent l`. This is free software,
    and you are welcome to redistribute it under certain
    conditions; type `fluent l -f` for details.
*/

//
// Created by rodrigo on 6/9/25.
//

#ifndef FLUENT_TYPE_H
#define FLUENT_TYPE_H

// ============= FLUENT LIB C =============
#include <fluent/types/types.h>
#include <fluent/std_bool/std_bool.h>
#include <fluent/vector/vector.h>

typedef enum
{
    TYPE_NUM,        ///< Numeric type
    TYPE_DEC,        ///< Decimal type
    TYPE_BOOL,       ///< Boolean type
    TYPE_STRING,     ///< String type
    TYPE_NOTHING,    ///< Represents no type (void)
    TYPE_CUSTOM      ///< Custom type (user-defined)
} type_enum_t;

typedef struct
{
    size_t pointers;   ///< Number of pointers in the type
    size_t arrays;     ///< Number of array dimensions in the type
    bool is_primitive; ///< Whether the type is a primitive type
    type_enum_t type;  ///< The type enumeration (e.g., TYPE_NUM, TYPE_STRING)
    char *name;        ///< Name of the custom type (if applicable)
    vector_generic_t *children; ///< Vector of child types (for complex types)
} type_t;

#endif //FLUENT_TYPE_H
