import Logger from "./Logger.ts";

export function catchException(why: string, help: string[], ...details: string[]) {
    Logger.err("<red_bright>" + why + "</red_bright>", help, ...details);
    Deno.exit(1);
}