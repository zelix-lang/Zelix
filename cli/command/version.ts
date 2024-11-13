import GenericCommand from "../generic/GenericCommand.ts";

export default class Version implements GenericCommand {
        
    run(): void {
        console.log("v1.0.0");
    }

}