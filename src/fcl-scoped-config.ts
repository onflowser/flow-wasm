import * as fcl from '@onflow/fcl';

type FclConfig = Record<string, unknown>;

export class FclScopedConfig<Identifier extends string> {
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
