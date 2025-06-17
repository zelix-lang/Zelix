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
// Created by rodrigo on 6/16/25.
//

#ifndef FLUENT_PARSER_QUEUE_EXPRESSION_H
#define FLUENT_PARSER_QUEUE_EXPRESSION_H

// ============= FLUENT LIB C =============
#include <fluent/alinked_queue/alinked_queue.h> // fluent_libc

typedef struct
{
    token_t **body;
    size_t start;
    size_t len;
} queue_expression_t;

DEFINE_ALINKED_NODE(queue_expression_t, expr);

#endif //FLUENT_PARSER_QUEUE_EXPRESSION_H
