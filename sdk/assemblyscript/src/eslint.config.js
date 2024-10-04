// @ts-check

import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import aseslint from "./tools/assemblyscript-eslint.js";

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  aseslint.config,
  {
    // generated
    ignores: ["transform/lib/**", "build/**"],
  },
);
