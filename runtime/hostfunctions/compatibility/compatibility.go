/*
 * Copyright 2024 Hypermode, Inc.
 */

package compatibility

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/plugins/metadata"
)

// This is a compatibility shim for Hypermode host functions that are not in the metadata.
// It is only needed to allow previous AssemblyScript SDK versions to work with the latest runtime.
// This should be removed when all plugins are updated to the latest SDK version that uses the V2 metadata format.

func GetImportMetadataShim(fnName string) (fn *metadata.Function, err error) {
	switch fnName {
	case "hypermode.log":
		fn = metadata.NewFunction(fnName).
			WithParameter("level", "~lib/string/String").
			WithParameter("message", "~lib/string/String")

	/* Model function imports */

	case "hypermode.lookupModel":
		fn = metadata.NewFunction(fnName).
			WithParameter("modelName", "~lib/string/String").
			WithResult("~lib/@hypermode/models-as/index/ModelInfo")

	case "hypermode.invokeModel":
		fn = metadata.NewFunction(fnName).
			WithParameter("modelName", "~lib/string/String").
			WithParameter("input", "~lib/string/String").
			WithResult("~lib/string/String | null")

	/* HTTP, GraphQL, and SQL function imports */

	case "hypermode.httpFetch":
		fn = metadata.NewFunction(fnName).
			WithParameter("url", "~lib/@hypermode/functions-as/assembly/http/Request").
			WithResult("~lib/@hypermode/functions-as/assembly/http/Response")

	case "hypermode.executeGQL":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithParameter("statement", "~lib/string/String").
			WithParameter("variables", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.databaseQuery":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithParameter("dbType", "~lib/string/String").
			WithParameter("statement", "~lib/string/String").
			WithParameter("paramsJson", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/database/HostQueryResponse")

	/* Dgraph function imports */

	case "hypermode.dgraphAlterSchema":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithParameter("schema", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.dgraphDropAll":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.dgraphDropAttr":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithParameter("attr", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.executeDQL":
		fn = metadata.NewFunction(fnName).
			WithParameter("hostName", "~lib/string/String").
			WithParameter("request", "~lib/@hypermode/functions-as/assembly/dgraph/Request").
			WithResult("~lib/@hypermode/functions-as/assembly/dgraph/Response")

	/* Collection function imports */

	case "hypermode.computeDistance_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("key1", "~lib/string/String").
			WithParameter("key2", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject")

	case "hypermode.deleteFromCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult")

	case "hypermode.getNamespacesFromCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithResult("~lib/array/Array<~lib/string/String>")

	case "hypermode.getTextFromCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.getTextsFromCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithResult("~lib/map/Map<~lib/string/String,~lib/string/String>")

	case "hypermode.getVector":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/array/Array<f32>")

	case "hypermode.getLabels":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/array/Array<~lib/string/String>")

	case "hypermode.nnClassifyCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("text", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResult")

	case "hypermode.recomputeSearchMethod_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult")

	case "hypermode.searchCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespaces", "~lib/array/Array<~lib/string/String>").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("text", "~lib/string/String").
			WithParameter("limit", "i32").
			WithParameter("returnText", "bool").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult")

	case "hypermode.searchCollectionByVector":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespaces", "~lib/array/Array<~lib/string/String>").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("vector", "~lib/array/Array<f32>").
			WithParameter("limit", "i32").
			WithParameter("returnText", "bool").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult")

	case "hypermode.upsertToCollection_v2":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("namespace", "~lib/string/String").
			WithParameter("keys", "~lib/array/Array<~lib/string/String>").
			WithParameter("texts", "~lib/array/Array<~lib/string/String>").
			WithParameter("labels", "~lib/array/Array<~lib/array/Array<~lib/string/String>>").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult")

	/* Older collection function imports */

	case "hypermode.computeDistance":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("key1", "~lib/string/String").
			WithParameter("key2", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject")

	case "hypermode.computeSimilarity":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("key1", "~lib/string/String").
			WithParameter("key2", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject")

	case "hypermode.deleteFromCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult")

	case "hypermode.getTextFromCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("key", "~lib/string/String").
			WithResult("~lib/string/String")

	case "hypermode.getTextsFromCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithResult("~lib/map/Map<~lib/string/String,~lib/string/String>")

	case "hypermode.nnClassifyCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("text", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResult")

	case "hypermode.recomputeSearchMethod":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult")

	case "hypermode.searchCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("searchMethod", "~lib/string/String").
			WithParameter("text", "~lib/string/String").
			WithParameter("limit", "i32").
			WithParameter("returnText", "bool").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult")

	case "hypermode.upsertToCollection":
		fn = metadata.NewFunction(fnName).
			WithParameter("collection", "~lib/string/String").
			WithParameter("keys", "~lib/array/Array<~lib/string/String>").
			WithParameter("texts", "~lib/array/Array<~lib/string/String>").
			WithResult("~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult")

	/* Legacy model function imports */

	case "hypermode.invokeClassifier":
		fn = metadata.NewFunction(fnName).
			WithParameter("modelName", "~lib/string/String").
			WithParameter("sentenceMap", "~lib/map/Map<~lib/string/String,~lib/string/String>").
			WithResult("~lib/map/Map<~lib/string/String,~lib/map/Map<~lib/string/String,f32>>")

	case "hypermode.computeEmbedding":
		fn = metadata.NewFunction(fnName).
			WithParameter("modelName", "~lib/string/String").
			WithParameter("sentenceMap", "~lib/map/Map<~lib/string/String,~lib/string/String>").
			WithResult("~lib/map/Map<~lib/string/String,~lib/array/Array<f64>>")

	case "hypermode.invokeTextGenerator":
		fn = metadata.NewFunction(fnName).
			WithParameter("modelName", "~lib/string/String").
			WithParameter("instruction", "~lib/string/String").
			WithParameter(("format"), "~lib/string/String").
			WithResult("~lib/string/String")

	default:
		return nil, fmt.Errorf("unknown host function: %s", fnName)
	}
	return fn, nil
}
