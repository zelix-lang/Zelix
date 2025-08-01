/*
        ==== The Fluent Programming Language ====
---------------------------------------------------------
  - This file is part of the Fluent Programming Language
    codebase. Fluent is a fast, statically-typed and
    memory-safe programming language that aims to
    match native speeds while staying highly performant.
---------------------------------------------------------
  - Fluent is categorized as free software; you can
    redistribute it and/or modify it under the terms of
    the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.
---------------------------------------------------------
  - Fluent is distributed in the hope that it will
    be useful, but WITHOUT ANY WARRANTY; without even
    the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE. See the GNU General Public
    License for more details.
---------------------------------------------------------
  - You should have received a copy of the GNU General
    Public License along with Fluent. If not, see
    <https://www.gnu.org/licenses/>.
*/

#include "lexer.h"

#include "ankerl/unordered_dense.h"

using namespace fluent;
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

bool push_token(container::vector<lexer::token> &tokens, const char *source)
{
    if (t_len == 0)
    {
        return true; // No token to push
    }

    const char *ptr = source + start;
    auto value_opt = container::optional<container::external_string>::emplace(ptr, t_len);
    const auto &value = value_opt.get();

    // Check if we have numbers or decimals
    if (num)
    {
        if (dec)
        {
            tokens.push_back(lexer::token{
                .value = value_opt,
                .type = lexer::token::DECIMAL_LITERAL,
                .line = line,
                .column = col - t_len
            });
        }
        else
        {
            tokens.push_back(lexer::token{
                .value = value_opt,
                .type = lexer::token::NUMBER_LITERAL,
                .line = line,
                .column = col - t_len
            });
        }
    }
    else if (str)
    {
        tokens.push_back(lexer::token{
            .value = value_opt,
            .type = lexer::token::STRING_LITERAL,
            .line = line,
            .column = col - t_len
        });
    }
    else
    {
        // Check if the token is known
        if (const auto it = token_map.find(value);
            it != token_map.end()
        )
        {
            tokens.push_back(lexer::token{
                .value = container::optional<container::external_string>::none(),
                .type = it->second, // Default to UNKNOWN for now
                .line = line,
                .column = col - t_len
            });
        }
        else
        {
            if (identifier)
            {
                tokens.push_back(lexer::token{
                    .value = value_opt,
                    .type = lexer::token::IDENTIFIER,
                    .line = line,
                    .column = col - t_len
                });
            }

            lexer::global_err.type = lexer::UNKNOWN_TOKEN;
            lexer::global_err.line = line;
            lexer::global_err.column = col - t_len;
            return false;
        }
    }

    // Reset the token length and flags
    reset_flags();

    return true;
}

container::optional<container::stream<lexer::token>> lexer::lex(
    const container::external_string &source
)
noexcept {
    container::vector<token> tokens; // Vector to hold the tokens
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

            push_token(tokens, ptr); // Push the current token if any
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
                return container::optional<container::stream<token>>::none();
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
                tokens.push_back(token{
                    .value = container::optional<container::external_string>::emplace(ptr + start, t_len),
                    .type = token::STRING_LITERAL,
                    .line = line,
                    .column = col - t_len
                });

                start = i + 1; // Move start to the next character
                t_len = 0; // Reset token length
                str = false; // Exit string mode
            }
            else
            {
                push_token(tokens, ptr);
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
        if (!str && c == '/' && ptr[i + 1] == '*')
        {
            // Push the current token if any
            push_token(tokens, ptr);
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
            !str && (c == '&' || c == '=' || c == '|')
            && ptr[i + 1] == c
        )
        {
            push_token(tokens, ptr);
            i++; // Skip the next character

            tokens.push_back(token{
                .value = container::optional<container::external_string>::none(),
                .type = c == '&' ? token::AND :
                        c == '|' ? token::OR :
                        token::BOOL_EQ,
                .line = line,
                .column = col
            });

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "&&"
            continue;
        }

        // Special cases: >=, <=, !=
        if (
            !str && (c == '>' || c == '<' || c == '!')
            && ptr[i + 1] == '='
        )
        {
            push_token(tokens, ptr);
            i++; // Skip the next character

            tokens.push_back(token{
                .value = container::optional<container::external_string>::none(),
                .type = c == '>' ? token::BOOL_GTE :
                        c == '<' ? token::BOOL_LTE :
                        token::BOOL_NEQ,
                .line = line,
                .column = col
            });

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "!=" or ">="
            continue;
        }

        // Special case: ->
        if (
            !str && c == '-' && ptr[i + 1] == '>'
        )
        {
            push_token(tokens, ptr);
            i++; // Skip the next character

            tokens.push_back(token{
                .value = container::optional<container::external_string>::none(),
                .type = token::ARROW,
                .line = line,
                .column = col
            });

            start = i + 1; // Move start to the next character
            reset_flags();
            col += 2; // Increment column for the "->"
            continue;
        }

        // Handle sing-char punctuation
        if (
            !str && c == '{' || c == '}' || c == '(' || c == ')'
            || c == '[' || c == ']' || c == ';' || c == ','
            || c == ':' || c == '=' || c == '+' || c == '-'
            || c == '*' || c == '/' || c == '!'
        )
        {
            push_token(tokens, ptr); // Push the current token if any

            tokens.push_back(token{
                .value = container::optional<container::external_string>::none(),
                .type = c == '{' ? token::OPEN_CURLY :
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
                            token::NOT,
                .line = line,
                .column = col
            });

            start = i + 1; // Move start to the next character
            reset_flags();
            col++; // Increment column for the punctuation
            continue;
        }

        // Handle decimals
        if (!str && c == '.')
        {
            // Prevent stuff like "1..2" or "1.2.3"
            if (dec)
            {
                global_err.type = UNKNOWN_TOKEN;
                global_err.line = line;
                global_err.column = col;
                return container::optional<container::stream<token>>::none();
            }

            dec = num; // If we were in a number, we are now in a decimal

            // Handle punctuation
            if (!dec)
            {
                push_token(tokens, ptr);
                col++; // Increment column for the dot

                tokens.push_back(token{
                    .value = container::optional<container::external_string>::none(),
                    .type = token::DOT,
                    .line = line,
                    .column = col
                });

                start = i + 1; // Move start to the next character
                reset_flags();
                continue;
            }
        }

        // Handle invalid identifiers
        if (!str && identifier && !isalnum(c) && c != '_')
        {
            global_err.type = UNKNOWN_TOKEN;
            global_err.line = line;
            global_err.column = col;
            return container::optional<container::stream<token>>::none();
        }

        if (!str && ((num && !isdigit(c) && c != '.') || (dec && !isdigit(c))))
        {
            global_err.type = UNKNOWN_TOKEN;
            global_err.line = line;
            global_err.column = col;
            return container::optional<container::stream<token>>::none();
        }

        col++; // Increment column for each character
        t_len++; // Increment token length
    }

    return container::optional<container::stream<token>>::emplace(tokens);
}