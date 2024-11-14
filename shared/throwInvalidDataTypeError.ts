import Token from "../lexer/defined/Token.ts";
import { catchException } from "./catchException.ts";

export function throwInvalidDataTypeError(trace: Token) {
    catchException(
        "Invalid data type.",
        ["The data type is not valid."],
        trace.buildTrace()
    )
}