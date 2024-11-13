import * as yargs from "https://deno.land/x/yargs@v17.7.2-deno/deno.ts";
import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";
import commandList from "./generic/commandList.ts";

const args : Arguments = yargs.default(Deno.args)
        .command(
            "run [file]",
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
        .command(
            "init [where]",
            "Initializes a new project"
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
        .option("show-full", {
            alias: "sf",
            type: "boolean",
            description: "Shows the full license notice",
        })
        .option("show-performance", {
            alias: "sp",
            type: "boolean",
            description: "Show performance information",
        })
        .strictCommands()
        .demandCommand(1)
        .parse()

const [command] = args._;
const commandRunnable = commandList[command];

if (!commandRunnable) {
    console.log("TODO!");
} else {
    commandRunnable.run(args);
}