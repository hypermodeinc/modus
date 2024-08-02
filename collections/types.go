package collections

type collectionMutationResult struct {
	Collection string
	Operation  string
	Status     string
	Keys       []string
	Error      string
}

type searchMethodMutationResult struct {
	Collection   string
	SearchMethod string
	Operation    string
	Status       string
	Error        string
}

type collectionSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []collectionSearchResultObject
	Error        string
}

type collectionSearchResultObject struct {
	Key      string
	Text     string
	Distance float64
	Score    float64
}

type collectionClassificationResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Label        string
	LabelsResult []string
	Cluster      []collectionClassificationResultObject
	Error        string
}

type collectionClassificationResultObject struct {
	Key      string
	Label    string
	Distance float64
	Score    float64
}
