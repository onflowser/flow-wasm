import FS from '@isomorphic-git/lightning-fs';
import {GoFileInfo, GoFileSystem, GoResult} from "@/go-interfaces";

export class LightningFileSystem implements GoFileSystem {
    private readonly fs: FS;
    private readonly rootDir: string;

    constructor(fs: FS, rootDir: string) {
        this.fs = fs;
        this.rootDir = rootDir;
    }

    async readFile(path: string): Promise<GoResult<string>> {
        try {
            const result = await this.fs.promises.readFile(this.scopedPath(path), { encoding: "utf8"})
            return {
                error: null,
                value: result.toString()
            }
        } catch (error) {
            return {
                error: String(error),
                value: null
            }
        }
    }

    async writeFile(path: string, data: string, perm: number): Promise<GoResult<null>> {
        try {
            await this.fs.promises.writeFile(this.scopedPath(path), data, {
                encoding: "utf8",
                mode: perm
            });
            return {
                error: null,
                value: null
            }
        } catch (error) {
            return {
                error: String(error),
                value: null
            }
        }
    }

    async mkdirAll(path: string, perm: number): Promise<GoResult<null>> {
        const pathSeparator = "/"
        const dirPaths = path.split(pathSeparator);

        for (let i = 1; i <= dirPaths.length; i++) {
            const pathUpToCurrentLevel = dirPaths.slice(0, i).join(pathSeparator);
            try {
                await this.fs.promises.mkdir(this.scopedPath(pathUpToCurrentLevel), { mode: perm });
            } catch (error) {
                // This function should be permissive, so this is not an error.
                // Move on to the next directory.
                if (String(error) === "Error: EEXIST") {
                    continue;
                }
                return {
                    error: String(error),
                    value: null
                }
            }
        }

        return {
            error: null,
            value: null
        }
    }

    async stat(path: string): Promise<GoResult<GoFileInfo>> {
        try {
            const result = await this.fs.promises.stat(this.scopedPath(path));
            return {
                error: null,
                value: {
                    name: path.split('/').pop() ?? "",
                    size: result.size,
                    mode: result.mode,
                    modTime: result.mtimeMs,
                    isDir: !result.isFile(),
                }
            }
        } catch (error) {
            return {
                error: String(error),
                value: null
            }
        }
    }

    private scopedPath(path: string): string {
        return `${this.rootDir}/${path}`
    }

}
