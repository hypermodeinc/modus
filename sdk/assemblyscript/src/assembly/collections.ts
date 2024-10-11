/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as utils from "./utils";

export type CollectionStatus = string;
// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace CollectionStatus {
  export const Success = "success";
  export const Error = "error";
}
abstract class CollectionResult {
  collection: string;
  status: CollectionStatus;
  error: string;
  get isSuccessful(): bool {
    return this.status == CollectionStatus.Success;
  }

  constructor(collection: string, status: CollectionStatus, error: string) {
    this.collection = collection;
    this.status = status;
    this.error = error;
  }
}
export class CollectionMutationResult extends CollectionResult {
  operation: string;
  keys: string[] = [];

  constructor(
    collection: string,
    status: CollectionStatus,
    error: string,
    operation: string,
  ) {
    super(collection, status, error);
    this.operation = operation;
  }
}
export class SearchMethodMutationResult extends CollectionResult {
  operation: string;
  searchMethod: string;

  constructor(
    collection: string,
    status: CollectionStatus,
    error: string,
    operation: string,
    searchMethod: string,
  ) {
    super(collection, status, error);
    this.operation = operation;
    this.searchMethod = searchMethod;
  }
}
export class CollectionSearchResult extends CollectionResult {
  searchMethod: string;
  objects: CollectionSearchResultObject[];

  constructor(
    collection: string,
    status: CollectionStatus,
    error: string,
    searchMethod: string,
    objects: CollectionSearchResultObject[],
  ) {
    super(collection, status, error);
    this.searchMethod = searchMethod;
    this.objects = objects;
  }
}

export class CollectionSearchResultObject {
  namespace: string;
  key: string;
  text: string;
  labels: string[];
  distance: f64;
  score: f64;

  constructor(
    namespace: string,
    key: string,
    text: string,
    labels: string[],
    distance: f64,
    score: f64,
  ) {
    this.namespace = namespace;
    this.key = key;
    this.text = text;
    this.labels = labels;
    this.distance = distance;
    this.score = score;
  }
}

export class CollectionClassificationResult extends CollectionResult {
  searchMethod: string;
  labelsResult: CollectionClassificationLabelObject[];
  cluster: CollectionClassificationResultObject[];

  constructor(
    collection: string,
    status: CollectionStatus,
    error: string,
    searchMethod: string,
    labelsResult: CollectionClassificationLabelObject[],
    cluster: CollectionClassificationResultObject[],
  ) {
    super(collection, status, error);
    this.searchMethod = searchMethod;
    this.labelsResult = labelsResult;
    this.cluster = cluster;
  }
}

export class CollectionClassificationLabelObject {
  label: string;
  confidence: f64;

  constructor(label: string, confidence: f64) {
    this.label = label;
    this.confidence = confidence;
  }
}

export class CollectionClassificationResultObject {
  key: string;
  labels: string[];
  distance: f64;
  score: f64;

  constructor(key: string, labels: string[], distance: f64, score: f64) {
    this.key = key;
    this.labels = labels;
    this.distance = distance;
    this.score = score;
  }
}

// @ts-expect-error: decorator
@external("modus_collections", "upsert")
declare function hostUpsertToCollection(
  collection: string,
  namespace: string,
  keys: string[],
  texts: string[],
  labels: string[][],
): CollectionMutationResult;

// @ts-expect-error: decorator
@external("modus_collections", "delete")
declare function hostDeleteFromCollection(
  collection: string,
  namespace: string,
  key: string,
): CollectionMutationResult;

// @ts-expect-error: decorator
@external("modus_collections", "search")
declare function hostSearchCollection(
  collection: string,
  namespaces: string[],
  searchMethod: string,
  text: string,
  limit: i32,
  returnText: bool,
): CollectionSearchResult;

// @ts-expect-error: decorator
@external("modus_collections", "classifyText")
declare function hostNnClassifyCollection(
  collection: string,
  namespace: string,
  searchMethod: string,
  text: string,
): CollectionClassificationResult;

// @ts-expect-error: decorator
@external("modus_collections", "recomputeIndex")
declare function hostRecomputeSearchMethod(
  collection: string,
  namespace: string,
  searchMethod: string,
): SearchMethodMutationResult;

// @ts-expect-error: decorator
@external("modus_collections", "computeDistance")
declare function hostComputeDistance(
  collection: string,
  namespace: string,
  searchMethod: string,
  key1: string,
  key2: string,
): CollectionSearchResultObject;

// @ts-expect-error: decorator
@external("modus_collections", "getText")
declare function hostGetTextFromCollection(
  collection: string,
  namespace: string,
  key: string,
): string;

// @ts-expect-error: decorator
@external("modus_collections", "dumpTexts")
declare function hostGetTextsFromCollection(
  collection: string,
  namespace: string,
): Map<string, string>;

// @ts-expect-error: decorator
@external("modus_collections", "getNamespaces")
declare function hostGetNamespacesFromCollection(collection: string): string[];

// @ts-expect-error: decorator
@external("modus_collections", "getVector")
declare function hostGetVector(
  collection: string,
  namespace: string,
  searchMethod: string,
  key: string,
): f32[];

// @ts-expect-error: decorator
@external("modus_collections", "getLabels")
declare function hostGetLabels(
  collection: string,
  namespace: string,
  key: string,
): string[];

// @ts-expect-error: decorator
@external("modus_collections", "searchByVector")
declare function hostSearchCollectionByVector(
  collection: string,
  namespaces: string[],
  searchMethod: string,
  vector: f32[],
  limit: i32,
  returnText: bool,
): CollectionSearchResult;

// add batch upsert
export function upsertBatch(
  collection: string,
  keys: string[] | null,
  texts: string[],
  labelsArr: string[][] = [],
  namespace: string = "",
): CollectionMutationResult {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Collection is empty.",
      "upsert",
    );
  }
  if (texts.length == 0) {
    console.error("Texts is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Texts is empty.",
      "upsert",
    );
  }
  let keysArr: string[] = [];
  if (keys != null) {
    keysArr = keys;
  }

  const result = hostUpsertToCollection(
    collection,
    namespace,
    keysArr,
    texts,
    labelsArr,
  );
  if (utils.resultIsInvalid(result)) {
    console.error("Error upserting to Text index.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Error upserting to Text index.",
      "upsert",
    );
  }
  return result;
}

// add data to in-mem storage, get all embedders for a collection, run text through it
// and insert the Text into the Text indexes for each search method
export function upsert(
  collection: string,
  key: string | null,
  text: string,
  labels: string[] = [],
  namespace: string = "",
): CollectionMutationResult {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Collection is empty.",
      "upsert",
    );
  }
  if (text.length == 0) {
    console.error("Text is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Text is empty.",
      "upsert",
    );
  }
  const keys: string[] = [];
  if (key != null) {
    keys.push(key);
  }

  const texts: string[] = [text];

  const labelsArr: string[][] = [labels];

  const result = hostUpsertToCollection(
    collection,
    namespace,
    keys,
    texts,
    labelsArr,
  );
  if (utils.resultIsInvalid(result)) {
    console.error("Error upserting to Text index.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Error upserting to Text index.",
      "upsert",
    );
  }
  return result;
}

// remove data from in-mem storage and indexes
export function remove(
  collection: string,
  key: string,
  namespace: string = "",
): CollectionMutationResult {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Collection is empty.",
      "delete",
    );
  }
  if (key.length == 0) {
    console.error("Key is empty.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Key is empty.",
      "delete",
    );
  }
  const result = hostDeleteFromCollection(collection, namespace, key);
  if (utils.resultIsInvalid(result)) {
    console.error("Error deleting from Text index.");
    return new CollectionMutationResult(
      collection,
      CollectionStatus.Error,
      "Error deleting from Text index.",
      "delete",
    );
  }
  return result;
}

// fetch embedders for collection & search method, run text through it and
// search Text index for similar Texts, return the result keys
// open question: how do i return a more expansive result from string array
export function search(
  collection: string,
  searchMethod: string,
  text: string,
  limit: i32,
  returnText: bool = false,
  namespaces: string[] = [],
): CollectionSearchResult {
  if (text.length == 0) {
    return new CollectionSearchResult(
      collection,
      CollectionStatus.Error,
      "Text is empty.",
      searchMethod,
      [],
    );
  }
  const result = hostSearchCollection(
    collection,
    namespaces,
    searchMethod,
    text,
    limit,
    returnText,
  );
  if (utils.resultIsInvalid(result)) {
    console.error("Error searching Text index.");
    return new CollectionSearchResult(
      collection,
      CollectionStatus.Error,
      "Error searching Text index.",
      searchMethod,
      [],
    );
  }
  return result;
}

export function searchByVector(
  collection: string,
  searchMethod: string,
  vector: f32[],
  limit: i32,
  returnText: bool = false,
  namespaces: string[] = [],
): CollectionSearchResult {
  if (vector.length == 0) {
    return new CollectionSearchResult(
      collection,
      CollectionStatus.Error,
      "Vector is empty.",
      searchMethod,
      [],
    );
  }
  const result = hostSearchCollectionByVector(
    collection,
    namespaces,
    searchMethod,
    vector,
    limit,
    returnText,
  );
  if (utils.resultIsInvalid(result)) {
    console.error("Error searching Text index by vector.");
    return new CollectionSearchResult(
      collection,
      CollectionStatus.Error,
      "Error searching Text index by vector.",
      searchMethod,
      [],
    );
  }
  return result;
}

// fetch embedders for collection & search method, run text through it and
// classify Text index for similar Texts, return the result keys
export function nnClassify(
  collection: string,
  searchMethod: string,
  text: string,
  namespace: string = "",
): CollectionClassificationResult {
  if (text.length == 0) {
    return new CollectionClassificationResult(
      collection,
      CollectionStatus.Error,
      "Text is empty.",
      searchMethod,
      [],
      [],
    );
  }
  const result = hostNnClassifyCollection(
    collection,
    namespace,
    searchMethod,
    text,
  );
  if (utils.resultIsInvalid(result)) {
    console.error("Error classifying Text index.");
    return new CollectionClassificationResult(
      collection,
      CollectionStatus.Error,
      "Error classifying Text index.",
      searchMethod,
      [],
      [],
    );
  }
  return result;
}

export function recomputeSearchMethod(
  collection: string,
  searchMethod: string,
  namespace: string = "",
): SearchMethodMutationResult {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new SearchMethodMutationResult(
      collection,
      CollectionStatus.Error,
      "Collection is empty.",
      "recompute",
      searchMethod,
    );
  }
  if (searchMethod.length == 0) {
    console.error("Search method is empty.");
    return new SearchMethodMutationResult(
      collection,
      CollectionStatus.Error,
      "Search method is empty.",
      "recompute",
      searchMethod,
    );
  }
  const result = hostRecomputeSearchMethod(collection, namespace, searchMethod);
  if (utils.resultIsInvalid(result)) {
    console.error("Error recomputing Text index.");
    return new SearchMethodMutationResult(
      collection,
      CollectionStatus.Error,
      "Error recomputing Text index.",
      "recompute",
      searchMethod,
    );
  }
  return result;
}

/**
 * @deprecated Use `collections.computeDistance` instead.
 */
export function computeSimilarity(
  collection: string,
  searchMethod: string,
  key1: string,
  key2: string,
  namespace: string = "",
): CollectionSearchResultObject {
  return computeDistance(collection, searchMethod, key1, key2, namespace);
}

export function computeDistance(
  collection: string,
  searchMethod: string,
  key1: string,
  key2: string,
  namespace: string = "",
): CollectionSearchResultObject {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new CollectionSearchResultObject("", "", "", [], 0.0, 0.0);
  }
  if (searchMethod.length == 0) {
    console.error("Search method is empty.");
    return new CollectionSearchResultObject("", "", "", [], 0.0, 0.0);
  }
  if (key1.length == 0) {
    console.error("Key1 is empty.");
    return new CollectionSearchResultObject("", "", "", [], 0.0, 0.0);
  }
  if (key2.length == 0) {
    console.error("Key2 is empty.");
    return new CollectionSearchResultObject("", "", "", [], 0.0, 0.0);
  }
  return hostComputeDistance(collection, namespace, searchMethod, key1, key2);
}

export function getText(
  collection: string,
  key: string,
  namespace: string = "",
): string {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return "";
  }
  if (key.length == 0) {
    console.error("Key is empty.");
    return "";
  }
  return hostGetTextFromCollection(collection, namespace, key);
}

export function getTexts(
  collection: string,
  namespace: string = "",
): Map<string, string> {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return new Map<string, string>();
  }
  return hostGetTextsFromCollection(collection, namespace);
}

export function getNamespaces(collection: string): string[] {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return [];
  }
  return hostGetNamespacesFromCollection(collection);
}

export function getVector(
  collection: string,
  searchMethod: string,
  key: string,
  namespace: string = "",
): f32[] {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return [];
  }
  if (searchMethod.length == 0) {
    console.error("Search method is empty.");
    return [];
  }
  if (key.length == 0) {
    console.error("Key is empty.");
    return [];
  }
  return hostGetVector(collection, namespace, searchMethod, key);
}

export function getLabels(
  collection: string,
  key: string,
  namespace: string = "",
): string[] {
  if (collection.length == 0) {
    console.error("Collection is empty.");
    return [];
  }
  if (key.length == 0) {
    console.error("Key is empty.");
    return [];
  }
  return hostGetLabels(collection, namespace, key);
}
