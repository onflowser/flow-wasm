import {describe, expect, it} from "vitest";
import * as fcl from "@onflow/fcl";
import * as fs from "node:fs/promises";
import * as path from "path";
import "../../dist/wasm_exec.js";
import {FclGateway, FlowWasm, GoWasmRuntime, WasmGlobal} from "./index";
import {NullPrompter} from "./prompter/null-prompter";
import {InMemoryFileSystem} from "./filesystem/in-memory-file-system";
import {memfs} from 'memfs';
import {Account} from "@onflow/typedefs";
import {HashAlgorithm, SignatureAlgorithm} from "@onflow/typedefs/src";

declare global {
    // This is implemented in wasm_exec.js
    class Go implements GoWasmRuntime {
        importObject: WebAssembly.Imports;

        run(instance: WebAssembly.Instance): Promise<void>;
    }
}

describe("Flow WASM", () => {
    it("should get account resource", async () => {
        const flowWasmBin = await fs.readFile(path.join(__dirname, "../../dist/flow.wasm"));
        const goRuntime = new Go();
        const wasmModule = await WebAssembly.instantiate(flowWasmBin, goRuntime.importObject);
        const rootDir = ".";
        const memFsInstance = memfs({
            // language=Cadence
            "cadence/contracts/HelloWorld.cdc": `
                pub contract HelloWorld {
                    pub fun sayHello(): String {
                        return "Hello World"
                    }
                }
            `,
            "flow.json": JSON.stringify({
                "contracts": {
                    "HelloWorld": "./cadence/contracts/HelloWorld.cdc"
                },
                "accounts": {
                    "emulator-account": {
                        "address": "f8d6e0586b0a20c7",
                        "key": "6d12eebfef9866c9b6fa92b97c6e705c26a1785b1e7944da701fc545a51d4673"
                    }
                },
                "networks": {
                    "emulator": "127.0.0.1:3569",
                    "mainnet": "access.mainnet.nodes.onflow.org:9000",
                    "sandboxnet": "access.sandboxnet.nodes.onflow.org:9000",
                    "testnet": "access.devnet.nodes.onflow.org:9000",
                    "previewnet": "access.previewnet.nodes.onflow.org:9000"
                },
            })
        }, rootDir);

        const flowWasm = new FlowWasm({
            flowWasm: wasmModule,
            fileSystem: new InMemoryFileSystem(memFsInstance.fs, rootDir),
            prompter: new NullPrompter(),
            global: global as unknown as WasmGlobal,
            gateways: {
                mainnet: new FclGateway("mainnet"),
                previewnet: new FclGateway("previewnet"),
                testnet: new FclGateway("testnet"),
            }
        });

        await flowWasm.run(goRuntime);

        fcl.config({"sdk.transport": flowWasm.getFclTransport()});

        const response = await fcl.send([fcl.getAccount("0xf8d6e0586b0a20c7")]).then(fcl.decode);

        const expectedAccount: Account = {
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

        expect(response).toMatchObject(expectedAccount);
    });
});
