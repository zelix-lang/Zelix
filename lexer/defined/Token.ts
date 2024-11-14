import * as path from "@std/path";
import type TokenType from "./TokenType.ts";

export default class Token {
    private readonly type: TokenType;
    private readonly value: string;
    private readonly file: string;
    private readonly line: number;
    private readonly column: number;

    constructor(
        type: TokenType,
        value: string,
        file: string,
        line: number,
        column: number,
    ) {
        this.type = type;
        this.value = value;
        this.file = path.basename(file);
        this.line = line;
        this.column = column;
    }

    public getType(): TokenType {
        return this.type;
    }

    public getValue(): string {
        return this.value;
    }

    public buildHelpExample() {
        return "You can check the docs";
    }

    public buildTrace() : string {
        return "At " + this.file + ":" + this.line + ":" + this.column;
    }
}
