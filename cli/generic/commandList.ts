import { Init } from "../command/init.ts";
import License from "../command/license.ts";
import Version from "../command/version.ts";
import GenericCommand from "./GenericCommand.ts";

const list: { [key: string]: GenericCommand } = {
    license: new License(),
    version: new Version(),
    init: new Init()
}

export default list;