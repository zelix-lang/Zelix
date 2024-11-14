export default class Import {

    private readonly path: string;
    private readonly alias: string;
    private readonly names: string;

    constructor(path: string, alias: string, names: string) {
        this.path = path;
        this.alias = alias;
        this.names = names;
    }

    public getPath(): string {
        return this.path;
    }

    public getAlias(): string {
        return this.alias;
    }

    public getNames(): string {
        return this.names;
    }

}