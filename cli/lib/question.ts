import Logger, { colorize } from "../../shared/Logger.ts";
import { emojiOf } from "@lambdalisue/github-emoji";
import { readKeypress } from "./readKeys.ts";
import getColor, { ANSI_RESET } from "../../shared/ansi/AnsiColors.ts";

const crossMark = emojiOf("heavy_multiplication_x");
const greenBrightColor = getColor("green_bright");
const yellowBrightColor = getColor("yellow_bright");
const checkMark = emojiOf("heavy_check_mark");
const arrowRightHook = emojiOf("arrow_right_hook");

async function buildColorizedPrompt(prompt: string, additional: string | void) : Promise<[string, string]> {
    const encoder = new TextEncoder();
    const successMessage = colorize("<green_bright>" + checkMark + " " + prompt + "</green_bright>");
    const errorMessage = colorize("<red_bright>" + crossMark + " " + prompt + "</red_bright>");

    await Deno.stdout.write(encoder.encode(
        colorize("<blue_bright>" + arrowRightHook + " "
            + prompt + (additional ? " " + additional : "") + ": </blue_bright>")
        + yellowBrightColor
    ));

    return [successMessage, errorMessage];
}

function checkForUserCancellation(keypress: { sequence: string }) {
    if (keypress.sequence === "\x03") {
        Logger.log("<red_bright>Process cancelled by user.</red_bright>");
        Deno.exit(1);
    }
}

export default async function question(prompt: string, fallback: string, checkForRegex: RegExp) : Promise<string> {
    while (true) {
        const encoder = new TextEncoder();
        const [successMessage, errorMessage] = await buildColorizedPrompt(prompt);

        let buffer = "";

        while (true) {
            const {value: keypress} = await readKeypress().next();
            if (keypress.key === "return") {
                if (!buffer.trim()) buffer = fallback;
                break;
            }
            else if (keypress.key === "backspace" || keypress.sequence === "\x17") {
                if (buffer.length === 0) continue;

                if (keypress.sequence === "\x17") {
                    const lastWord = buffer.split(" ").pop();

                    if (!lastWord || !lastWord.trim()) continue;
                    
                    buffer = buffer.substring(0, buffer.length - lastWord.length);

                    for (let i = 1; i < lastWord.length; i++) {
                        await Deno.stdout.write(encoder.encode("\b \b"));
                    }
                } else {
                    buffer = buffer.slice(0, buffer.length - 1);
                }

                await Deno.stdout.write(encoder.encode("\b \b"));
                continue;
            }

            checkForUserCancellation(keypress);
            if (keypress.key.length !== 1 && keypress.key !== "space") continue;
            const add = keypress.key === "space" ? " " : keypress.key;
            
            buffer += add;
            await Deno.stdout.write(encoder.encode(add));
        }

        if (!checkForRegex.test(buffer)) {
            Deno.stdout.write(encoder.encode("\r" + errorMessage + " ".repeat(buffer.length + 2) +"\n"));
            Logger.log("<red_bright>Invalid value entered, please try again.</red_bright>");
            continue;
        }
        
        Deno.stdout.write(encoder.encode(
            "\r" + successMessage + greenBrightColor + " (" + (buffer || fallback || "empty") + ")" + ANSI_RESET + "\n"
        ));
        
        return buffer || fallback;
    }
}

export async function booleanQuestion(
    prompt: string,
    cancelOnNo: boolean = false,
    onCancel?: () => void
) : Promise<boolean> {
    const encoder = new TextEncoder();
    const [successMessage] = await buildColorizedPrompt(prompt, "[Y/n]");

    const {value: keypress} = await readKeypress().next();
    const key = keypress.key.toLowerCase();

    const value = key === "y" || key === "return";

    if (!value && cancelOnNo) {
        if (onCancel) onCancel();
        Logger.log("<red_bright>Operation cancelled.</red_bright>");
        Deno.exit(0);
    }
    
    checkForUserCancellation(keypress);
    Deno.stdout.write(encoder.encode(
        "\r" + successMessage + greenBrightColor + " (" + (value ? "yes" : "no") + ")"
            // Add 5 whitespaces to make dissappear the [Y/n] part
            + " ".repeat(5)
            + ANSI_RESET + "\n"
    ));

    return value;
}