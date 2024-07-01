import {NetworkId} from "./fcl-gateway";
import {GoFileSystem, GoFlowGateway} from "@/go-interfaces";
import * as fclTypes from "@onflow/typedefs";

export { FclGateway } from "./fcl-gateway"
export { LightningFileSystem } from "./lightning-file-system"

type FlowWasmOptions = {
    gateways: Record<NetworkId, GoFlowGateway>;
    fileSystem: GoFileSystem;
    flowWasmUrl: string;
}

/**
 * Global properties consumed or provided by go code.
 */
declare global {
    interface Window {
        // Consumed by Go runtime
        flowFileSystem: GoFileSystem;
        testnetGateway: GoFlowGateway;
        mainnetGateway: GoFlowGateway;
        previewnetGateway: GoFlowGateway;
        // Provided by Go runtime
        Install: () => void;
        GetAccount: (address: string) => fclTypes.Account;
        GetLogs: () => string;
    }
}

interface GoWasmRuntime {
    run(instance: WebAssembly.Instance): Promise<void>;
    importObject: WebAssembly.Imports;
}

export class FlowWasm {
    constructor(private readonly options: FlowWasmOptions) {}

    public async run(goRuntime: GoWasmRuntime) {
        // Configure runtime environment
        window.flowFileSystem = this.options.fileSystem;
        window.testnetGateway = this.options.gateways.testnet;
        window.mainnetGateway = this.options.gateways.mainnet;
        window.previewnetGateway = this.options.gateways.previewnet;

        const wasm = await WebAssembly.instantiateStreaming(fetch(this.options.flowWasmUrl), goRuntime.importObject);
        await goRuntime.run(wasm.instance);
    }

    public async install(): Promise<void> {
        return window.Install();
    }

    public getAccount(address: string): fclTypes.Account {
        return window.GetAccount(address);
    }

    public getLogs(): string[] {
        return JSON.parse(window.GetLogs());
    }
}
