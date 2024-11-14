import TokenType from "./TokenType.ts";

export default class DataTypes {
    public static list: TokenType[] = [
        TokenType.Num,
        TokenType.String,
        TokenType.Bool,
        TokenType.Nothing,
        TokenType.StringArray,
        TokenType.NumArray,
        TokenType.BoolArray,
        TokenType.Discrete,
    ];
}
