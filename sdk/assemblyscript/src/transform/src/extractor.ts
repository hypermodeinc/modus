/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import binaryen from "assemblyscript/lib/binaryen.js";
import {
  ArrayLiteralExpression,
  Class,
  ElementKind,
  Expression,
  FloatLiteralExpression,
  Function as Func,
  FunctionDeclaration,
  IntegerLiteralExpression,
  LiteralExpression,
  LiteralKind,
  NodeKind,
  Program,
  Property,
  StringLiteralExpression,
} from "assemblyscript/dist/assemblyscript.js";
import {
  FunctionSignature,
  JsonLiteral,
  Parameter,
  ProgramInfo,
  TypeDefinition,
  typeMap,
} from "./types.js";
import ModusTransform from "./index.js";

export class Extractor {
  binaryen: typeof binaryen;
  module: binaryen.Module;
  program: Program;
  transform: ModusTransform;

  constructor(transform: ModusTransform, module: binaryen.Module) {
    this.program = transform.program;
    this.binaryen = transform.binaryen;
    this.module = module;
  }

  getProgramInfo(): ProgramInfo {
    const exportedFunctions = this.getExportedFunctions()
      .map((e) => this.convertToFunctionSignature(e))
      .sort((a, b) => a.name.localeCompare(b.name));

    const importedFunctions = this.getImportedFunctions()
      .map((e) => this.convertToFunctionSignature(e))
      .sort((a, b) => a.name.localeCompare(b.name));

    const allTypes = new Map<string, TypeDefinition>(
      Array.from(this.program.managedClasses.values())
        .filter((c) => c.id > 2) // skip built-in classes
        .map((c) => {
          return new TypeDefinition(
            c.type.toString(),
            c.id,
            this.getClassFields(c),
          );
        })
        .map((t) => [t.name, t]),
    );
    const typePathsUsed = new Set(
      exportedFunctions
        .concat(importedFunctions)
        .flatMap((f) =>
          f.parameters.map((p) => p.type).concat(f.results[0]?.type),
        )
        .map((p) => makeNonNullable(p)),
    );

    const typesUsed = new Map<string, TypeDefinition>();

    allTypes.forEach((t) => {
      if (typePathsUsed.has(t.name)) {
        typesUsed.set(t.name, t);
      }
    });

    typesUsed.forEach((t) => {
      this.expandDependentTypes(t, allTypes, typesUsed);
    });

    const types = Array.from(typesUsed.values()).sort((a, b) =>
      a.name.localeCompare(b.name),
    );

    return {
      exportFns: exportedFunctions,
      importFns: importedFunctions,
      types,
    };
  }

  private expandDependentTypes(
    type: TypeDefinition,
    allTypes: Map<string, TypeDefinition>,
    typesUsed: Map<string, TypeDefinition>,
  ) {
    // collect dependent types into this set
    const dependentTypes = new Set<TypeDefinition>();

    // include fields
    if (type.fields) {
      type.fields.forEach((f) => {
        const path = makeNonNullable(f.type);
        const typeDef = allTypes.get(path);
        if (typeDef) {
          dependentTypes.add(typeDef);
        }
      });
    }

    // include generic type arguments
    const cls = this.program.managedClasses.get(type.id);
    if (cls.typeArguments) {
      cls.typeArguments.forEach((t) => {
        const typeDef = allTypes.get(t.toString());
        if (typeDef) {
          dependentTypes.add(typeDef);
        }
      });
    }

    // recursively expand dependencies of dependent types
    dependentTypes.forEach((t) => {
      if (!typesUsed.has(t.name)) {
        typesUsed.set(t.name, t);
        this.expandDependentTypes(t, allTypes, typesUsed);
      }
    });
  }

  private getClassFields(c: Class) {
    if (
      c.isArrayLike ||
      typeMap.has(c.type.toString()) ||
      typeMap.has(c.prototype.internalName)
    ) {
      return undefined;
    }

    return Array.from(c.members.values())
      .filter((m) => m.kind === ElementKind.PropertyPrototype)
      .map((m) => {
        const instance = this.program.instancesByName.get(m.internalName);
        return instance as Property;
      })
      .filter((p) => p && p.isField)
      .map((f) => ({
        name: f.name,
        type: f.type.toString(),
      }));
  }

  private getExportedFunctions() {
    const results: importExportInfo[] = [];

    for (let i = 0; i < this.module.getNumExports(); ++i) {
      const ref = this.module.getExportByIndex(i);
      const info = this.binaryen.getExportInfo(ref);

      if (info.kind !== binaryen.ExternalFunction) {
        continue;
      }

      const exportName = info.name;
      if (exportName.startsWith("_")) {
        continue;
      }

      const functionName = info.value.replace(/^export:/, "");
      results.push({ name: exportName, function: functionName });
    }

    return results;
  }

  private getImportedFunctions() {
    const results: importExportInfo[] = [];

    this.program.moduleImports.forEach((module, modName) => {
      module.forEach((e, fnName) => {
        if (modName != "env" && !modName.startsWith("wasi")) {
          results.push({
            name: `${modName}.${fnName}`,
            function: e.internalName,
          });
        }
      });
    });

    return results;
  }

  private convertToFunctionSignature(e: importExportInfo): FunctionSignature {
    const f = this.program.instancesByName.get(e.function) as Func;
    const d = f.declaration as FunctionDeclaration;
    const params: Parameter[] = [];
    for (let i = 0; i < f.signature.parameterTypes.length; i++) {
      const param = d.signature.parameters[i];
      const type = f.signature.parameterTypes[i];
      const name = param.name.text;
      const defaultValue = getLiteral(param.initializer);
      params.push({
        name,
        type: type.toString(),
        default: defaultValue,
      });
    }

    return new FunctionSignature(e.name, params, [
      { type: f.signature.returnType.toString() },
    ]);
  }
}

interface importExportInfo {
  name: string;
  function: string;
}

export function getTypeName(path: string): string {
  if (path.startsWith("~lib/array/Array")) {
    const type = getTypeName(
      path.slice(path.indexOf("<") + 1, path.lastIndexOf(">")),
    );
    if (isNullable(type)) {
      return "(" + type + ")[]";
    }
    return type + "[]";
  }

  if (isNullable(path)) {
    return makeNullable(getTypeName(makeNonNullable(path)));
  }

  const name = typeMap.get(path);
  if (name) return name;

  if (path.startsWith("~lib/map/Map")) {
    const [keyType, valueType] = getMapSubtypes(path);
    return "Map<" + getTypeName(keyType) + ", " + getTypeName(valueType) + ">";
  }

  if (path.startsWith("~lib/@hypermode")) {
    const lastIndex = path.lastIndexOf("/");
    const module = path.slice(
      path.lastIndexOf("/", lastIndex - 1) + 1,
      lastIndex,
    );
    const ty = path.slice(lastIndex + 1);
    return module + "." + ty;
  }

  return path.slice(path.lastIndexOf("/") + 1);
}

export function getLiteral(node: Expression | null): JsonLiteral {
  if (!node) return undefined;
  switch (node.kind) {
    case NodeKind.True: {
      return true;
    }
    case NodeKind.False: {
      return false;
    }
    case NodeKind.Null: {
      return null;
    }
    case NodeKind.Literal: {
      const _node = node as LiteralExpression;
      switch (_node.literalKind) {
        case LiteralKind.Integer: {
          return i64_to_f64((_node as IntegerLiteralExpression).value);
        }
        case LiteralKind.Float: {
          return (_node as FloatLiteralExpression).value;
        }
        case LiteralKind.String: {
          return (_node as StringLiteralExpression).value;
        }
        case LiteralKind.Array: {
          const out = [];
          const literals = (_node as ArrayLiteralExpression).elementExpressions;
          for (let i = 0; i < literals.length; i++) {
            const lit = getLiteral(literals[i]);
            out.push(lit);
          }
          return out;
        }
      }
    }
  }
  return "";
}

const nullableTypeRegex = /\s?\|\s?null$/;

function isNullable(type: string) {
  return nullableTypeRegex.test(type);
}

function makeNonNullable(type?: string) {
  return type?.replace(nullableTypeRegex, "");
}

function makeNullable(type?: string) {
  if (isNullable(type)) {
    return type;
  } else {
    return type + " | null";
  }
}

function getMapSubtypes(type: string): [string, string] {
  const prefix = "~lib/map/Map<";
  if (!type.startsWith(prefix)) {
    return ["", ""];
  }

  let n = 1;
  let c = 0;
  for (let i = prefix.length; i < type.length; i++) {
    switch (type.charAt(i)) {
      case "<":
        n++;
        break;
      case ",":
        if (n == 1) {
          c = i;
        }
        break;
      case ">":
        n--;
        if (n == 0) {
          return [type.slice(prefix.length, c), type.slice(c + 1, i)];
        }
    }
  }

  return ["", ""];
}
