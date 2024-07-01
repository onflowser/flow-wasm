import {NetworkId} from "./fcl-gateway";
import {GoFileSystem, GoFlowGateway} from "@/go-interfaces";

export { FclGateway } from "./fcl-gateway"
export { LightningFileSystem } from "./lightning-file-system"

type FlowWasmOptions = {
    gateways: Record<NetworkId, GoFlowGateway>;
    fileSystem: GoFileSystem;
    flowWasmUrl: string;
}

/**
 * Global properties expected by go code.
 */
declare global {
    interface Window {
        flowFileSystem: GoFileSystem;
        testnetGateway: GoFlowGateway;
        mainnetGateway: GoFlowGateway;
        previewnetGateway: GoFlowGateway;
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
}
