import {NetworkId} from "./gateways/fcl-gateway";
import {GoFileSystem, GoFlowGateway, GoPrompter} from "@/go-interfaces";
import * as fclTypes from "@onflow/typedefs";

export { FclGateway } from "./gateways/fcl-gateway"
export { WindowPrompter }  from "./prompter/window-prompter";
export { LightningFileSystem } from "./filesystem/lightning-file-system"

type FlowWasmOptions = {
    gateways: Record<NetworkId, GoFlowGateway>;
    fileSystem: GoFileSystem;
    flowWasm: WebAssembly.WebAssemblyInstantiatedSource;
    prompter: GoPrompter;
    global: WasmGlobal;
}

/**
 * Global properties consumed or provided by go code.
 */
export interface WasmGlobal {
    // Consumed by Go runtime
    flowFileSystem: GoFileSystem;
    testnetGateway: GoFlowGateway;
    mainnetGateway: GoFlowGateway;
    previewnetGateway: GoFlowGateway;
    prompter: GoPrompter;
    // Provided by Go runtime
    Install: () => void;
    GetAccount: (address: string) => fclTypes.Account;
    GetLogs: () => string;
}

export interface GoWasmRuntime {
    run(instance: WebAssembly.Instance): Promise<void>;
    importObject: WebAssembly.Imports;
}

export class FlowWasm {
    constructor(private readonly options: FlowWasmOptions) {}

    public async run(goRuntime: GoWasmRuntime) {
        const { global } = this.options;

        // Configure runtime environment
        global.flowFileSystem = this.options.fileSystem;
        global.testnetGateway = this.options.gateways.testnet;
        global.mainnetGateway = this.options.gateways.mainnet;
        global.previewnetGateway = this.options.gateways.previewnet;
        global.prompter = this.options.prompter;

        await goRuntime.run(this.options.flowWasm.instance);
    }

    public async install(): Promise<void> {
        return this.options.global.Install();
    }

    public getAccount(address: string): fclTypes.Account {
        return this.options.global.GetAccount(address);
    }

    public getLogs(): string[] {
        return JSON.parse(this.options.global.GetLogs());
    }
}
