import { NetworkId } from "./gateways/fcl-gateway";
import { GoFileSystem, GoFlowGateway, GoPrompter } from "@/go-interfaces";
import { buildWasmTransport, InternalGateway } from "@/fcl-transport";
import { InteractionAccount } from "@onflow/typedefs";

export { FclGateway } from "./gateways/fcl-gateway";
export { WindowPrompter } from "./prompter/window-prompter";
export { LightningFileSystem } from "./filesystem/lightning-file-system";

type FlowWasmOptions = {
  gateways: Record<NetworkId, GoFlowGateway>;
  fileSystem: GoFileSystem;
  flowWasm: WebAssembly.WebAssemblyInstantiatedSource;
  prompter: GoPrompter;
  global: WasmGlobal;
};

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
  // Called when the emulator starts and initializes APIs
  onStarted: () => void;
  // Provided by Go runtime
  gateway: InternalGateway;
  install: () => void;
  getLogs: () => string;
  deploy: () => string;
}

export interface GoWasmRuntime {
  run(instance: WebAssembly.Instance): Promise<void>;
  importObject: WebAssembly.Imports;
}

export class FlowWasm {
  constructor(private readonly options: FlowWasmOptions) {}

  public async run(goRuntime: GoWasmRuntime) {
    return new Promise<void>(resolve => {
      const { global } = this.options;

      // Configure runtime environment
      global.flowFileSystem = this.options.fileSystem;
      global.testnetGateway = this.options.gateways.testnet;
      global.mainnetGateway = this.options.gateways.mainnet;
      global.previewnetGateway = this.options.gateways.previewnet;
      global.prompter = this.options.prompter;
      global.onStarted = resolve;

      goRuntime.run(this.options.flowWasm.instance);
    });
  }

  public fclTransport() {
    return buildWasmTransport(this.options.global.gateway);
  }

  public async install(): Promise<void> {
    return this.options.global.install();
  }

  public getLogs(): string[] {
    return JSON.parse(this.options.global.getLogs());
  }

  public deploy() {
    this.options.global.deploy();
  }

  // Authorization function for signing with service account
  // https://developers.flow.com/tools/clients/fcl-js/api#authz
  public serviceAccountAuthz() {
    return function (authAccount: InteractionAccount): InteractionAccount {
      return {
        ...authAccount,
        addr: "f8d6e0586b0a20c7",
        keyId: 0,
        signingFunction: () => ({
          addr: "0xf8d6e0586b0a20c7",
          keyId: 0,
          signature: "",
        }),
      };
    };
  }
}
