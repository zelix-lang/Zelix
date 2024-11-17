#[derive(Debug, Clone, PartialEq)]
pub enum TokenType {
    Function, 
    NativeFunction, 
    NativeObject, 

    Let, 
    Const, 
    If, 
    Else, 
    ElseIf, 

    Return, 
    Assign, 
    AssignAdd, 
    AssignSub, 
    AssignSlash, 
    AssignAsterisk, 

    Plus, 
    Minus, 
    Increment, 
    Decrement, 
    Unsafe, 

    Asterisk, 
    Slash, 
    LessThan, 
    GreaterThan, 
    Equal, 
    NotEqual, 
    GreaterThanOrEqual, 
    LessThanOrEqual, 
    Percent,

    Arrow, 
    Comma, 

    Semicolon, 

    OpenParen, 
    CloseParen, 

    OpenCurly, 
    CloseCurly, 

    Colon, 
    Xor, 
    Not, 

    OpenBracket, 
    CloseBracket, 

    Discrete, 
    Dot, 

    StringArray, 
    NumArray, 
    BoolArray, 

    StringArrayLiteral, 
    NumArrayLiteral, 
    BoolArrayLiteral, 
    DiscreteLiteral, 

    String, 
    Num, 
    Nothing, 
    Bool, 

    StringLiteral, 
    NumLiteral, 
    BoolLiteral, 

    While, 
    For, 
    In, 

    Break, 
    Continue, 
    Pub, 

    Import, 
    From, 
    As, 

    Ampersand, 
    Bar, 

    Unknown, 
}