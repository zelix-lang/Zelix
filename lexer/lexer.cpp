/*
        ==== The Zelix Programming Language ====
---------------------------------------------------------
  - This file is part of the Zelix Programming Language
    codebase. Zelix is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Zelix is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Zelix is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Zelix. If not, see
    <https://www.gnu.org/licenses/>.
*/

#include "lexer.h"

#include "ankerl/unordered_dense.h"
#include "memory/allocator.h"

using namespace zelix;
ankerl::unordered_dense::map<
    container::external_string,
    lexer::token::t_type,
    container::external_string_hash
> token_map = {
    {
        container::external_string("import", 6),
        lexer::token::IMPORT
    },
    {
        container::external_string("fun", 3),
        lexer::token::FUNCTION
    },
    {
        container::external_string("mod", 3),
        lexer::token::MOD
    },
    {
        container::external_string("str", 3),
        lexer::token::STRING
    },
    {
        container::external_string("num", 3),
        lexer::token::NUMBER
    },
    {
        container::external_string("dec", 3),
        lexer::token::DECIMAL
    },
    {
        container::external_string("nothing", 7),
        lexer::token::NOTHING
    },
    {
        container::external_string("bool", 4),
        lexer::token::BOOL
    },
    {
        container::external_string("step", 4),
        lexer::token::STEP
    },
    {
        container::external_string("true", 4),
        lexer::token::TRUE
    },
    {
        container::external_string("false", 5),
        lexer::token::FALSE
    },
    {
        container::external_string("let", 3),
        lexer::token::LET
    },
    {
        container::external_string("const", 5),
        lexer::token::CONST
    },
    {
        container::external_string("pub", 3),
        lexer::token::PUB
    },
    {
        container::external_string("if", 2),
        lexer::token::IF
    },
    {
        container::external_string("else", 4),
        lexer::token::ELSE
    },
    {
        container::external_string("elseif", 6),
        lexer::token::ELSEIF
    },
    {
        container::external_string("for", 3),
        lexer::token::FOR
    },
    {
        container::external_string("while", 5),
        lexer::token::WHILE
    },
    {
        container::external_string("return", 6),
        lexer::token::RETURN
    },
    {
        container::external_string("to", 2),
        lexer::token::TO
    },
    {
        container::external_string("in", 2),
        lexer::token::IN
    },
};

bool lexer::is_err()
noexcept {
    return global_err.type != NONE;
}

inline size_t line = 1;
inline size_t col = 1;
inline size_t start = 0; // Start index for the current token
inline size_t t_len = 0; // Length of the current token
inline bool str = false; // Flag to track if we are in a string literal
inline bool identifier = false; // Flag to track if we have an identifier
inline bool num = false; // Flag to track if we are in a number literal
inline bool dec = false; // Flag to track if we are in a decimal literal
inline bool block_comment = false; // Flag to track if we are in a block comment

/// \brief Resets all tokenization flags and the current token length.
///
/// This function is called after a token is pushed to reset the state
/// for the next token. It clears the flags for string, identifier,
/// number, and decimal detection, and sets the token length to zero.
void reset_flags()
{
    str = false;
    identifier = false;
    num = false;
    dec = false;
    t_len = 0;
}

/// \brief Pushes the current token to the token vector if a valid token is present.
///
/// This function examines the current tokenization state (number, decimal, string, identifier)
/// and pushes the appropriate token to the provided token vector. It also checks for known
/// keywords using the token map. If the token is unknown and not an identifier, it sets a
/// global error and returns false.
///
/// \param tokens Reference to the vector where tokens are stored.
/// \param source Pointer to the source code string.
/// \param allocator Lazy allocator for managing memory for the tokens.
void push_token(
    container::stream<lexer::token *> &tokens,
    const char *source,
    memory::lazy_allocator<lexer::token> &allocator
)
{
    if (t_len == 0)
    {
        return; // No token to push
    }

    const char *ptr = source + start;
    auto value_opt = container::optional<container::external_string>::emplace(ptr, t_len);
    const auto &value = value_opt.get();

    // Check if we have numbers or decimals
    if (num)
    {
        if (dec)
        {
            const auto t = allocator.alloc();
            t->value = value_opt;
            t->type = lexer::token::DECIMAL_LITERAL;
            t->line = line;
            t->column = col - t_len;
            tokens.push(t);
        }
        else
        {
            const auto t = allocator.alloc();
            t->value = value_opt;
            t->type = lexer::token::NUMBER_LITERAL;
            t->line = line;
            t->column = col - t_len;
            tokens.push(t);
        }
    }
    else if (str)
    {
        const auto t = allocator.alloc();
        t->value = value_opt;
        t->type = lexer::token::STRING_LITERAL;
        t->line = line;
        t->column = col - t_len;
        tokens.push(t);
    }
    else
    {
        // Check if the token is known
        if (const auto it = token_map.find(value);
            it != token_map.end()
        )
        {
            const auto t = allocator.alloc();
            t->value = container::optional<container::external_string>::none();
            t->type = it->second;
            t->line = line;
            t->column = col - t_len;
            tokens.push(t);
        }
        else
        {
            if (__builtin_expect(identifier, true))
            {
                const auto t = allocator.alloc();
                t->value = value_opt;
                t->type = lexer::token::IDENTIFIER;
                t->line = line;
                t->column = col - t_len;
                tokens.push(t);
            }
            else
            {
                lexer::global_err.type = lexer::UNKNOWN_TOKEN;
                lexer::global_err.line = line;
                lexer::global_err.column = col - t_len;
                throw except::exception("Unknown token");
            }
        }
    }

    // Reset the token length and flags
    reset_flags();
}

/// \brief Lexical analyzer for the Fluent Programming Language.
///
/// This function tokenizes the given source code string into a stream of tokens.
/// It handles whitespace, newlines, string literals, line and block comments,
/// identifiers, numbers, decimals, and various punctuation and operator tokens.
///
/// The lexer maintains state for line and column tracking, as well as flags for
/// string, identifier, number, decimal, and block comment detection. It pushes
/// tokens to the output vector as they are recognized, and sets global error
/// information if an invalid or unknown token is encountered.
///
/// \param source The source code to tokenize, as an external string.
/// \param allocator A lazy allocator for managing memory for the tokens.
/// \return An optional stream of tokens. Returns none if a lexical error occurs.
///
/// \note This function is not thread-safe due to the use of global state for
///       error reporting and tokenization flags.
container::stream<lexer::token *> lexer::lex(
    const container::external_string &source,
    memory::lazy_allocator<token> &allocator
)
{
    container::vector<token *> vec; // Vector to hold the tokens
    container::stream<token *> tokens(vec); // Stream to return
    const auto ptr = source.ptr();

    for (size_t i = 0; i < source.size(); i++)
    {
        const char c = ptr[i];

        // Handle spaces
        if (!block_comment && c == ' ')
        {
            // If we are in a string, continue
            if (str)
            {
                t_len++; // Increment token length
                continue;
            }

           push_token(tokens, ptr, allocator);
            start = i + 1; // Move start to the next character
            continue;
        }

        // Handle newlines
        if (c == '\n')
        {
            if (block_comment)
            {
                line++;
                col = 1;
                continue;
            }

            if (str)
            {
                global_err.type = UNCLOSED_STRING;
                global_err.line = line;
                global_err.column = col;
                throw except::exception("Unclosed string literal");
            }

            line++;
            col = 1;
            t_len = 0; // Reset token length
            start = i + 1; // Move start to the next character
            continue; // Skip to the next character
        }

        // Handle string literals
        if (!block_comment && c == '"')
        {
            if (str)
            {
                // If we are already in a string, push the token
                auto t = allocator.alloc();
                t->value = container::optional<container::external_string>::emplace(ptr + start, t_len);
                t->type = token::STRING_LITERAL;
                t->line = line;
                t->column = col - t_len;
                tokens.push(t);

                start = i + 1; // Move start to the next character
                t_len = 0; // Reset token length
                str = false; // Exit string mode
            }
            else
            {
                push_token(tokens, ptr, allocator);
                // Start a new string literal
                str = true;
                start = i + 1; // Start after the quote
                t_len = 0; // Reset token length
            }

            col++;
            continue;
        }

        // Handle comments
        // Since ptr is null-terminated, we can safely check i + 1
        // since we'd get the null character if we go out of bounds
        if (!block_comment && !str && c == '/' && ptr[i + 1] == '/')
        {
            // Get the next newline
            if (
                const char *next_newline = strchr(ptr + i + 2, '\n');
                next_newline != nullptr
            )
            {
                const size_t comment_length = next_newline - (ptr + i + 2);
                i += comment_length + 1; // Move past the comment
                start = i; // Reset start to the next character
                reset_flags();
            }
            else
            {
                // If no newline is found, we are at the end of the file
                break;
            }

            continue; // Skip the rest of the line
        }

        if (block_comment && c == '*' && ptr[i + 1] == '/')
        {
            // End of block comment
            block_comment = false; // Exit block comment mode
            i++; // Skip the closing slash
            col++; // Increment column for the closing slash
            start = i + 1; // Move start to the next character
            t_len = 0;
            reset_flags();
            continue; // Skip the rest of the line
        }

        if (block_comment)
        {
            col++; // Increment column for each character in the block comment
            continue;
        }

        // Ignore characters if we have an open string
        if (str)
        {
            col++;
            t_len++; // Increment token length for the string
            continue;
        }

        // Handle block comments
        if (c == '/' && ptr[i + 1] == '*')
        {
            // Push the current token if any
            push_token(tokens, ptr, allocator);
            block_comment = true; // Enter block comment mode
            continue; // Skip the rest of the line
        }

        // Handle identifier detection
        if (t_len == 0)
        {
            if (isalpha(c) || c == '_')
            {
                identifier = true; // Start of an identifier
                start = i; // Set start to the current index
            }
            else if (isdigit(c))
            {
                num = true; // Start of a number
                start = i; // Set start to the current index
            }
            else if (c == '.')
            {
                dec = true; // Start of a decimal
                start = i; // Set start to the current index
            }
        }

        // Handle punctuation signs
        if (
            (c == '&' || c == '=' || c == '|')
            && ptr[i + 1] == c
        )
        {
            push_token(tokens, ptr, allocator);
            i++; // Skip the next character

            const auto t = allocator.alloc();
            t->value = container::optional<container::external_string>::none();
            t->line = line;
            t->column = col;
            t->type = c == '&' ? token::AND :
                        c == '|' ? token::OR :
                        token::BOOL_EQ,
            tokens.push(t);

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "&&"
            continue;
        }

        // Special cases: >=, <=, !=
        if (
            (c == '>' || c == '<' || c == '!')
            && ptr[i + 1] == '='
        )
        {
            push_token(tokens, ptr, allocator);
            i++; // Skip the next character

            const auto t = allocator.alloc();
            t->value = container::optional<container::external_string>::none();
            t->line = line;
            t->column = col;
            t->type = c == '>' ? token::BOOL_GTE :
                        c == '<' ? token::BOOL_LTE :
                        token::BOOL_NEQ;
            tokens.push(t);

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "!=" or ">="
            continue;
        }

        // Special case: ->
        if (
            c == '-' && ptr[i + 1] == '>'
        )
        {
            push_token(tokens, ptr, allocator);
            i++; // Skip the next character

            const auto t = allocator.alloc();
            t->value = container::optional<container::external_string>::none();
            t->type = token::ARROW;
            t->line = line;
            t->column = col;
            tokens.push(t);

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "->"
            continue;
        }

        // Handle sing-char punctuation
        if (
            c == '{' || c == '}' || c == '(' || c == ')'
            || c == '[' || c == ']' || c == ';' || c == ','
            || c == ':' || c == '=' || c == '+' || c == '-'
            || c == '*' || c == '/' || c == '!' || c == '&'
        )
        {
            push_token(tokens, ptr, allocator);

            const auto t = allocator.alloc();
            t->value = container::optional<container::external_string>::none();
            t->line = line;
            t->column = col;
            t->type = c == '{' ? token::OPEN_CURLY :
                        c == '}' ? token::CLOSE_CURLY :
                        c == '(' ? token::OPEN_PAREN :
                        c == ')' ? token::CLOSE_PAREN :
                        c == '[' ? token::OPEN_BRACKET :
                        c == ']' ? token::CLOSE_BRACKET :
                        c == ';' ? token::SEMICOLON :
                        c == ',' ? token::COMMA :
                        c == ':' ? token::COLON :
                        c == '=' ? token::EQUALS :
                        c == '+' ? token::PLUS :
                        c == '-' ? token::MINUS :
                        c == '*' ? token::MULTIPLY :
                        c == '/' ? token::DIVIDE :
                        c == '&' ? token::AMPERSAND :
                            token::NOT;

            tokens.push(t);

            start = i + 1; // Move start to the next character
            reset_flags();
            col++; // Increment column for the punctuation
            continue;
        }

        // Handle decimals
        if (c == '.')
        {
            // Prevent stuff like "1..2" or "1.2.3"
            if (dec)
            {
                global_err.type = UNKNOWN_TOKEN;
                global_err.line = line;
                global_err.column = col;
                throw except::exception("Unexpected decimal point in number");
            }

            dec = num; // If we were in a number, we are now in a decimal

            // Handle punctuation
            if (!dec)
            {
                push_token(tokens, ptr, allocator);
                col++; // Increment column for the dot

                const auto t = allocator.alloc();
                t->value = container::optional<container::external_string>::none();
                t->type = token::DOT;
                t->line = line;
                t->column = col - 1; // Column is the current position minus one for
                tokens.push(t);

                start = i + 1; // Move start to the next character
                reset_flags();
                continue;
            }
        }

        // Handle invalid identifiers
        if (identifier && !isalnum(c) && c != '_')
        {
            global_err.type = UNKNOWN_TOKEN;
            global_err.line = line;
            global_err.column = col;
            throw except::exception("Invalid character in identifier");
        }

        if (!str && ((num && !isdigit(c) && c != '.') || (dec && !isdigit(c))))
        {
            global_err.type = UNKNOWN_TOKEN;
            global_err.line = line;
            global_err.column = col;
            throw except::exception("Invalid character in number or decimal");
        }

        col++; // Increment column for each character
        t_len++; // Increment token length
    }

    // Handle unclosed comments or strings
    if (str || block_comment)
    {
        global_err.type = str ? UNCLOSED_STRING : UNCLOSED_COMMENT;
        global_err.line = line;
        global_err.column = col - t_len; // Column where the unclosed string/comment started
        return tokens;
    }

    return tokens;
}