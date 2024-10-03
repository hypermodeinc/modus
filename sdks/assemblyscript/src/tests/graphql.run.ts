import { readFileSync } from "fs";
import { instantiate } from "../build/graphql.spec.js";
const binary = readFileSync("./build/graphql.spec.wasm");
const module = new WebAssembly.Module(binary);
instantiate(module, {
  env: {},
  hypermode: {},
});
