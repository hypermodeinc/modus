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
  ClassDeclaration,
  CommentKind,
  CommentNode,
  ElementKind,
  Expression,
  FieldDeclaration,
  FloatLiteralExpression,
  Function as Func,
  FunctionDeclaration,
  IntegerLiteralExpression,
  LiteralExpression,
  LiteralKind,
  NodeKind,
  Program,
  Property,
  Range,
  StringLiteralExpression,
} from "assemblyscript/dist/assemblyscript.js";
import {
  Docs,
  Field,
  FunctionSignature,
  JsonLiteral,
  Parameter,
  ProgramInfo,
  TypeDefinition,
  typeMap,
} from "./types.js";
import ModusTransform from "./index.js";
import { Visitor } from "./visitor.js";

export class Extractor {
  binaryen: typeof binaryen;
  module: binaryen.Module;
  program: Program;
  transform: ModusTransform;

  constructor(transform: ModusTransform) {
    this.binaryen = transform.binaryen;
  }

  initHook(program: Program): void {
    this.program = program;
  }

  compileHook(module: binaryen.Module): void {
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
            this.getClassFields(c)
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
      types: types.map(v => this.getTypeDocs(v)),
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
      .map((f) => <Field>{
        name: f.name,
        type: f.type.toString()
      });
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

    const signature = new FunctionSignature(e.name, params, [
      { type: f.signature.returnType.toString() },
    ]);

    signature.docs = this.getDocsFromFunction(signature);
    return signature;
  }
  private getDocsFromFunction(signature: FunctionSignature) {
    const visitor = new Visitor();
    let docs: Docs | null = null;

    visitor.visitFunctionDeclaration = (node: FunctionDeclaration) => {
      const source = node.range.source;
      // Exported/Imported name may differ from real defined name
      // TODO: Track renaming and aliasing of function identifiers

      if (node.name.text == signature.name) {
        const nodeIndex = source.statements.indexOf(node);
        const prevNode = source.statements[Math.max(nodeIndex - 1, 0)];

        const start = nodeIndex > 0 ? prevNode.range.end : 0;
        const end = node.range.start;

        const newRange = new Range(start, end);
        newRange.source = source;
        const commentNodes = this.parseComments(newRange);
        if (!commentNodes.length) return;
        docs = Docs.from(commentNodes);
      }
    }
    visitor.visit(this.program.sources);
    return docs;
  }
  private getTypeDocs(type: TypeDefinition): TypeDefinition {
    const name = (() => {
      if (type.name.startsWith("~lib/")) return null;
      return type.name.slice(
        Math.max(
          type.name.lastIndexOf("<"),
          type.name.lastIndexOf("/") + 1
        ),
        Math.max(
          type.name.indexOf(">"),
          type.name.length
        )
      );
    })();
    if (!name) return type;
    for (const _node of Array.from(this.program.managedClasses.values())) {
      if (_node.name != name) continue;
      const node = _node.declaration as ClassDeclaration;
      const source = node.range.source;
      const nodeIndex = source.statements.indexOf(node);
      const prevNode = source.statements[Math.max(nodeIndex - 1, 0)];

      const start = nodeIndex > 0 ? prevNode.range.end : 0;
      const end = node.range.start;

      const newRange = new Range(start, end);
      newRange.source = source;
      const commentNodes = this.parseComments(newRange);
      if (!commentNodes.length) break;
      type.docs = Docs.from(commentNodes);

      if (node.members.length) {
        const memberDocs = this.getFieldsDocs(node.members.filter(v => v.kind == NodeKind.FieldDeclaration) as FieldDeclaration[], node);
        if (!memberDocs) continue;
        for (let i = 0; i < memberDocs.length; i++) {
          const docs = memberDocs[i];
          if (docs) {
            console.log("Got docs: ", docs.description);
            const index = type.fields.findIndex(v => v.name == node.members[i].name.text);
            if (index < 0) continue;
            type.fields[index].docs = docs;
          }
        }
      }
    }
    return type;
  }
  private getFieldsDocs(nodes: FieldDeclaration[], parent: ClassDeclaration): (Docs | null)[] {
    const docs = new Array<Docs | null>(nodes.length).fill(null);
    console.log("Members: " + nodes.length)
    for (let i = 0; i < nodes.length; i++) {
      const node = nodes[i];
      const source = node.range.source;

      const start = i == 0 ? parent.range.start : nodes[i-1].range.end;
      const end = node.range.start;

      console.log("text: " + source.text.slice(start, end))
      const newRange = new Range(start, end);
      newRange.source = source;
      const commentNodes = this.parseComments(newRange);
      if (!commentNodes.length) continue;
      docs[i] = Docs.from(commentNodes);
    }
    return docs;
  }
  private parseComments(range: Range): CommentNode[] {
    const nodes: CommentNode[] = [];
    let text = range.source.text.slice(range.start, range.end).trim();
    const start = Math.min(
      text.indexOf("/*") === -1 ? Infinity : text.indexOf("/*"),
      text.indexOf("//") === -1 ? Infinity : text.indexOf("//")
    );
    if (start !== Infinity) text = text.slice(start);
    let commentKind: CommentKind;

    if (text.startsWith("//")) {
      commentKind = text.startsWith("///") ? CommentKind.Triple : CommentKind.Line;

      const end = range.source.text.indexOf("\n", range.start + 1);
      if (end === -1) return [];
      range.start = range.source.text.indexOf("//", range.start);
      const newRange = new Range(range.start, end);
      newRange.source = range.source;
      const node = new CommentNode(commentKind, newRange.source.text.slice(newRange.start, newRange.end), newRange);

      nodes.push(node);

      if (end < range.end) {
        const newRange = new Range(end, range.end);
        newRange.source = range.source;
        nodes.push(...this.parseComments(newRange));
      }
    } else if (text.startsWith("/*")) {
      commentKind = CommentKind.Block;
      const end = range.source.text.indexOf("*/", range.start) + 2;
      if (end === 1) return [];

      range.start = range.source.text.indexOf("/**", range.start);
      const newRange = new Range(range.start, end);
      newRange.source = range.source;
      const node = new CommentNode(commentKind, newRange.source.text.slice(newRange.start, newRange.end), newRange);

      nodes.push(node);

      if (end < range.end) {
        const newRange = new Range(end, range.end);
        newRange.source = range.source;
        nodes.push(...this.parseComments(newRange));
      }
    } else {
      return [];
    }

    return nodes;
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
