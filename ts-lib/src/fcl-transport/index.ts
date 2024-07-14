import {
  Account,
  InteractionTag,
  Transaction,
  TransactionStatus,
} from "@onflow/typedefs";
import * as fcl from "@onflow/fcl";
import * as rlp from "@onflow/rlp";
import { Interaction } from "@onflow/typedefs/types/interaction";

// https://developers.flow.com/tools/clients/fcl-js/api#collectionobject
type Collection = {
  id: string;
  transactionIds: string[];
};

type NetworkParameters = {
  chainId: string;
};

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
  getTransaction: (id: string) => JsResponse<Transaction>;
  getTransactionsByBlockId: (blockId: string) => JsResponse<Transaction[]>;
  getTransactionResultsByBlockId: (
    blockId: string
  ) => JsResponse<TransactionStatus[]>;
  getTransactionResult: (
    transactionId: string
  ) => JsResponse<TransactionStatus>;
  getCollection: (id: string) => JsResponse<Collection[]>;
  // JSON encoded object matching SendSignedTransactionRequest Go struct
  sendSignedTransaction: (request: string) => JsResponse<string>;
  // https://github.com/onflow/fcl-js/pull/1420
  getNetworkParameters: () => JsResponse<NetworkParameters>;
  executeScript: (request: string) => JsResponse<string>;
  executeScriptAtHeight: (request: string) => JsResponse<string>;
  executeScriptAtId: (request: string) => JsResponse<string>;
}

// TODO: Expand to handle error cases
type JsResponse<Value> = Value;

export function buildWasmTransport(internalGateway: InternalGateway) {
  return async function transportWasm(
    _ix: Interaction | Promise<Interaction>,
    context: FclContext,
    _options: FclOptions
  ) {
    const ix = await _ix;
    switch (ix.tag) {
      case InteractionTag.GET_ACCOUNT:
        return {
          ...context.response(),
          tag: ix.tag,
          account: internalGateway.getAccount(ix.account.addr!),
        };
      case InteractionTag.GET_BLOCK:
        if (ix.block.isSealed) {
          return {
            ...context.response(),
            tag: ix.tag,
            block: internalGateway.getLatestBlock(),
          };
        }
        if (ix.block.id) {
          return {
            ...context.response(),
            tag: ix.tag,
            block: internalGateway.getBlockById(ix.block.id),
          };
        }
        if (ix.block.height !== undefined && ix.block.height !== null) {
          return {
            ...context.response(),
            tag: ix.tag,
            block: internalGateway.getBlockByHeight(Number(ix.block.height)),
          };
        }
        // No parameters are provided when fetching the reference block for a transaction.
        // See: https://github.com/onflow/fcl-js/blob/9c7873140015c9d1e28712aed93c56654f656639/packages/sdk/src/resolve/resolve.js#L83-L89
        return {
          ...context.response(),
          tag: ix.tag,
          block: internalGateway.getLatestBlock(),
        };
      case InteractionTag.GET_TRANSACTION:
        if (ix.transaction.id) {
          return {
            ...context.response(),
            tag: ix.tag,
            transaction: internalGateway.getTransaction(ix.transaction.id),
          };
        }
        throw new Error("Unreachable");
      case InteractionTag.GET_TRANSACTION_STATUS:
        if (ix.transaction.id) {
          return {
            ...context.response(),
            tag: ix.tag,
            transactionStatus: internalGateway.getTransactionResult(
              ix.transaction.id
            ),
          };
        }
        throw new Error("Unreachable");
      case InteractionTag.GET_COLLECTION:
        if (ix.collection.id) {
          return {
            ...context.response(),
            tag: ix.tag,
            collection: internalGateway.getCollection(ix.collection.id),
          };
        }
        throw new Error("Unreachable");
      case InteractionTag.TRANSACTION:
        return {
          ...context.response(),
          tag: ix.tag,
          transaction: internalGateway.sendSignedTransaction(
            JSON.stringify({
              gasLimit: Number(ix.message.computeLimit ?? 10),
              payer: ix.message.payer ?? "",
              referenceBlockId: ix.message.refBlock ?? "",
              script: ix.message.cadence ?? "",
              arguments: ix.message.arguments.map(argumentId =>
                JSON.stringify(ix.arguments[argumentId].asArgument)
              ),
            })
          ),
        };
      case InteractionTag.SCRIPT:
        return {
          ...context.response(),
          tag: ix.tag,
          encodedData: JSON.parse(
            internalGateway.executeScript(
              JSON.stringify({
                script: ix.message.cadence,
                arguments: JSON.stringify(
                  ix.message.arguments.map(
                    argumentId => ix.arguments[argumentId].asArgument
                  )
                ),
              })
            )
          ),
        };
      case InteractionTag.GET_NETWORK_PARAMETERS:
        return {
          ...context.response(),
          tag: ix.tag,
          networkParameters: internalGateway.getNetworkParameters(),
        };
      default:
        throw new Error(`Unimplemented interaction: ${JSON.stringify(ix)}`);
    }
  };
}
