import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";
import GenericCommand from "../generic/GenericCommand.ts";
import Logger from "../../shared/Logger.ts";

export default class License implements GenericCommand {
    
    run(args: Arguments): void {
        if (args.sf) {
            Logger.log(
                "You may find the complete license notie in",
                "  <yellow_bright>|</yellow_bright> <black_bright>the LICENSE file of the GitHub repository</black_bright>",
                "  <blue_bright>|</blue_bright> <black_bright>the official GNU website:</black_bright>",
                "    <blue_bright>https://www.gnu.org/licenses/gpl-3.0.en.html</blue_bright>"
            );
            return;
        }

        Logger.log(
            "<blue_bright>Surf</blue_bright> Copyright (C) 2024 Rodrigo R. & all Surf contributors",
            "This program comes with ABSOLUTELY NO WARRANTY; for details type `surf license'.",
            "This is free software, and you are welcome to redistribute it",
            "under certain conditions; type `surf license --show-full' for details."
        )
    }

} 