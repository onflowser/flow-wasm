import FS from '@isomorphic-git/lightning-fs';

export class LightningFileSystem {
    private readonly fs: FS;
    private readonly rootDir: string;

    constructor(fs: FS, rootDir: string) {
        this.fs = fs;
        this.rootDir = rootDir;
    }

    async readFile(path: string) {
        console.log("readFile", path)
        const result = await this.fs.promises.readFile(`${this.rootDir}/${path}`, { encoding: "utf8"})
        console.log("readFile", result)
        return result;
    }

    // TODO: Implement other functions
}
