import * as yargs from "https://deno.land/x/yargs@v17.7.2-deno/deno.ts";
import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";

const args : Arguments = yargs.default(Deno.args)
        .command(
            "run <file>",
            "Runs a file",
            // deno-lint-ignore no-explicit-any
            (yargs: any) => yargs.positional("file", { type: "string" }),
        )
        .command(
            "license",
            "Shows the license notice"
        )
        .command(
            "version",
            "Shows the currently installed version"
        )
        .option("save-errors", {
            alias: "se",
            type: "boolean",
            description: "Save errors to a file",
        })
        .option("dont-warn-style", {
            alias: "style",
            type: "boolean",
            description: "Don't warn about code style",
        })
        .strictCommands()
        .demandCommand(1)
        .parse()

console.log(args);
