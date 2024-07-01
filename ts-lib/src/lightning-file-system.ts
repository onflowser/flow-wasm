import FS from '@isomorphic-git/lightning-fs';
import {GoFileSystem, GoResult} from "@/go-interfaces";

export class LightningFileSystem implements GoFileSystem {
    private readonly fs: FS;
    private readonly rootDir: string;

    constructor(fs: FS, rootDir: string) {
        this.fs = fs;
        this.rootDir = rootDir;
    }

    async readFile(path: string): Promise<GoResult<string>> {
        const result = await this.fs.promises.readFile(`${this.rootDir}/${path}`, { encoding: "utf8"})
        return {
            error: null,
            value: result.toString()
        }
    }

}
