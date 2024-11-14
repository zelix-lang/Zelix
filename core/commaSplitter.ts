import Token from "../lexer/defined/Token.ts";
import TokenType from "../lexer/defined/TokenType.ts";

export default function commaSplitter(tokens: Token[]) : Token[][] {
    const builder = [];
    const result = [];

    for (let i = 0; i < tokens.length; i++) {
        const token = tokens[i];

        if (token.getType() === TokenType.Comma) {
            result.push([...builder]);
            builder.length = 0;

            continue;
        }

        builder.push(token);
    }

    if (builder.length > 0) {
        result.push(builder);
    }

    return result;
}