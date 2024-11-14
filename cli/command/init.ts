import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";
import GenericCommand from "../generic/GenericCommand.ts";
import Logger, { colorize } from "../../shared/Logger.ts";
import { stringify } from "@std/yaml";
import { fileExists } from "../../shared/util/FileUtil.ts";
import { join } from "@std/path";
import messageHeader from "../message/header.ts";
import question, { booleanQuestion } from "../lib/question.ts";
import { step } from "jsr:@sylc/step-spinner"
import { emojiOf } from "@lambdalisue/github-emoji";
import exampleSurfScript from "../generic/exampleSurfScript.ts";
import retrievePath from "../lib/retrievePath.ts";

const rocketEmoji = emojiOf("rocket");
const checkingSpinner = step("Checking...");
const writingSpinner = step("Writing...");

export class Init implements GenericCommand {

    async run(args: Arguments) {
        const cwd = Deno.cwd();
        const where = retrievePath(args.where);

        if (!fileExists(where)) {
            Deno.mkdirSync(where, { recursive: true });
        }

        const files = Array.from(Deno.readDirSync(where));

        if (files.length > 0) {
            Logger.err(
                "The directory is not empty.",
                ["You can create an empty directory and run the command again."],
                "You can't initialize a project in a non-empty directory.",
                "Directory specified: " + where
            )

            Deno.exit(1);
        }

        messageHeader();
        Logger.log(
            "<black_bright>You may cancel this process at any time by pressing Ctrl + C.</black_bright>",
            "<black_bright>This will not affect the directory in any way unless you confirm the operation.</black_bright>",
            ""
        );

        const textEncoder = new TextEncoder();
        const name = await question("Project name", "Surf-Project", /^[a-zA-Z\-_]+((\d+)?)$/);
        const description = await question("Description", "A simple Surf project", /.+/);
        const version = await question("Version", "1.0.0", /^\d+\.\d+\.\d+$/);
        const author = await question("Author", "John Doe", /^[a-zA-Z\.\s0-9]+$/);
        const license = await question("License", "MIT", /.+/);
        const repository = await question("Git Repository", "None", /.+/);
        const initGitRepo = await booleanQuestion("Initialize a Git repository?",);

        // Logger.log() not needed here
        // We're printing an empty string
        // there's no necessity to pass it through the colorizer
        console.log();
        
        checkingSpinner.start();

        const isGitInstalled = await new Deno.Command("git", {
            args: ["--version"],
            stdout: "piped"
        }).spawn().status;

        checkingSpinner.message = colorize("<green_bright>Checked</green_bright>");
        if (!isGitInstalled.success) {
            Logger.err(
                "Git is not installed.",
                ["You can install Git from the official website."],
                "You can't initialize a Git repository without Git installed.",
                "The project will still be initialized, but without a Git repository."
            );
            checkingSpinner.fail();
        } else {
            checkingSpinner.succeed();

            if (initGitRepo) {
                const gitInit = await new Deno.Command("git", {
                    args: ["init"],
                    cwd: where
                }).spawn().status;
    
                if (!gitInit.success) {
                    Logger.err(
                        "An error occurred while initializing the Git repository.",
                        ["You can try again."],
                        "The project was not initialized."
                    );
    
                    Deno.exit(1);
                }
            }
        }

        const data = stringify({
            name,
            version,
            author,
            description,
            git: repository,
            license,
            repositories: [],
            dependencies: [],
            main_script: "src/main.surf",
            bind: [
                { name: "project_version", value: "$self.version" },
                { name: "project_name", value: "$self.name" }
            ]
        })

        Logger.log(
            "The following contents are going to be written to Surf.yml:",
            "",
            data
        );

        await booleanQuestion("Is this okay?", true, () => {
            Deno.removeSync(where, { recursive: true });
        });

        writingSpinner.start();

        try {
            Deno.mkdirSync(join(where, "src"), { recursive: true });
            Deno.writeFileSync(
                join(where, "src", "main.surf"),
                textEncoder.encode(exampleSurfScript)
            );

            Deno.writeFileSync(
                join(where, "Surf.yml"),
                textEncoder.encode(
                    "# ---------------------------------------------------\n" + 
                    "#   This file is generated by the Surf CLI.\n" + 
                    "#   Feel free to modify the file to suit your needs.\n" +
                    "# ----------------------------------------------------\n\n" +
                    data
                )
            );
    
            writingSpinner.message = colorize("<green_bright>Wrote</green_bright>");
            writingSpinner.succeed();
        } catch {
            Deno.removeSync(where, { recursive: true });
            writingSpinner.message = colorize("<red_bright>Failed</red_bright>");
            writingSpinner.fail();
            Logger.err(
                "An error occurred while writing the file.",
                ["You can try again."],
                "The project was not initialized."
            );

            Deno.exit(1);
        }

        Logger.log(
            "",
            rocketEmoji + " Your project is ready!",
            "",
            "  Next steps:",
            ""
        );

        const printCd = where === cwd;

        if (!printCd) {
            Logger.log(
                "    <magenta_bright>1. cd " + args.where + "</magenta_bright>"
            );
        }

        Logger.log(
            `    <magenta_bright>${printCd ? 1 : 2}. surf run</magenta_bright>`,
            ""
        );
    }

}