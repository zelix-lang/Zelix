import { isAbsolute } from "@std/path/is-absolute";
import { join } from "@std/path/join";

export default function retrievePath(path: string) : string {
    const cwd = Deno.cwd();
    return isAbsolute(path) ? path : join(cwd, path);
}