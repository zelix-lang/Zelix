import Token from "../../lexer/defined/Token.ts";
import TokenType from "../../lexer/defined/TokenType.ts";
import Parameter from "./Parameter.ts";

export default class Function {

    private readonly body: Token[];
    private readonly parameters: Parameter[];
    private readonly returnType: TokenType;
    private readonly returnedNative: string;

    constructor(body: Token[], parameters: Parameter[], returnType: TokenType, returnedNative: string) {
        this.body = body;
        this.parameters = parameters;
        this.returnType = returnType;
        this.returnedNative = returnedNative;
    }

    getBody(): Token[] {
        return this.body;
    }

    getParameters(): Parameter[] {
        return this.parameters;
    }

    getReturnType(): TokenType {
        return this.returnType;
    }

    getReturnedNative(): string {
        return this.returnedNative;
    }
    
}