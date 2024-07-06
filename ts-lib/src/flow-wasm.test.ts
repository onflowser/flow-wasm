import {describe, it} from "vitest";
import * as fcl from "@onflow/fcl";
import {transportWasm} from "./fcl-transport";
import * as fs from "node:fs/promises";
import * as path from "path";
import "../../dist/wasm_exec.js";
import {FclGateway, FlowWasm, GoWasmRuntime, WasmGlobal} from "./index";
import {NullPrompter} from "./prompter/null-prompter";
import {InMemoryFileSystem} from "./filesystem/in-memory-file-system";
import {memfs} from 'memfs';

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

        flowWasm.run(goRuntime);

        fcl.config({
            "sdk.transport": transportWasm
        });

        await fcl.send([fcl.getAccount("0x7e60df042a9c0868")]);
    });
});
