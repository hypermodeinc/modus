package collections

import "hmruntime/plugins"

type collectionMutationResult struct {
	Collection string
	Operation  string
	Status     string
	Keys       []string
	Error      string
}

func (r *collectionMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	}
}

type searchMethodMutationResult struct {
	Collection   string
	SearchMethod string
	Operation    string
	Status       string
	Error        string
}

func (r *searchMethodMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "SearchMethodMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	}
}

type collectionSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []collectionSearchResultObject
	Error        string
}

func (r *collectionSearchResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	}
}

type collectionSearchResultObject struct {
	Key      string
	Text     string
	Distance float64
	Score    float64
}

func (r *collectionSearchResultObject) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResultObject",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	}
}
