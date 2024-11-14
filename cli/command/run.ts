import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";
import GenericCommand from "../generic/GenericCommand.ts";
import { join } from "@std/path/join";
import retrievePath from "../lib/retrievePath.ts";
import Logger from "../../shared/Logger.ts";
import { parse } from "@std/yaml/parse";
import Lexer from "../../lexer/index.ts";
import { fileExists } from "../../shared/util/FileUtil.ts";
import asc from "assemblyscript/asc";
import extractFunction from "../../core/functionExtractor.ts";
import { catchException } from "../../shared/catchException.ts";

export default class Run implements GenericCommand {
    
    run(args: Arguments): void {
        const where = retrievePath(args.file ?? ".");

        if (!fileExists(where) || !Deno.statSync(where).isDirectory) {
            catchException(
                "Invalid directory.",
                ["The directory is not valid."],
                "You can't run a project without a valid directory path."
            )
        };

        const files = Array.from(Deno.readDirSync(where));

        if (!files.some(a => a.name === "Surf.yml")) {
            Logger.err(
                "The directory does not contain a Surf.yml file.",
                ["Initialize a project with</black_bright><blue_bright> surf init</blue_bright><black_bright> or create a Surf.yml file manually."],
                "You can't run a project without a Surf configuration file."
            );

            Deno.exit(1);
        }
        
        const config: object = parse(Deno.readTextFileSync(join(where, "Surf.yml"))) as object;
        
        if (!("main_script" in config) || typeof config["main_script"] !== "string") {
            Logger.err(
                "The Surf.yml file does not contain a main_script field.",
                ["Add a main_script field to the Surf.yml file."],
                "You can't run a project without a main_script field."
            );

            Deno.exit(1);
        }

        if (!fileExists(join(where, config["main_script"]))) {
            Logger.err(
                "The main script file does not exist.",
                ["Create the main script file or change the main_script field in the Surf.yml file."],
                "You can't run a project without the main script file."
            );

            Deno.exit(1);
        }

        const scriptPath = join(where, config["main_script"]);
        
        
        const tokens = new Lexer().tokenize(scriptPath, Deno.readTextFileSync(scriptPath));
        const [ functions, imports, declaredFunctions ] = extractFunction(tokens);
    }

}