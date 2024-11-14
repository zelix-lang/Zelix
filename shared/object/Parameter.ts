import TokenType from "../../lexer/defined/TokenType.ts";

export default class Parameter {

    private readonly name: string;
    private readonly nativeType: string;
    private readonly type: TokenType;

    constructor(name: string, type: TokenType, nativeType: string) {
        this.name = name;
        this.type = type;
        this.nativeType = nativeType;
    }

    public getName(): string {
        return this.name;
    }

    public getType(): TokenType {
        return this.type;
    }

    public getNativeType(): string {
        return this.nativeType;
    }

}