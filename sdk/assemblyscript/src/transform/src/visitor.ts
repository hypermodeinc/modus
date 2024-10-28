/*
 * Copyright 2024 Jairus Tanaka. (Taken with permission from https://github.com/JairusSW/as-compiler-utils)
 * Licensed under the terms of the MIT License
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Jairus Tanaka. <me@jairus.dev>
 * SPDX-License-Identifier: MIT
 */

/* eslint-disable @typescript-eslint/no-unused-vars */
import {
    ArrayLiteralExpression,
    AssertionExpression,
    BinaryExpression,
    CallExpression,
    ElementAccessExpression,
    FloatLiteralExpression,
    FunctionTypeNode,
    IdentifierExpression,
    NamedTypeNode,
    Node,
    ObjectLiteralExpression,
    Source,
    TypeNode,
    TypeParameterNode,
    BlockStatement,
    BreakStatement,
    ClassDeclaration,
    ClassExpression,
    CommaExpression,
    ConstructorExpression,
    ContinueStatement,
    DecoratorNode,
    DoStatement,
    EmptyStatement,
    EnumDeclaration,
    EnumValueDeclaration,
    ExportDefaultStatement,
    ExportImportStatement,
    ExportMember,
    ExportStatement,
    Expression,
    ExpressionStatement,
    FalseExpression,
    FieldDeclaration,
    ForStatement,
    FunctionDeclaration,
    FunctionExpression,
    IfStatement,
    ImportDeclaration,
    ImportStatement,
    IndexSignatureNode,
    InstanceOfExpression,
    IntegerLiteralExpression,
    InterfaceDeclaration,
    LiteralExpression,
    MethodDeclaration,
    NamespaceDeclaration,
    NewExpression,
    NullExpression,
    ParameterNode,
    ParenthesizedExpression,
    PropertyAccessExpression,
    RegexpLiteralExpression,
    ReturnStatement,
    Statement,
    StringLiteralExpression,
    SuperExpression,
    SwitchCase,
    SwitchStatement,
    TemplateLiteralExpression,
    TernaryExpression,
    ThisExpression,
    ThrowStatement,
    TrueExpression,
    TryStatement,
    TypeDeclaration,
    TypeName,
    UnaryExpression,
    UnaryPostfixExpression,
    UnaryPrefixExpression,
    VariableDeclaration,
    VariableStatement,
    VoidStatement,
    WhileStatement,
    NodeKind,
    CommentNode,
    LiteralKind,
} from "assemblyscript/dist/assemblyscript.js";

export declare type Collection<T> =
    | Node
    | T[]
    | Map<string, T | T[] | Iterable<T>>
    | Iterable<T>;

export class Visitor {
    public currentSource: Source | null = null;
    public depth = 0;
    visit<T>(node: T | T[] | Map<string, T | T[] | Iterable<T>> | Iterable<T>): void {
        if (node == null) {
            return;
        } else if (Array.isArray(node)) {
            for (const element of node) {
                this.visit(element);
            }
        } else if (node instanceof Map) {
            for (const element of node.values()) {
                this.visit(element);
            }
        } else if (typeof node[Symbol.iterator] === "function") {
            // @ts-expect-error: Node can be of type [Symbol.iterator]
            for (const element of node) {
                this.visit(element);
            }
        } else {
            // @ts-expect-error: Node can be of type [Symbol.iterator]
            this._visit(node);
        }
    }
    private _visit(node: Node): void {
        if (node.kind == NodeKind.Source) {
            this.visitSource(node as Source);
        } else if (node.kind == NodeKind.NamedType) {
            this.visitNamedTypeNode(node as NamedTypeNode);
        } else if (node.kind == NodeKind.FunctionType) {
            this.visitFunctionTypeNode(node as FunctionTypeNode);
        } else if (node.kind == NodeKind.TypeName) {
            this.visitTypeName(node as TypeName);
        } else if (node.kind == NodeKind.TypeParameter) {
            this.visitTypeParameter(node as TypeParameterNode);
        } else if (node.kind == NodeKind.Identifier) {
            this.visitIdentifierExpression(node as IdentifierExpression);
        } else if (node.kind == NodeKind.Assertion) {
            this.visitAssertionExpression(node as AssertionExpression);
        } else if (node.kind == NodeKind.Binary) {
            this.visitBinaryExpression(node as BinaryExpression);
        } else if (node.kind == NodeKind.Call) {
            this.visitCallExpression(node as CallExpression);
        } else if (node.kind == NodeKind.Class) {
            this.visitClassExpression(node as ClassExpression);
        } else if (node.kind == NodeKind.Comma) {
            this.visitCommaExpression(node as CommaExpression);
        } else if (node.kind == NodeKind.Comment) {
            this.visitComment(node as CommentNode)
        } else if (node.kind == NodeKind.ElementAccess) {
            this.visitElementAccessExpression(node as ElementAccessExpression);
        } else if (node.kind == NodeKind.Function) {
            this.visitFunctionExpression(node as FunctionExpression);
        } else if (node.kind == NodeKind.InstanceOf) {
            this.visitInstanceOfExpression(node as InstanceOfExpression);
        } else if (node.kind == NodeKind.Literal) {
            this.visitLiteralExpression(node as LiteralExpression);
        } else if (node.kind == NodeKind.New) {
            this.visitNewExpression(node as NewExpression);
        } else if (node.kind == NodeKind.Parenthesized) {
            this.visitParenthesizedExpression(node as ParenthesizedExpression);
        } else if (node.kind == NodeKind.PropertyAccess) {
            this.visitPropertyAccessExpression(node as PropertyAccessExpression);
        } else if (node.kind == NodeKind.Ternary) {
            this.visitTernaryExpression(node as TernaryExpression);
        } else if (node.kind == NodeKind.UnaryPostfix) {
            this.visitUnaryPostfixExpression(node as UnaryPostfixExpression);
        } else if (node.kind == NodeKind.UnaryPrefix) {
            this.visitUnaryPrefixExpression(node as UnaryPrefixExpression);
        } else if (node.kind == NodeKind.Block) {
            this.visitBlockStatement(node as BlockStatement);
        } else if (node.kind == NodeKind.Break) {
            this.visitBreakStatement(node as BreakStatement);
        } else if (node.kind == NodeKind.Continue) {
            this.visitContinueStatement(node as ContinueStatement);
        } else if (node.kind == NodeKind.Do) {
            this.visitDoStatement(node as DoStatement);
        } else if (node.kind == NodeKind.Empty) {
            this.visitEmptyStatement(node as EmptyStatement);
        } else if (node.kind == NodeKind.Export) {
            this.visitExportStatement(node as ExportStatement);
        } else if (node.kind == NodeKind.ExportDefault) {
            this.visitExportDefaultStatement(node as ExportDefaultStatement);
        } else if (node.kind == NodeKind.ExportImport) {
            this.visitExportImportStatement(node as ExportImportStatement);
        } else if (node.kind == NodeKind.Expression) {
            this.visitExpressionStatement(node as ExpressionStatement);
        } else if (node.kind == NodeKind.For) {
            this.visitForStatement(node as ForStatement);
        } else if (node.kind == NodeKind.If) {
            this.visitIfStatement(node as IfStatement);
        } else if (node.kind == NodeKind.Import) {
            this.visitImportStatement(node as ImportStatement);
        } else if (node.kind == NodeKind.Return) {
            this.visitReturnStatement(node as ReturnStatement);
        } else if (node.kind == NodeKind.Switch) {
            this.visitSwitchStatement(node as SwitchStatement);
        } else if (node.kind == NodeKind.Throw) {
            this.visitThrowStatement(node as ThrowStatement);
        } else if (node.kind == NodeKind.Try) {
            this.visitTryStatement(node as TryStatement);
        } else if (node.kind == NodeKind.Variable) {
            this.visitVariableStatement(node as VariableStatement);
        } else if (node.kind == NodeKind.While) {
            this.visitWhileStatement(node as WhileStatement);
        } else if (node.kind == NodeKind.ClassDeclaration) {
            this.visitClassDeclaration(node as ClassDeclaration);
        } else if (node.kind == NodeKind.EnumDeclaration) {
            this.visitEnumDeclaration(node as EnumDeclaration);
        } else if (node.kind == NodeKind.EnumValueDeclaration) {
            this.visitEnumValueDeclaration(node as EnumValueDeclaration);
        } else if (node.kind == NodeKind.FieldDeclaration) {
            this.visitFieldDeclaration(node as FieldDeclaration);
        } else if (node.kind == NodeKind.FunctionDeclaration) {
            this.visitFunctionDeclaration(node as FunctionDeclaration);
        } else if (node.kind == NodeKind.ImportDeclaration) {
            this.visitImportDeclaration(node as ImportDeclaration);
        } else if (node.kind == NodeKind.InterfaceDeclaration) {
            this.visitInterfaceDeclaration(node as InterfaceDeclaration);
        } else if (node.kind == NodeKind.MethodDeclaration) {
            this.visitMethodDeclaration(node as MethodDeclaration);
        } else if (node.kind == NodeKind.NamespaceDeclaration) {
            this.visitNamespaceDeclaration(node as NamespaceDeclaration);
        } else if (node.kind == NodeKind.TypeDeclaration) {
            this.visitTypeDeclaration(node as TypeDeclaration);
        } else if (node.kind == NodeKind.VariableDeclaration) {
            this.visitVariableDeclaration(node as VariableDeclaration);
        } else if (node.kind == NodeKind.Decorator) {
            this.visitDecoratorNode(node as DecoratorNode);
        } else if (node.kind == NodeKind.ExportMember) {
            this.visitExportMember(node as ExportMember);
        } else if (node.kind == NodeKind.SwitchCase) {
            this.visitSwitchCase(node as SwitchCase);
        } else if (node.kind == NodeKind.IndexSignature) {
            this.visitIndexSignature(node as IndexSignatureNode);
        } else if (node.kind == NodeKind.Null) {
            this.visitNullExpression(node as NullExpression);
        } else if (node.kind == NodeKind.Constructor) {
            this.visitConstructorExpression(node as ConstructorExpression);
        } else if (node.kind == NodeKind.Super) {
            this.visitSuperExpression(node as SuperExpression);
        } else if (node.kind == NodeKind.True) {
            this.visitTrueExpression(node as TrueExpression);
        } else if (node.kind == NodeKind.False) {
            this.visitFalseExpression(node as FalseExpression);
        } else if (node.kind == NodeKind.This) {
            this.visitThisExpression(node as ThisExpression);
        } else if (node.kind == NodeKind.Void) {
            this.visitVoidStatement(node as VoidStatement);
        }
    }
    visitSource(node: Source): void {
        this.currentSource = node;
        for (const stmt of node.statements) {
            this.depth++;
            this._visit(stmt);
            this.depth--;
        }
        this.currentSource = null;
    }
    visitTypeNode(node: TypeNode): void { }
    visitTypeName(node: TypeName): void {
        this.visit(node.identifier);
        this.visit(node.next);
    }
    visitNamedTypeNode(node: NamedTypeNode): void {
        this.visit(node.name);
        this.visit(node.typeArguments);
    }
    visitFunctionTypeNode(node: FunctionTypeNode): void {
        this.visit(node.parameters);
        this.visit(node.returnType);
        this.visit(node.explicitThisType);
    }
    visitTypeParameter(node: TypeParameterNode): void {
        this.visit(node.name);
        this.visit(node.extendsType);
        this.visit(node.defaultType);
    }
    visitIdentifierExpression(node: IdentifierExpression): void { }
    visitArrayLiteralExpression(node: ArrayLiteralExpression) {
        this.visit(node.elementExpressions);
    }
    visitObjectLiteralExpression(node: ObjectLiteralExpression) {
        this.visit(node.names);
        this.visit(node.values);
    }
    visitAssertionExpression(node: AssertionExpression) {
        this.visit(node.toType);
        this.visit(node.expression);
    }
    visitBinaryExpression(node: BinaryExpression) {
        this.visit(node.left);
        this.visit(node.right);
    }
    visitCallExpression(node: CallExpression) {
        this.visit(node.expression);
        this.visitArguments(node.typeArguments, node.args);
    }
    visitArguments(typeArguments: TypeNode[] | null, args: Expression[]) {
        this.visit(typeArguments);
        this.visit(args);
    }
    visitClassExpression(node: ClassExpression) {
        this.visit(node.declaration);
    }
    visitCommaExpression(node: CommaExpression) {
        this.visit(node.expressions);
    }
    visitElementAccessExpression(node: ElementAccessExpression) {
        this.visit(node.elementExpression);
        this.visit(node.expression);
    }
    visitFunctionExpression(node: FunctionExpression) {
        this.visit(node.declaration);
    }
    visitLiteralExpression(node: LiteralExpression) {
        if (node.literalKind == LiteralKind.Float) {
            this.visitFloatLiteralExpression(node as FloatLiteralExpression);
        } else if (node.literalKind == LiteralKind.Integer) {
            this.visitIntegerLiteralExpression(node as IntegerLiteralExpression);
        } else if (node.literalKind == LiteralKind.String) {
            this.visitStringLiteralExpression(node as StringLiteralExpression);
        } else if (node.literalKind == LiteralKind.Template) {
            this.visitTemplateLiteralExpression(node as TemplateLiteralExpression);
        } else if (node.literalKind == LiteralKind.RegExp) {
            this.visitRegexpLiteralExpression(node as RegexpLiteralExpression);
        } else if (node.literalKind == LiteralKind.Array) {
            this.visitArrayLiteralExpression(node as ArrayLiteralExpression);
        } else if (node.literalKind == LiteralKind.Object) {
            this.visitObjectLiteralExpression(node as ObjectLiteralExpression);
        }
    }
    visitFloatLiteralExpression(node: FloatLiteralExpression) { }
    visitInstanceOfExpression(node: InstanceOfExpression) {
        this.visit(node.expression);
        this.visit(node.isType);
    }
    visitIntegerLiteralExpression(node: IntegerLiteralExpression) { }
    visitStringLiteral(str: string, singleQuoted: boolean = false) { }
    visitStringLiteralExpression(node: StringLiteralExpression) {
        this.visitStringLiteral(node.value);
    }
    visitTemplateLiteralExpression(node: TemplateLiteralExpression) { }
    visitRegexpLiteralExpression(node: RegexpLiteralExpression) { }
    visitNewExpression(node: NewExpression) {
        this.visit(node.typeArguments);
        this.visitArguments(node.typeArguments, node.args);
        this.visit(node.args);
    }
    visitParenthesizedExpression(node: ParenthesizedExpression) {
        this.visit(node.expression);
    }
    visitPropertyAccessExpression(node: PropertyAccessExpression) {
        this.visit(node.property);
        this.visit(node.expression);
    }
    visitTernaryExpression(node: TernaryExpression) {
        this.visit(node.condition);
        this.visit(node.ifThen);
        this.visit(node.ifElse);
    }
    visitUnaryExpression(node: UnaryExpression) {
        this.visit(node.operand);
    }
    visitUnaryPostfixExpression(node: UnaryPostfixExpression) {
        this.visit(node.operand);
    }
    visitUnaryPrefixExpression(node: UnaryPrefixExpression) {
        this.visit(node.operand);
    }
    visitSuperExpression(node: SuperExpression) { }
    visitFalseExpression(node: FalseExpression) { }
    visitTrueExpression(node: TrueExpression) { }
    visitThisExpression(node: ThisExpression) { }
    visitNullExpression(node: NullExpression) { }
    visitConstructorExpression(node: ConstructorExpression) { }
    visitNodeAndTerminate(statement: Statement) { }
    visitBlockStatement(node: BlockStatement) {
        this.depth++;
        this.visit(node.statements);
        this.depth--;
    }
    visitBreakStatement(node: BreakStatement) {
        this.visit(node.label);
    }
    visitContinueStatement(node: ContinueStatement) {
        this.visit(node.label);
    }
    visitClassDeclaration(node: ClassDeclaration, isDefault: boolean = false) {
        this.visit(node.name);
        this.depth++;
        this.visit(node.decorators);
        if (
            node.isGeneric ? node.typeParameters != null : node.typeParameters == null
        ) {
            this.visit(node.typeParameters);
            this.visit(node.extendsType);
            this.visit(node.implementsTypes);
            this.visit(node.members);
            this.depth--;
        } else {
            throw new Error(
                "Expected to type parameters to match class declaration, but found type mismatch instead!",
            );
        }
    }
    visitDoStatement(node: DoStatement) {
        this.visit(node.condition);
        this.visit(node.body);
    }
    visitEmptyStatement(node: EmptyStatement) { }
    visitEnumDeclaration(node: EnumDeclaration, isDefault: boolean = false) {
        this.visit(node.name);
        this.visit(node.decorators);
        this.visit(node.values);
    }
    visitEnumValueDeclaration(node: EnumValueDeclaration) {
        this.visit(node.name);
        this.visit(node.initializer);
    }
    visitExportImportStatement(node: ExportImportStatement) {
        this.visit(node.name);
        this.visit(node.externalName);
    }
    visitExportMember(node: ExportMember) {
        this.visit(node.localName);
        this.visit(node.exportedName);
    }
    visitExportStatement(node: ExportStatement) {
        this.visit(node.path);
        this.visit(node.members);
    }
    visitExportDefaultStatement(node: ExportDefaultStatement) {
        this.visit(node.declaration);
    }
    visitExpressionStatement(node: ExpressionStatement) {
        this.visit(node.expression);
    }
    visitFieldDeclaration(node: FieldDeclaration) {
        this.visit(node.name);
        this.visit(node.type);
        this.visit(node.initializer);
        this.visit(node.decorators);
    }
    visitForStatement(node: ForStatement) {
        this.visit(node.initializer);
        this.visit(node.condition);
        this.visit(node.incrementor);
        this.visit(node.body);
    }
    visitFunctionDeclaration(
        node: FunctionDeclaration,
        isDefault: boolean = false,
    ) {
        this.visit(node.name);
        this.visit(node.decorators);
        this.visit(node.typeParameters);
        this.visit(node.signature);
        this.depth++;
        this.visit(node.body);
        this.depth--;
    }
    visitIfStatement(node: IfStatement) {
        this.visit(node.condition);
        this.visit(node.ifTrue);
        this.visit(node.ifFalse);
    }
    visitImportDeclaration(node: ImportDeclaration) {
        this.visit(node.foreignName);
        this.visit(node.name);
        this.visit(node.decorators);
    }
    visitImportStatement(node: ImportStatement) {
        this.visit(node.namespaceName);
        this.visit(node.declarations);
    }
    visitIndexSignature(node: IndexSignatureNode) {
        this.visit(node.keyType);
        this.visit(node.valueType);
    }
    visitInterfaceDeclaration(
        node: InterfaceDeclaration,
        isDefault: boolean = false,
    ) {
        this.visit(node.name);
        this.visit(node.typeParameters);
        this.visit(node.implementsTypes);
        this.visit(node.extendsType);
        this.depth++;
        this.visit(node.members);
        this.depth--;
    }
    visitMethodDeclaration(node: MethodDeclaration) {
        this.visit(node.name);
        this.visit(node.typeParameters);
        this.visit(node.signature);
        this.visit(node.decorators);
        this.depth++;
        this.visit(node.body);
        this.depth--;
    }
    visitNamespaceDeclaration(
        node: NamespaceDeclaration,
        isDefault: boolean = false,
    ) {
        this.visit(node.name);
        this.visit(node.decorators);
        this.visit(node.members);
    }
    visitReturnStatement(node: ReturnStatement) {
        this.visit(node.value);
    }
    visitSwitchCase(node: SwitchCase) {
        this.visit(node.label);
        this.visit(node.statements);
    }
    visitSwitchStatement(node: SwitchStatement) {
        this.visit(node.condition);
        this.depth++;
        this.visit(node.cases);
        this.depth--;
    }
    visitThrowStatement(node: ThrowStatement) {
        this.visit(node.value);
    }
    visitTryStatement(node: TryStatement) {
        this.visit(node.bodyStatements);
        this.visit(node.catchVariable);
        this.visit(node.catchStatements);
        this.visit(node.finallyStatements);
    }
    visitTypeDeclaration(node: TypeDeclaration) {
        this.visit(node.name);
        this.visit(node.decorators);
        this.visit(node.type);
        this.visit(node.typeParameters);
    }
    visitVariableDeclaration(node: VariableDeclaration) {
        this.visit(node.name);
        this.visit(node.type);
        this.visit(node.initializer);
    }
    visitVariableStatement(node: VariableStatement) {
        this.visit(node.decorators);
        this.visit(node.declarations);
    }
    visitWhileStatement(node: WhileStatement) {
        this.visit(node.condition);
        this.depth++;
        this.visit(node.body);
        this.depth--;
    }
    visitVoidStatement(node: VoidStatement) { }
    visitComment(node: CommentNode) { }
    visitDecoratorNode(node: DecoratorNode) {
        this.visit(node.name);
        this.visit(node.args);
    }
    visitParameter(node: ParameterNode) {
        this.visit(node.name);
        this.visit(node.implicitFieldDeclaration);
        this.visit(node.initializer);
        this.visit(node.type);
    }
}