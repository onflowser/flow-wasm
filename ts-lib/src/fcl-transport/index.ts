import {InteractionTag, Account} from "@onflow/typedefs";
import * as fcl from "@onflow/fcl";
import * as rlp from "@onflow/rlp";
import {Interaction} from "@onflow/typedefs/types/interaction";

type FclContext = {
    config: typeof fcl.config;
    ix: unknown;
    response: () => any;
    Buffer: typeof rlp.Buffer;
};

type FclOptions = unknown;

export interface InternalGateway {
    getAccount: (address: string) => JsResponse<Account>;
    getBlockById: (id: string) => JsResponse<Account>;
    getLatestBlock: () => JsResponse<Account>;
    getBlockByHeight: (height: number) => JsResponse<Account>;
}

type JsResponse<Value> = {
    error: string;
    value: Value;
}

export function buildWasmTransport(internalGateway: InternalGateway) {
    return function transportWasm(
        ix: Interaction,
        context: FclContext,
        _options: FclOptions
    ) {
        switch (ix.tag) {
            case InteractionTag.GET_ACCOUNT:
                return {
                    ...context.response(),
                    tag: ix.tag,
                    account: internalGateway.getAccount(ix.account.addr!)
                };
            case InteractionTag.GET_BLOCK:
                if (ix.block.isSealed) {
                    return {
                        ...context.response(),
                        tag: ix.tag,
                        block: internalGateway.getLatestBlock()
                    };
                }
                if (ix.block.id) {
                    return {
                        ...context.response(),
                        tag: ix.tag,
                        block: internalGateway.getBlockById(ix.block.id)
                    };
                }
                if (ix.block.height !== undefined && ix.block.height !== null) {
                    return {
                        ...context.response(),
                        tag: ix.tag,
                        block: internalGateway.getBlockByHeight(Number(ix.block.height))
                    };
                }
                throw new Error("Unreachable")
            default:
                throw new Error(`Unimplemented interaction: ${ix.tag}`)
        }
    }
}

