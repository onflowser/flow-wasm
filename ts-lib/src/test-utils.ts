import fs from "node:fs/promises";
import path from "path";
import { memfs } from "memfs";
import { FclGateway, FlowWasm, GoWasmRuntime, WasmGlobal } from "@/index";
import { InMemoryFileSystem } from "@/filesystem/in-memory-file-system";
import { NullPrompter } from "@/prompter/null-prompter";

declare global {
  // This is implemented in wasm_exec.js
  class Go implements GoWasmRuntime {
    importObject: WebAssembly.Imports;

    run(instance: WebAssembly.Instance): Promise<void>;
  }
}

export async function runTestingFlowWasm(): Promise<FlowWasm> {
  const flowWasmBin = await fs.readFile(
    path.join(__dirname, "../../dist/flow.wasm")
  );
  const goRuntime = new Go();
  const wasmModule = await WebAssembly.instantiate(
    flowWasmBin,
    goRuntime.importObject
  );
  const rootDir = ".";
  const memFsInstance = memfs(
    {
      // language=Cadence
      "cadence/contracts/HelloWorld.cdc": `
                pub contract HelloWorld {
                    pub fun sayHello(): String {
                        return "Hello World"
                    }
                }
            `,
      "flow.json": JSON.stringify({
        contracts: {
          HelloWorld: "./cadence/contracts/HelloWorld.cdc",
        },
        accounts: {
          "emulator-account": {
            address: "f8d6e0586b0a20c7",
            key: "6d12eebfef9866c9b6fa92b97c6e705c26a1785b1e7944da701fc545a51d4673",
          },
        },
        networks: {
          emulator: "127.0.0.1:3569",
          mainnet: "access.mainnet.nodes.onflow.org:9000",
          sandboxnet: "access.sandboxnet.nodes.onflow.org:9000",
          testnet: "access.devnet.nodes.onflow.org:9000",
          previewnet: "access.previewnet.nodes.onflow.org:9000",
        },
      }),
    },
    rootDir
  );

  const flowWasm = new FlowWasm({
    flowWasm: wasmModule,
    fileSystem: new InMemoryFileSystem(memFsInstance.fs, rootDir),
    prompter: new NullPrompter(),
    global: global as unknown as WasmGlobal,
    gateways: {
      mainnet: new FclGateway("mainnet"),
      previewnet: new FclGateway("previewnet"),
      testnet: new FclGateway("testnet"),
    },
  });

  await flowWasm.run(goRuntime);

  return flowWasm;
}
