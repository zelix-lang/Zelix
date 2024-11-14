import Token from "../lexer/defined/Token.ts";
import TokenType, { name } from "../lexer/defined/TokenType.ts";
import { catchException } from "../shared/catchException.ts";

export default function extractTokens(
    tokens: Token[],
    delimiter: TokenType
) : [Token[], number] {
    const builder = [];
    let hasMetDelimiter = false;
    let lastIdx : number = 0;

    for (let i = 0; i < tokens.length; i++) {
        const token = tokens[i];

        if (token.getType() === delimiter) {
            hasMetDelimiter = true;
            lastIdx = i;
            break;
        }

        builder.push(token);
    }

    if (!hasMetDelimiter || builder.length === 0) {
        catchException(
            "Invalid statement!",
            ["Expected a delimiter (" + name(delimiter) + ") but found none."],
            tokens[lastIdx]?.buildTrace() ?? "No trace available"
        );
    }

    return [builder, lastIdx];
}