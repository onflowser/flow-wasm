import * as fcl from "@onflow/fcl";
import * as types from "@onflow/typedefs";
import { GoFlowAccount, GoFlowGateway, GoResult } from "@/go-interfaces";

export type NetworkId = "testnet" | "mainnet" | "previewnet";

export class FclGateway implements GoFlowGateway {
  private readonly scopedConfig: FclScopedConfig<NetworkId>;

  constructor(private readonly network: NetworkId) {
    this.scopedConfig = new FclScopedConfig<NetworkId>();

    this.scopedConfig.setConfig("testnet", {
      "flow.network": "testnet",
      "accessNode.api": "https://rest-testnet.onflow.org",
    });

    this.scopedConfig.setConfig("mainnet", {
      "flow.network": "mainnet",
      "accessNode.api": "https://rest-mainnet.onflow.org",
    });

    this.scopedConfig.setConfig("previewnet", {
      "flow.network": "previewnet",
      "accessNode.api": "https://rest-previewnet.onflow.org",
    });
  }

  async getAccount(address: string): Promise<GoResult<GoFlowAccount>> {
    this.prepare();
    return resolveToGoResult<types.Account, GoFlowAccount>(
      () => fcl.account(address),
      {
        transform: value => {
          return {
            address: value.address,
            balance: value.balance,
            code: String(value.code), // TODO: Report an issue for mismatched type?
            contracts: JSON.stringify(value.contracts),
            keys: JSON.stringify(value.keys),
          };
        },
      }
    );
  }

  private prepare() {
    this.scopedConfig.useConfig(this.network);
  }
}

async function resolveToGoResult<Value, TValue>(
  asyncFunction: () => Promise<Value>,
  options: {
    transform: (value: Value) => TValue;
  }
): Promise<GoResult<TValue>> {
  try {
    const result = options.transform(await asyncFunction());
    return {
      error: null,
      value: result,
    };
  } catch (error) {
    return {
      error: String(error),
      value: null,
    };
  }
}

type FclConfig = Record<string, unknown>;

class FclScopedConfig<Identifier extends string> {
  private readonly configurations: Map<Identifier, FclConfig>;

  constructor() {
    this.configurations = new Map();
  }

  setConfig(identifier: Identifier, config: FclConfig) {
    this.configurations.set(identifier, config);
  }

  useConfig(identifier: Identifier) {
    const config = this.configurations.get(identifier);
    if (!config) {
      throw new Error(`Configuration '${identifier}' not found.`);
    }
    // Apply configuration to fcl locally within the closure
    const fclConfig = fcl.config();
    Object.keys(config).forEach(key => {
      fclConfig.put(key, config[key]);
    });
    return fclConfig;
  }
}
