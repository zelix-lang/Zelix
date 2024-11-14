import { Init } from "../command/init.ts";
import License from "../command/license.ts";
import Run from "../command/run.ts";
import Version from "../command/version.ts";
import GenericCommand from "./GenericCommand.ts";

const list: { [key: string]: GenericCommand } = {
    license: new License(),
    version: new Version(),
    init: new Init(),
    run: new Run()
}

export default list;