import { Arguments } from "https://deno.land/x/yargs@v17.7.2-deno/deno-types.ts";

export default interface GenericCommand {
    run(args: Arguments) : void;
}