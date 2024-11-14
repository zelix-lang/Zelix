enum TokenType {
    Function, // Implemented
    NativeFunction, // Implemented
    NativeObject, // Implemented

    Let, // Implemented
    Const, // Implemented
    If, // Implemented
    Else, // Implemented
    ElseIf, // Implemented

    Return, // Implemented
    Assign, // Implemented
    AssignAdd, // Implemented
    AssignSub, // Implemented
    AssignSlash, // Implemented
    AssignAsterisk, // Implemented

    Plus, // Implementation not needed, handled by Engine
    Minus, // Implementation not needed, handled by Engine
    Increment, // Implemented
    Decrement, // Implemented
    Unsafe, // Implemented

    Asterisk, // Implementation not needed, handled by Engine
    Slash, // Implementation not needed, handled by Engine
    LessThan, // Implementation not needed, handled by Engine
    GreaterThan, // Implementation not needed, handled by Engine
    Equal, // Implementation not needed, handled by Engine
    NotEqual, // Implementation not needed, handled by Engine
    GreaterThanOrEqual, // Implementation not needed, handled by Engine
    LessThanOrEqual, // Implementation not needed, handled by Engine
    Percent,

    Arrow, // Implemented
    Comma, // Implemented

    Semicolon, // Implemented

    OpenParen, // Implementation not needed, handled by Engine
    CloseParen, // Implementation not needed, handled by Engine

    OpenCurly, // Implementation not needed, handled by Engine
    CloseCurly, // Implementation not needed, handled by Engine

    Colon, // Implemented
    Xor, // Implementation not needed, handled by Engine
    Not, // Implementation not needed, handled by Engine

    OpenBracket, // Implemented
    CloseBracket, // Implemented

    Discrete, // Implemented
    Dot, // Implementation not needed, handled by Engine

    StringArray, // Implemented
    NumArray, // Implemented
    BoolArray, // Implemented

    StringArrayLiteral, // Implementation not needed, handled by Engine
    NumArrayLiteral, // Implementation not needed, handled by Engine
    BoolArrayLiteral, // Implementation not needed, handled by Engine
    DiscreteLiteral, // Implementation not needed, handled by Engine

    String, // Implemented
    Num, // Implemented
    Nothing, // Implemented
    Bool, // Implemented

    StringLiteral, // Implemented
    NumLiteral, // Implemented
    BoolLiteral, // Implemented

    While, // Implemented
    For, // Implemented
    In, // Implemented

    Break, // Implemented
    Continue, // Implemented
    Pub, // Implemented

    Import, // Implemented
    From, // Implemented
    As, // Implemented

    Ampersand, // Implementation not needed, handled by Engine
    Bar, // Implementation not needed, handled by Engine

    Unknown, // Implemented
}

export function name(tokenType: TokenType): string {
    return TokenType[tokenType];
}

export default TokenType;
