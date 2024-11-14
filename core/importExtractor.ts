import Token from "../lexer/defined/Token.ts";
import TokenType from "../lexer/defined/TokenType.ts";
import { catchException } from "../shared/catchException.ts";
import commaSplitter from "./commaSplitter.ts";
import Import from "./object/Import.ts";
import extractTokens from "./tokensExtractor.ts";

export default function extractImport(tokens: Token[]) : [Import[], number] {
    const [everything, skipped] = extractTokens(tokens, TokenType.Semicolon);
    const [beforeFrom, beforeFromSkipped] = extractTokens(everything, TokenType.From);
    const afterFrom = everything.slice(beforeFromSkipped + 1);

    if (afterFrom.length !== 1 || afterFrom[0].getType() !== TokenType.StringLiteral) {
        catchException(
            "Invalid import statement!",
            ["Expected a string literal after 'from' but found none."],
            tokens[beforeFromSkipped + 1].buildTrace()
        );
    }

    const path = afterFrom[0].getValue();

    if (beforeFrom.length === 1) {
        return [[new Import(
            path,
            beforeFrom[0].getValue(),
            beforeFrom[0].getValue(),
            // Add 1 because we skipped the "import" token
        )], skipped + 1];
    }

    const result : Import[] = [];
    const parts = commaSplitter(beforeFrom);
    
    for (const part of parts) {
        if (part.length !== 1 && part.length !== 3) {
            catchException(
                "Invalid import statement!",
                ["Use either 'import A from B' or 'import A as B, C as D from E'."],
                part[0].buildTrace(),
                "Expected either one or three tokens in an import statement."
            );
        }

        const [name, asToken, alias] = part;

        if (
            name.getType() !== TokenType.Unknown
            || (asToken && asToken.getType() !== TokenType.As)
            || (alias && alias.getType() !== TokenType.Unknown)
        ) {
            catchException(
                "Invalid import statement!",
                ["Use either 'import A from B' or 'import A as B, C as D from E'."],
                part[0].buildTrace(),
                "Expected either 'import A from B' or 'import A as B, C as D from E'."
            );
        }

        if (part.length === 1) {
            result.push(new Import(
                path,
                name.getValue(),
                name.getValue(),
            ));
        } else {
            result.push(new Import(
                path,
                alias.getValue(),
                name.getValue(),
            ));
        }
    }

    // Add 1 because we skipped the "import" token
    return [result, skipped + 1];
}