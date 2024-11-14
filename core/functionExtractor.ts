import DataTypes from "../lexer/defined/DataTypes.ts";
import Token from "../lexer/defined/Token.ts";
import TokenType from "../lexer/defined/TokenType.ts";
import { catchException } from "../shared/catchException.ts";
import * as _Function from "../shared/object/Function.ts";
import Parameter from "../shared/object/Parameter.ts";
import { throwInvalidDataTypeError } from "../shared/throwInvalidDataTypeError.ts";
import extractImport from "./importExtractor.ts";
import Import from "./object/Import.ts";

const NAME_REGEX = /^[a-zA-Z_][a-zA-Z0-9_]*$/;

export default function extractFunction(tokens: Token[]): [
    Map<string, _Function.default>,
    Import[],
    Map<string, _Function.default>
] {
    const functions: Map<string, _Function.default> = new Map();
    const declaredFunctions = new Map();

    let inFunction: boolean = false;
    let nestedBrackets: number = 0;

    let expectingOpenParen: boolean = false;
    let expectingArgs: boolean = false;
    let expectingFunName: boolean = false;
    let expectingColon: boolean = false;
    let expectingArgType: boolean = false;
    let expectingComma: boolean = false;
    let expectingDataType: boolean = false;
    let expectingArrow: boolean = false;
    let expectingOpenCurly: boolean = false;
    let hasFunctionEnded: boolean = false;
    let hasFunctionReturned: boolean = false;
    let expectingFunKeyword: boolean = true;
    let isCurrentFunPublic: boolean = false;
    let lastFunctionName: string = "";
    let lastParameterName: string = "";
    let lastFunctionReturnedNative: string = "";
    let lastFunctionArgs: Parameter[] = [];
    let lastFunctionTokens: Token[] = [];

    let lastFunctionReturnType: TokenType = TokenType.Nothing;
    const imports : Import[] = [];

    for (let i = 0; i < tokens.length; i++) {
        const token: Token = tokens[i];
        const type: TokenType = token.getType();
        const value: string = token.getValue();

        if (type === TokenType.Function && !inFunction) {
            inFunction = true;
            expectingFunKeyword = false;
            expectingFunName = true;
        } else if (type !== TokenType.Function && !inFunction) {
            if (type === TokenType.Pub) {
                isCurrentFunPublic = true;
                expectingFunKeyword = true;
                continue;
            }

            if (type === TokenType.Import) {
                const [imported, skipped] = extractImport(tokens.slice(i + 1));
                i += skipped;
                
                imports.push(...imported);

                continue;
            }

            catchException(
                "You can't have anything outside a function definition!",
                [],
                token.buildTrace()
            );
        } else if (expectingOpenCurly) {
            if (type !== TokenType.OpenCurly) {
                catchException(
                    "Expected an open curly brace",
                    ["You need to open the function body with a curly brace."],
                    token.buildTrace()
                );
            }

            expectingOpenCurly = false;
        } else if (expectingDataType) {
            if (type === TokenType.Unknown) {
                if (!imports.some(a => a.getAlias() === value)) {
                    catchException(
                        "Unknown type",
                        ["You need to import the data type before using it."],
                        token.buildTrace()
                    );
                }

                lastFunctionReturnedNative = value;
                continue;
            }

            if (!DataTypes.list.includes(type))
                throwInvalidDataTypeError(token);

            if (type !== TokenType.Nothing && lastFunctionName === "main") {
                catchException(
                    "The main function can't return a value!",
                    ["Use nothing as the return type."],
                    token.buildTrace()
                );
            }

            lastFunctionReturnType = type;
            expectingDataType = false;
            expectingOpenCurly = true;
        } else if (expectingArrow) {
            if (type !== TokenType.Arrow) {
                catchException(
                    "Expected an arrow",
                    ["You need to separate the function arguments from the function's return type."],
                    token.buildTrace()
                )
            }

            if (lastFunctionArgs.length > 0 && lastFunctionName === "main") {
                catchException(
                    "The main function can't have arguments!",
                    [],
                    token.buildTrace()
                );
            }

            expectingArrow = false;
            expectingDataType = true;
        } else if (expectingComma) {
            if (type === TokenType.CloseParen) {
                expectingArgs = false;
                expectingArrow = true;
                expectingComma = false;

                continue;
            }

            if (type !== TokenType.Comma) {
                catchException(
                    "Expected a comma",
                    ["You need to separate the arguments with a comma."],
                    token.buildTrace()
                );
            }

            expectingComma = false;
            expectingArgs = true;
        } else if (expectingArgType) {
            if (type === TokenType.Unknown) {
                if (!imports.some(a => a.getAlias() === value)) {
                    catchException(
                        "Unknown type",
                        ["You need to import the data type before using it."],
                        token.buildTrace()
                    );
                }

                lastFunctionArgs.push(new Parameter(lastParameterName, type, value));

                expectingArgType = false;
                expectingComma = true;
                lastParameterName = "";
                continue;
            } else if (!DataTypes.list.includes(type)) throwInvalidDataTypeError(token);

            if (type === TokenType.Nothing) {
                catchException(
                    "You can't use the 'nothing' type as an argument type.",
                    [],
                    token.buildTrace()
                );
            }

            lastFunctionArgs.push(new Parameter(lastParameterName, type, ""));

            expectingArgType = false;
            expectingComma = true;
            lastParameterName = "";
        } else if (expectingColon) {
            if (type !== TokenType.Colon) {
                catchException(
                    "Expected a colon",
                    ["You need to separate the argument name from its type."],
                    token.buildTrace()
                );
            }

            expectingColon = false;
            expectingArgType = true;
        } else if (expectingArgs) {
            if (type === TokenType.CloseParen) {
                expectingArgs = false;
                expectingArrow = true;

                continue;
            }

            if (type !== TokenType.Unknown) {
                catchException(
                    "Expected an argument name",
                    ["You need to specify the name of the argument."],
                    token.buildTrace()
                );
            }

            lastParameterName = value;
            expectingColon = true;
            expectingArgs = false;
        } else if (expectingOpenParen) {
            if (type !== TokenType.OpenParen) {
                catchException(
                    "Expected an open parenthesis",
                    ["You need to open the function parameters with a parenthesis."],
                    token.buildTrace()
                )
            }

            expectingOpenParen = false;
            expectingArgs = true;
        } else if (expectingFunName) {
            if (type !== TokenType.Unknown) {
                catchException(
                    "Expected a function name",
                    ["You need to specify the name of the function."],
                    token.buildTrace()
                );
            }
            
            if (functions.has(value)) {
                catchException(
                    "Function already defined",
                    ["You can't define a function with the same name as another."],
                    token.buildTrace()
                );
            }

            if (!NAME_REGEX.test(value)) {
                catchException(
                    "Invalid function name",
                    ["Use a valid function name."],
                    token.buildTrace()
                );
            }

            lastFunctionName = value;
            expectingFunName = false;
            expectingOpenParen = true;
        } else if (expectingFunKeyword) {
            if (type !== TokenType.Function) {
                catchException(
                    "Expected the keyword 'fun'",
                    ["You need to start the function definition with the 'fun' keyword."],
                    token.buildTrace()
                );
            }

            expectingFunKeyword = false;
            inFunction = true;
            expectingFunName = true;
            hasFunctionEnded = false;
        } else {
            if (type === TokenType.CloseCurly) {
                if (nestedBrackets === 0) {
                    inFunction = false;
                    hasFunctionEnded = true;

                    if (!hasFunctionReturned && lastFunctionReturnType !== TokenType.Nothing) {
                        catchException(
                            "Function doesn't return a value",
                            ["Because you specified a return type, you need to return a value of that type."],
                            token.buildTrace()
                        );
                    }

                    const fun = new _Function.default(
                        lastFunctionTokens,
                        lastFunctionArgs,
                        lastFunctionReturnType,
                        lastFunctionReturnedNative
                    );

                    functions.set(
                        lastFunctionName,
                        fun    
                    );

                    if (isCurrentFunPublic) {
                        declaredFunctions.set(lastFunctionName, fun);
                    }

                    lastFunctionName = "";
                    lastFunctionArgs = [];
                    lastFunctionTokens = [];

                    hasFunctionReturned = false;
                    lastFunctionReturnType = TokenType.Nothing;
                    expectingArgType = false;
                    expectingArrow = false;
                    expectingColon = false;
                    expectingComma = false;
                    expectingDataType = false;
                    expectingOpenCurly = false;
                    expectingOpenParen = false;
                    expectingArgs = false;
                    expectingFunName = false;
                    isCurrentFunPublic = false;
                    hasFunctionEnded = true;

                    continue;
                }

                nestedBrackets--;
            } else if (type === TokenType.OpenCurly) nestedBrackets++;
            else if (type === TokenType.Return) {
                if (lastFunctionReturnType === TokenType.Nothing) {
                    catchException(
                        "Function doesn't return a value",
                        ["You can't return a value from a function that returns nothing."],
                        token.buildTrace()
                    );
                }

                hasFunctionReturned = true;
            }

            lastFunctionTokens.push(token);
        }
    }

    if (!hasFunctionEnded) {
        catchException(
            "You need to end the function with a curly brace",
            ["You need to close the function body with a curly brace."],
            tokens[tokens.length - 1].buildTrace()
        );
    }

    if (expectingFunKeyword) {
        catchException(
            "Expected the keyword 'fun'",
            ["You need to start the function definition with the 'fun' keyword."],
            tokens[tokens.length - 1].buildTrace()
        );
    }

    return [functions, imports, declaredFunctions];
}
