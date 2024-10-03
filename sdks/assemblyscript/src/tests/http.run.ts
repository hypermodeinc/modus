import { readFileSync } from "fs";
import { instantiate } from "../build/http.spec.js";
const binary = readFileSync("./build/http.spec.wasm");
const module = new WebAssembly.Module(binary);
instantiate(module, {
  env: {},
  hypermode: {},
});
