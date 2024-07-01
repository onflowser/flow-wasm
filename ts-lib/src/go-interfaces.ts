import * as fclTypes from "@onflow/typedefs";

/**
 * Every function should return an object with both value and error fields,
 * so that errors can be handled appropriately on the go side.
 */
export interface GoResult<Data> {
    value: Data;
    error: string | null;
}

// Equivalent to os.FileMode in go.
type FileMode = number;

/**
 * Defines file system interface as implemented in /js/filesystem.go.
 */
export interface GoFileSystem {
    readFile(): Promise<GoResult<string>>;
    writeFile(filename: string, data: string, perm: FileMode): Promise<GoResult<null>>;
    mkdirAll(path: string, perm: FileMode): Promise<GoResult<null>>
    stat(path: string): Promise<GoResult<GoFileInfo>>
}

/**
 * Defines file info interface as implemented in /js/file_info.go
 */
export interface GoFileInfo {
    name: string;
    size: number;
    mode: FileMode;
    // Unix timestamp (in milliseconds).
    modTime: number;
    isDir: boolean;
}

/**
 * Defines Flow gateway interface as implemented in /js/gateway.go
 */
export interface GoFlowGateway {
    getAccount(address: string): Promise<GoResult<fclTypes.Account>>;
    // TODO: Define other functions
}
