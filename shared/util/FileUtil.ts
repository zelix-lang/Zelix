export function fileExists(path: string): boolean {
    try {
        Deno.lstatSync(path);
        return true;
    } catch {
        return false;
    }
}