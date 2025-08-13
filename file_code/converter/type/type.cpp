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

//
// Created by rodri on 8/12/25.
//

#include "type.h"

using namespace zelix;

struct queue_el
{
    code::type &type;
    const parser::ast *node;
};

code::type code::converter::type(const parser::ast *root)
{
    // Create the root type
    code::type ret{};
    container::vector<queue_el> queue;
    queue.emplace_back(ret, root);

    while (!queue.empty())
    {
        // Get the last element in the queue
        const auto &el = queue.back();
        const auto [type, node] = el;
        queue.pop_back();

        // Parse the type
        size_t i = 0;
        const auto &children = node->children;
        const auto max = children.size();

        // Parse pointers first
        while (i < max && children[i]->rule == parser::ast::PTR)
        {
            type.pointers++;
            i++;
        }

        // Parse the base type
        const auto &base_node = children[i];
        switch (base_node->rule)
        {
            case parser::ast::STR:
            {
                type.base = type::STR;
                break;
            }

            case parser::ast::NOTHING:
            {
                type.base = type::NOTHING;
                break;
            }

            case parser::ast::NUM:
            {
                type.base = type::NUM;
                break;
            }

            case parser::ast::DEC:
            {
                type.base = type::DEC;
                break;
            }

            case parser::ast::BOOL:
            {
                type.base = type::BOOL;
                break;
            }

            default:
            {
                type.base = type::USER_DEFINED;
                type.name = base_node->value.get(); // Set the name of the type
                break;
            }
        }

        // Parse nested types
        for (
            const auto &nested_children = base_node->children;
            parser::ast *ast : nested_children
        )
        {
            type.children.emplace_back();
            queue.emplace_back(type.children.back(), ast);
        }
    }

    return ret;
}
