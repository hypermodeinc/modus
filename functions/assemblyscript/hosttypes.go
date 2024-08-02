/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import "hmruntime/plugins"

var mappedStructs = map[string]plugins.TypeInfo{
	"hmruntime/collections.collectionMutationResult": {
		Name: "CollectionMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	},
	"hmruntime/collections.searchMethodMutationResult": {
		Name: "SearchMethodMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	},
	"hmruntime/collections.collectionSearchResult": {
		Name: "CollectionSearchResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	},
	"hmruntime/collections.collectionSearchResultObject": {
		Name: "CollectionSearchResultObject",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	},
	"hmruntime/httpclient.HttpRequest": {
		Name: "Request",
		Path: "~lib/@hypermode/functions-as/assembly/http/Request",
	},
	"hmruntime/httpclient.HttpResponse": {
		Name: "Response",
		Path: "~lib/@hypermode/functions-as/assembly/http/Response",
	},
	"hmruntime/httpclient.HttpHeaders": {
		Name: "Headers",
		Path: "~lib/@hypermode/functions-as/assembly/http/Headers",
	},
	"hmruntime/httpclient.HttpHeader": {
		Name: "Header",
		Path: "~lib/@hypermode/functions-as/assembly/http/Header",
	},
	"hmruntime/models.ModelInfo": {
		Name: "ModelInfo",
		Path: "~lib/@hypermode/models-as/index/ModelInfo",
	},
	"hmruntime/sqlclient.hostQueryResponse": {
		Name: "HostQueryResponse",
		Path: "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	},
}
