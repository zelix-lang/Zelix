import TokenType from "./defined/TokenType.ts";
import Token from "./defined/Token.ts";
import DataTypes from "./defined/DataTypes.ts";
import { catchException } from "../shared/catchException.ts";

const types = new Map<string, TokenType>([
    ["return", TokenType.Return],
    ["fun", TokenType.Function],
    ["let", TokenType.Let],
    ["const", TokenType.Const],
    ["while", TokenType.While],
    ["for", TokenType.For],
    ["break", TokenType.Break],
    ["continue", TokenType.Continue],
    ["in", TokenType.In],
    ["if", TokenType.If],
    ["else", TokenType.Else],
    ["elseif", TokenType.ElseIf],
    ["unsafe", TokenType.Unsafe],
    ["=", TokenType.Assign],
    ["+", TokenType.Plus],
    ["-", TokenType.Minus],
    ["++", TokenType.Increment],
    ["--", TokenType.Decrement],
    ["*", TokenType.Asterisk],
    ["/", TokenType.Slash],
    ["<", TokenType.LessThan],
    [">", TokenType.GreaterThan],
    ["+=", TokenType.AssignAdd],
    ["-=", TokenType.AssignSub],
    ["*=", TokenType.AssignAsterisk],
    ["/=", TokenType.AssignSlash],
    ["==", TokenType.Equal],
    ["!=", TokenType.NotEqual],
    ["<=", TokenType.LessThanOrEqual],
    [">=", TokenType.GreaterThanOrEqual],
    ["&", TokenType.Ampersand],
    ["|", TokenType.Bar],
    ["^", TokenType.Xor],
    ["!", TokenType.Not],
    [",", TokenType.Comma],
    [";", TokenType.Semicolon],
    ["(", TokenType.OpenParen],
    [")", TokenType.CloseParen],
    ["{", TokenType.OpenCurly],
    ["}", TokenType.CloseCurly],
    [":", TokenType.Colon],
    ["->", TokenType.Arrow],
    ["[", TokenType.OpenBracket],
    ["]", TokenType.CloseBracket],
    [".", TokenType.Dot],
    ["%", TokenType.Percent],
    ["string", TokenType.String],
    ["num", TokenType.Num],
    ["nothing", TokenType.Nothing],
    ["bool", TokenType.Bool],
    ["pub", TokenType.Pub],
    ["import", TokenType.Import],
    ["from", TokenType.From],
    ["as", TokenType.As],
    ["string[]", TokenType.StringArray],
    ["num[]", TokenType.NumArray],
    ["bool[]", TokenType.BoolArray],
    ["[discrete]", TokenType.Discrete],
]);

export default class Lexer {
    private inString: boolean = false;
    private inEscape: boolean = false;
    private inComment: boolean = false;
    private inBlockComment: boolean = false;

    private static NumBER_REGEX: RegExp = /^\d+((\.\d+)?)$/;
    private static BoolEAN_REGEX: RegExp = /^(true|false)$/;
    private static PUNCTUATION_CHARACTERS = "[;,(){}:+->%<.=![]/|*^&]";

    calculate(
        string: string,
        file: string,
        line: number,
        column: number,
    ): Token {
        string = string.trim();

        const type = types.get(string) ?? (() => {
            if (Lexer.NumBER_REGEX.test(string)) return TokenType.NumLiteral;
            if (Lexer.BoolEAN_REGEX.test(string)) return TokenType.BoolLiteral;

            return TokenType.Unknown;
        })();

        return new Token(type, string, file, line, column);
    }

    private pushTokenIfBuilderNotEmpty(tokens: Token[], builder: string, file: string, line: number, column: number) {
        if (builder) {
            tokens.push(this.calculate(builder, file, line, column));
            return "";
        }
        
        return builder;
    }


    tokenize(file: string, code: string): Token[] {
        const tokens: Token[] = [];
        let builder: string = "";
        const characters: string[] = code.split("");
        let line: number = 1;
        let column: number = 1;

        for (let i = 0; i < characters.length; i++) {
            const character = characters[i];
            if (character === "\n") {
                line++;
                column = 1;
            } else {
                column++;
            }

            if (!this.inBlockComment && character === "/") {
                if (characters[i + 1] === "*") {
                    this.inBlockComment = true;
                    
                    continue;
                }
            } else if (this.inBlockComment) {
                if (character === "/" && characters[i - 1] === "*") {
                    this.inBlockComment = false;
                }

                continue;
            }
            
            if (this.inComment) {
                if (character === "\n") {
                    this.inComment = false;
                }
                
                continue;
            } else if (this.inString) {
                if (character === '"' && !this.inEscape) {
                    this.inString = false;
                    tokens.push(
                        new Token(
                            TokenType.StringLiteral,
                            builder,
                            file,
                            line,
                            column,
                        ),
                    );
                    builder = "";
                } else if (character === "\\" && !this.inEscape) {
                    this.inEscape = true;
                    builder += character;
                } else {
                    builder += character;
                    this.inEscape = character === "\\";
                }
                continue;
            } else if (character === '"') {
                this.inString = true;
                continue;
            } else if (character === "/" && characters[i + 1] === "/") {
                this.inComment = true;
                builder = "";
                continue;
            } else if (Lexer.PUNCTUATION_CHARACTERS.includes(character)) {
                if (character === "." && !isNaN(parseInt(builder))) {
                    builder += character;
                    continue;
                }

                if (builder) {
                    builder = this.pushTokenIfBuilderNotEmpty(tokens, builder, file, line, column);
                }

                if (character === "]") {
                    const lastToken = tokens[tokens.length - 1];
                    const lastLastToken = tokens[tokens.length - 2];

                    if (
                        DataTypes.list.includes(lastLastToken?.getType()) &&
                        lastToken?.getType() === TokenType.OpenBracket
                        // num[], string[], bool[]
                    ) {
                        tokens.pop();
                        tokens.pop();

                        tokens.push(this.calculate(lastLastToken.getValue()
                            +lastToken.getValue()
                            +character, file, line, column));
                        continue;
                    } else if (
                        lastLastToken?.getType() === TokenType.OpenBracket &&
                        lastToken.getValue() === "discrete"
                    ) {
                        tokens.pop();
                        tokens.pop();
                        
                        tokens.push(this.calculate(lastLastToken.getValue()
                            +lastToken.getValue()
                            +character, file, line, column));

                        continue;
                    }
                } else if (
                    character === "=" ||
                    character === "-" ||
                    character === "+"
                ) {
                    const lastToken = tokens[tokens.length - 1];
                    
                    if (
                        (
                            character === "=" &&
                            (
                                lastToken?.getType() === TokenType.Assign
                                || lastToken?.getType() === TokenType.LessThan
                                || lastToken?.getType() === TokenType.GreaterThan
                                || lastToken?.getType() === TokenType.Not
                                || lastToken?.getType() === TokenType.Plus
                                || lastToken?.getType() === TokenType.Minus
                                || lastToken?.getType() === TokenType.Asterisk
                                || lastToken?.getType() === TokenType.Slash
                            )
                        )
                        || (
                            (character === "-" && lastToken?.getType() === TokenType.Minus) ||
                            (character === "+" && lastToken?.getType() === TokenType.Plus)
                        )
                    ) {
                        tokens.pop();
                        tokens.push(this.calculate(lastToken.getValue() + character, file, line, column));
                        continue;
                    }
                } else if (character === ">" && tokens[tokens.length - 1]?.getType() === TokenType.Minus) {
                    tokens.pop();
                    tokens.push(this.calculate("->", file, line, column));
                    continue;
                }

                tokens.push(this.calculate(character, file, line, column));
            } else if (!character.trim()) {
                builder = this.pushTokenIfBuilderNotEmpty(tokens, builder, file, line, column);
            } else {
                builder += character;
            }
        }

        if (this.inString) {
            catchException(
                "Unterminated string!",
                ["You need to close the string with a double quote"],
                tokens[tokens.length - 1].buildTrace()
            );
        }
        if (builder) this.pushTokenIfBuilderNotEmpty(tokens, builder, file, line, column);

        return tokens;
    }
}
