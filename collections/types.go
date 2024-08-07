package collections

type CollectionMutationResult struct {
	Collection string
	Operation  string
	Status     string
	Keys       []string
	Error      string
}

type SearchMethodMutationResult struct {
	Collection   string
	SearchMethod string
	Operation    string
	Status       string
	Error        string
}

type CollectionSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []*CollectionSearchResultObject
	Error        string
}

type CollectionSearchResultObject struct {
	Key      string
	Text     string
	Distance float64
	Score    float64
}

type CollectionClassificationResult struct {
	Collection   string
	SearchMethod string
	Status       string
	LabelsResult []*CollectionClassificationLabelObject
	Cluster      []*CollectionClassificationResultObject
	Error        string
}

type CollectionClassificationLabelObject struct {
	Label      string
	Confidence float64
}

type CollectionClassificationResultObject struct {
	Key      string
	Labels   []string
	Distance float64
	Score    float64
}
