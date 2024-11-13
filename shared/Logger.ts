import * as AnsiColors from "./ansi/AnsiColors.ts";

function withPrefix(
    prefix: string | void = "<magenta_bright>   [info] | </magenta_bright>",
    ...details: string[]
): string {
    return details.map(a => {
        return prefix + "<black_bright>" + a + "</black_bright>";
    }).join("\n");
}

export function colorize(message: string): string {
    return message.replace(
        /<(\w+(:[#\w:]+)?)>(.*?)<\/\1>/g,
        (_, color: string, __, text: string) => {
            const _color: string | undefined = AnsiColors.default(color);
            return _color ? _color + text + AnsiColors.ANSI_RESET : text;
        },
    );
}

export default class Logger {
    static log(...messages: string[]) {
        for (const message of messages) {
            console.log(colorize(message));
        }
    }

    static warn(message: string, ...details: string[]) {
        this.log(
            "<yellow_bright>[warning] | </yellow_bright>" + message,
            withPrefix(undefined, ...details),
        );
        
        console.log();
    }

    static err(why: string, help: string[], ...details: string[]) {
        this.log(
            "<red_bright>[error] | </red_bright>" + why,
            withPrefix(undefined, ...details),
            withPrefix("<blue_bright>   [help] | </blue_bright>", ...help),
            
        );
    }
}
