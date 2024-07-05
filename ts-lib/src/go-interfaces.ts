/**
 * Every function should return an object with both value and error fields,
 * so that errors can be handled appropriately on the go side.
 */
export interface GoResult<Data> {
    value: Data | null;
    error: string | null;
}

// Equivalent to os.FileMode in go.
type FileMode = number;

/**
 * Defines file system interface as implemented in /js/filesystem.go.
 */
export interface GoFileSystem {
    readFile(filename: string): Promise<GoResult<string>>;
    // TODO: Implement
    // writeFile(filename: string, data: string, perm: FileMode): Promise<GoResult<null>>;
    // mkdirAll(path: string, perm: FileMode): Promise<GoResult<null>>
    // stat(path: string): Promise<GoResult<GoFileInfo>>
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

export type GoFlowAccount = {
    address: string
    balance: number
    code: string
    // JSON encoded map of contracts
    contracts: string;
    // JSON encoded map of keys
    keys: string;
}

/**
 * Defines Flow gateway interface as implemented in /js/gateway.go
 */
export interface GoFlowGateway {
    getAccount(address: string): Promise<GoResult<GoFlowAccount>>;
    // TODO: Define other functions
}

/**
 * Utilities for getting user input as defined in /js/prompter.go.
 */
export interface GoPrompter {
    shouldUpdateDependency(contractName: string): boolean;
    // Must return selected account name.
    addContractToDeployment(networkName: string, accounts: PrompterAccount[], contractName: string): string;
    addressPromptOrEmpty(label: string): string;
}

export type PrompterAccount = {
    Name: string;
    Address: string;
    Key: string;
}
