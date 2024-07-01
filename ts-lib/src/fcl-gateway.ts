import * as fcl from "@onflow/fcl";
import * as types from "@onflow/typedefs";
import {FclScopedConfig} from "./fcl-scoped-config";

export type NetworkId = "testnet" | "mainnet" | "previewnet"

export class FclGateway {
    private readonly scopedConfig: FclScopedConfig<NetworkId>;

    constructor(private readonly network: NetworkId) {
        this.scopedConfig = new FclScopedConfig<NetworkId>();

        this.scopedConfig.setConfig("testnet", {
            'flow.network': 'testnet',
            'accessNode.api': 'https://rest-testnet.onflow.org',
        });

        this.scopedConfig.setConfig("mainnet", {
            'flow.network': 'mainnet',
            'accessNode.api': 'https://rest-mainnet.onflow.org',
        });

        this.scopedConfig.setConfig("previewnet", {
            'flow.network': 'previewnet',
            'accessNode.api': 'https://rest-previewnet.onflow.org',
        });
    }


    async getAccount(address: string): Promise<types.Account> {
        this.prepare();
        return fcl.account(address);
    }

    private prepare() {
        this.scopedConfig.useConfig(this.network);
    }
}
