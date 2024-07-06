import {InteractionTag} from "@onflow/typedefs";
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

export function transportWasm(
    ix: Interaction,
    _context: FclContext,
    _options: FclOptions
) {
    switch (ix.tag) {
        case InteractionTag.GET_ACCOUNT:
        default:
            throw new Error(`Unimplemented interaction: ${ix.tag}`)
    }
}

