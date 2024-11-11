/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { CommentKind, CommentNode } from "assemblyscript/dist/assemblyscript.js";
import { getTypeName } from "./extractor.js";

export class ProgramInfo {
  exportFns: FunctionSignature[];
  importFns: FunctionSignature[];
  types: TypeDefinition[];
}

export class Result {
  public name?: string;
  public type: string;
}

export class Docs {
  constructor(public lines: string[]) { }
  static from(nodes: CommentNode[]): Docs | null {
    for (const node of nodes.reverse()) {
      if (node.commentKind != CommentKind.Block || !node.text.startsWith("/**")) continue;
      const lines = node.text.split("\n").filter(v => v.trim().startsWith("* ")).map(v => v.trim().slice(2));
      return new Docs(lines);
    }
    return null;
  }
}

export class FunctionSignature {
  constructor(
    public name: string,
    public parameters: Parameter[],
    public results: Result[],
    public docs: Docs | undefined = undefined
  ) { }

  toString() {
    let params = "";
    for (let i = 0; i < this.parameters.length; i++) {
      const param = this.parameters[i]!;
      const defaultValue = param.default;
      if (i > 0) params += ", ";
      params += `${param.name}: ${getTypeName(param.type)}`;
      if (defaultValue !== undefined) {
        params += ` = ${JSON.stringify(defaultValue)}`;
      }
    }
    return `${this.name}(${params}): ${getTypeName(this.results[0].type)}`;
  }

  toJSON() {
    const output = {};

    // always omit the function name

    // omit empty parameters
    if (this.parameters.length > 0) {
      output["parameters"] = this.parameters;
    }

    // omit void result types
    if (this.results[0].type !== "void") {
      output["results"] = this.results;
    }

    if (this.docs) {
      output["docs"] = this.docs;
    }

    return output;
  }
}

export class TypeDefinition {
  constructor(
    public name: string,
    public id: number,
    public fields?: Field[],
    public docs: Docs | undefined = undefined
  ) { }

  toString() {
    const name = getTypeName(this.name);
    if (!this.fields || this.fields.length === 0) {
      return name;
    }

    const fields = this.fields
      .map((f) => `${f.name}: ${getTypeName(f.type)}`)
      .join(", ");
    return `${name} { ${fields} }`;
  }

  toJSON() {
    return {
      id: this.id,
      fields: this.fields,
      docs: this.docs
    };
  }

  isHidden() {
    return this.name.startsWith("~lib/");
  }
}

export type JsonLiteral =
  | null
  | boolean
  | number
  | string
  | Array<JsonLiteral>
  | { [key: string]: JsonLiteral };

export interface Parameter {
  name: string;
  type: string;
  default?: JsonLiteral;
}

export class Field {
  constructor(
    public name: string,
    public type: string,
    public docs: Docs | undefined = undefined
  ) { }
  toJSON() {
    return {
      name: this.name,
      type: this.type,
      docs: this.docs
    }
  }
}

export const typeMap = new Map<string, string>([
  ["~lib/string/String", "string"],
  ["~lib/array/Array", "Array"],
  ["~lib/map/Map", "Map"],
  ["~lib/date/Date", "Date"],
  ["~lib/wasi_date/wasi_Date", "Date"],
]);
