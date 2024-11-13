const AnsiColors: { [key: string]: string } = {
    black: "\u001b[30m",
    red: "\u001b[31m",
    green: "\u001b[32m",
    yellow: "\u001b[33m",
    blue: "\u001b[34m",
    magenta: "\u001b[35m",
    cyan: "\u001b[36m",
    white: "\u001b[37m",
    black_bright: "\u001b[90m",
    red_bright: "\u001b[91m",
    green_bright: "\u001b[92m",
    yellow_bright: "\u001b[93m",
    blue_bright: "\u001b[94m",
    magenta_bright: "\u001b[95m",
    cyan_bright: "\u001b[96m",
    white_bright: "\u001b[97m",
    reset: "\u001b[0m",
};

export const ANSI_RESET: string = AnsiColors.reset;
export default function get(name: string): string | undefined {
    return AnsiColors[name];
}
