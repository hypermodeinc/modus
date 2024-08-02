/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

var hostTypes = map[string]string{
	"hmruntime/collections.collectionMutationResult":             "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	"hmruntime/collections.searchMethodMutationResult":           "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	"hmruntime/collections.collectionSearchResult":               "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	"hmruntime/collections.collectionSearchResultObject":         "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	"hmruntime/collections.collectionClassificationResult":       "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResult",
	"hmruntime/collections.collectionClassificationLabelObject":  "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationLabelObject",
	"hmruntime/collections.collectionClassificationResultObject": "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResultObject",
	"hmruntime/httpclient.HttpRequest":                           "~lib/@hypermode/functions-as/assembly/http/Request",
	"hmruntime/httpclient.HttpResponse":                          "~lib/@hypermode/functions-as/assembly/http/Response",
	"hmruntime/httpclient.HttpHeaders":                           "~lib/@hypermode/functions-as/assembly/http/Headers",
	"hmruntime/httpclient.HttpHeader":                            "~lib/@hypermode/functions-as/assembly/http/Header",
	"hmruntime/sqlclient.hostQueryResponse":                      "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	"hmruntime/models.ModelInfo":                                 "~lib/@hypermode/models-as/index/ModelInfo",
}
