/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

var hostTypes = map[string]string{
	"hmruntime/collections.CollectionMutationResult":             "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	"hmruntime/collections.SearchMethodMutationResult":           "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	"hmruntime/collections.CollectionSearchResult":               "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	"hmruntime/collections.CollectionSearchResultObject":         "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	"hmruntime/collections.CollectionClassificationResult":       "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResult",
	"hmruntime/collections.CollectionClassificationLabelObject":  "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationLabelObject",
	"hmruntime/collections.CollectionClassificationResultObject": "~lib/@hypermode/functions-as/assembly/collections/CollectionClassificationResultObject",
	"hmruntime/httpclient.HttpRequest":                           "~lib/@hypermode/functions-as/assembly/http/Request",
	"hmruntime/httpclient.HttpResponse":                          "~lib/@hypermode/functions-as/assembly/http/Response",
	"hmruntime/httpclient.HttpHeaders":                           "~lib/@hypermode/functions-as/assembly/http/Headers",
	"hmruntime/httpclient.HttpHeader":                            "~lib/@hypermode/functions-as/assembly/http/Header",
	"hmruntime/sqlclient.HostQueryResponse":                      "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	"hmruntime/models.ModelInfo":                                 "~lib/@hypermode/models-as/index/ModelInfo",
}
