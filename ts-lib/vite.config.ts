/// <reference types="vitest" />
// @ts-ignore
import path from "path";
import { defineConfig } from "vite";
// @ts-ignore
import packageJson from "../package.json";

const getPackageName = () => {
  return packageJson.name;
};

const getPackageNameCamelCase = () => {
  try {
    return getPackageName().replace(/-./g, char => char[1].toUpperCase());
  } catch (err) {
    throw new Error("Name property in package.json is missing.");
  }
};

type Format = "es" | "cjs";

const fileNameLookup = new Map<Format, string>([
  ["es", `${getPackageName()}.mjs`],
  ["cjs", `${getPackageName()}.cjs`]
])

const formats = Array.from(fileNameLookup.keys());

module.exports = defineConfig({
  base: "./",
  build: {
    outDir: "./dist",
    lib: {
      entry: path.resolve(__dirname, "src/index.ts"),
      name: getPackageNameCamelCase(),
      formats,
      fileName: (format) => {
        const filename = fileNameLookup.get(format as Format);

        if (!filename) {
          throw new Error("Filename not found for format: " + format)
        }

        return filename;
      },
    },
  },
  test: {
    setupFiles: ['./setupTests.js'],
  },
  resolve: {
    alias: [
      { find: "@", replacement: path.resolve(__dirname, "src") },
      { find: "@@", replacement: path.resolve(__dirname) },
    ],
  },
});
