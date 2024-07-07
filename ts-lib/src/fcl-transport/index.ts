import {InteractionTag, Account} from "@onflow/typedefs";
import * as fcl from "@onflow/fcl";
import * as rlp from "@onflow/rlp";
import {Interaction} from "@onflow/typedefs/types/interaction";

type FclContext = unknown;

type FclOptions = {
    config: typeof fcl.config;
    ix: unknown;
    response: unknown;
    Buffer: typeof rlp.Buffer;
};

export interface InternalGateway {
    getAccount: (address: string) => JsResponse<Account>;
}

type JsResponse<Value> = {
    error: string;
    value: Value;
}

export function buildWasmTransport(internalGateway: InternalGateway) {
    return function transportWasm(
        ix: Interaction,
        _context: FclContext,
        _options: FclOptions
    ) {
        switch (ix.tag) {
            case InteractionTag.GET_ACCOUNT:
                return internalGateway.getAccount(ix.account.addr!);
            default:
                throw new Error(`Unimplemented interaction: ${ix.tag}`)
        }
    }
}

