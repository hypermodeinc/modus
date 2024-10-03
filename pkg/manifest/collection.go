package manifest

type CollectionInfo struct {
	SearchMethods map[string]SearchMethodInfo `json:"searchMethods"`
}

type SearchMethodInfo struct {
	Embedder string    `json:"embedder"`
	Index    IndexInfo `json:"index"`
}

type IndexInfo struct {
	Type    string      `json:"type"`
	Options OptionsInfo `json:"options"`
}

type OptionsInfo struct {
	EfConstruction int `json:"efConstruction"`
	MaxLevels      int `json:"maxLevels"`
}
