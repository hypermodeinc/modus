/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

var hostTypes = map[string]string{
	"hypruntime/collections.CollectionMutationResult":             "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	"hypruntime/collections.SearchMethodMutationResult":           "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	"hypruntime/collections.CollectionSearchResult":               "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	"hypruntime/collections.CollectionSearchResultObject":         "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	"hypruntime/collections.CollectionClassificationResult":       "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResult",
	"hypruntime/collections.CollectionClassificationLabelObject":  "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationLabelObject",
	"hypruntime/collections.CollectionClassificationResultObject": "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResultObject",
	"hypruntime/httpclient.HttpRequest":                           "~lib/@hypermode/functions-as/assembly/http/Request",
	"hypruntime/httpclient.HttpResponse":                          "~lib/@hypermode/functions-as/assembly/http/Response",
	"hypruntime/httpclient.HttpHeaders":                           "~lib/@hypermode/functions-as/assembly/http/Headers",
	"hypruntime/httpclient.HttpHeader":                            "~lib/@hypermode/functions-as/assembly/http/Header",
	"hypruntime/sqlclient.HostQueryResponse":                      "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	"hypruntime/models.ModelInfo":                                 "~lib/@hypermode/models-as/index/ModelInfo",
	"hypruntime/dgraphclient.Request":                             "~lib/@hypermode/functions-as/assembly/dgraph/Request",
	"hypruntime/dgraphclient.Response":                            "~lib/@hypermode/functions-as/assembly/dgraph/Response",
	"hypruntime/dgraphclient.Query":                               "~lib/@hypermode/functions-as/assembly/dgraph/Query",
	"hypruntime/dgraphclient.Mutation":                            "~lib/@hypermode/functions-as/assembly/dgraph/Mutation",
}
