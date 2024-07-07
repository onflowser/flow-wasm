import {beforeAll, describe, expect, it} from "vitest";
import * as fcl from "@onflow/fcl";
import "../../../dist/wasm_exec.js";
import {
    Account,
    Block,
    CollectionGuarantee,
    InteractionAccount,
    Transaction,
    TransactionStatus
} from "@onflow/typedefs";
import {HashAlgorithm, SignatureAlgorithm} from "@onflow/typedefs/src";
import {runTestingFlowWasm} from "../test-utils";
import {Collection} from "../../../dist";

async function prepareTests() {
    const flowWasm = await runTestingFlowWasm();

    fcl.config({
        "sdk.transport": flowWasm.fclTransport(),
        "fcl.authz": flowWasm.serviceAccountAuthz(),
        // This is here because FCL expects it, but it's not actually being used
        'flow.network': 'testnet',
        'accessNode.api': 'https://rest-testnet.onflow.org',
        'discovery.wallet': 'https://fcl-discovery.onflow.org/testnet/authn'
    });
}


describe("FCL transport - accounts", () => {
    beforeAll(prepareTests);

    it("should get account by address", async () => {
        const actual = await fcl.send([fcl.getAccount("0xf8d6e0586b0a20c7")]).then(fcl.decode);

        const expected: Account = {
            address: "f8d6e0586b0a20c7",
            balance: 100000000000000000,
            // @ts-ignore incorrect type for 'code' field in FCL
            code: "",
            contracts: {},
            keys: [
                {
                    hashAlgoString: "SHA3_256",
                    hashAlgo: HashAlgorithm.SHA3_256,
                    weight: 1000,
                    sequenceNumber: 0,
                    revoked: false,
                    index: 0,
                    publicKey: "0x43661ddd40c0510b2097a5ad583607f4780876184308a325516951fac6a816fe4e522c9278d3ef3d67c6d903291d0501f9a9bd5b4dc2c5af26c2ad0597bac97a",
                    signAlgoString: "ECDSA_P256",
                    signAlgo: SignatureAlgorithm.ECDSA_P256
                }
            ]
        }

        expect(actual).toMatchObject(expected);
    });
});

describe("FCL transport - blocks", async () => {
    beforeAll(prepareTests);

    it('should get latest sealed block', async () => {
        const actual = await fcl.send([fcl.getBlock(true)]).then(fcl.decode);

        const expected: Block = {
            "blockSeals": [],
            "collectionGuarantees": [],
            "signatures": [], // Not implemented
            "height": 0,
            "id": "a20c602fbee6fe4491e116403e3258e7b7924609696ab2edb9a93eed2c29e445",
            "parentId": "0000000000000000000000000000000000000000000000000000000000000000",
            "timestamp": "2018-12-19 22:32:30.000000042 +0000 UTC",
        }

        expect(actual).toMatchObject(expected);
    });

    it('should get block by id', async () => {
        const actual = await fcl.send([fcl.getBlock(), fcl.atBlockId("a20c602fbee6fe4491e116403e3258e7b7924609696ab2edb9a93eed2c29e445")]).then(fcl.decode);

        const expected: Block = {
            "blockSeals": [],
            "collectionGuarantees": [],
            "signatures": [], // Not implemented
            "height": 0,
            "id": "a20c602fbee6fe4491e116403e3258e7b7924609696ab2edb9a93eed2c29e445",
            "parentId": "0000000000000000000000000000000000000000000000000000000000000000",
            "timestamp": "2018-12-19 22:32:30.000000042 +0000 UTC",
        }

        expect(actual).toMatchObject(expected);
    });

    it('should get block by height', async () => {
        const actual = await fcl.send([fcl.getBlock(), fcl.atBlockHeight(0)]).then(fcl.decode);

        const expected: Block = {
            "blockSeals": [],
            "collectionGuarantees": [],
            "signatures": [], // Not implemented
            "height": 0,
            "id": "a20c602fbee6fe4491e116403e3258e7b7924609696ab2edb9a93eed2c29e445",
            "parentId": "0000000000000000000000000000000000000000000000000000000000000000",
            "timestamp": "2018-12-19 22:32:30.000000042 +0000 UTC",
        }

        expect(actual).toMatchObject(expected);
    });
})

// language=Cadence
const helloWorldCadenceTx = `
    transaction {
        execute {
            log("Hello World")
        }
    }
`

describe("FCL transport - collections", async () => {
    beforeAll(prepareTests)

    it('should return collection by ID', async () => {
        const transactionId = await fcl.mutate({
            cadence: helloWorldCadenceTx,
            limit: 10
        });

        const status = await fcl.tx(transactionId).snapshot();

        const block = await fcl.send([fcl.getBlock(), fcl.atBlockId(status.blockId)]).then(fcl.decode);

        expect(block.collectionGuarantees).toHaveLength(1);

        const collectionId = block.collectionGuarantees[0].collectionId;

        const actual = await fcl.send([fcl.getCollection(collectionId)]).then(fcl.decode);

        const expected: Collection = {
            id: collectionId,
            transactionIds: [transactionId]
        }

        expect(actual).toMatchObject(expected)
    });
})

describe("FCL transport - transactions", async () => {
    beforeAll(prepareTests);

    it('should send a transaction and return status', async () => {
        const transactionId = await fcl.mutate({
            cadence: helloWorldCadenceTx,
            limit: 10
        });

        const status = await fcl.send([fcl.getTransactionStatus(transactionId)]).then(fcl.decode);

        expect(status).toMatchObject({
            blockId: expect.any(String),
            statusCode: 0,
            status: 4,
            statusString: "SEALED",
            events: [],
            errorMessage: "",
        })

        expect(transactionId).toBeTypeOf("string");
    });

    it('should get transaction by ID', async () => {
        const transactionId = await fcl.mutate({
            cadence: helloWorldCadenceTx,
            limit: 10
        });

        const actual = await fcl.send([fcl.getTransaction(transactionId)]).then(fcl.decode);

        // Omit some fields that are not specified in the docs.
        // See: https://developers.flow.com/tools/clients/fcl-js/api#transactionobject
        const expected: Omit<Transaction, "proposer" | "address" | "keyId" | "sequenceNumber"> = {
            script: helloWorldCadenceTx,
            authorizers: [],
            envelopeSignatures: [],
            gasLimit: 10,
            proposalKey: {
                sequenceNumber: 0,
                keyId: 0,
                address: "0000000000000000",
            },
            args: [],
            referenceBlockId: "a20c602fbee6fe4491e116403e3258e7b7924609696ab2edb9a93eed2c29e445",
            payer: "f8d6e0586b0a20c7",
            payloadSignatures: []
        }

        expect(actual).toMatchObject(expected);
    });
})
