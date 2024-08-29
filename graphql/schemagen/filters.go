/*
 * Copyright 2024 Hypermode, Inc.
 */

package schemagen

import "hypruntime/manifestdata"

func getFnFilter() func(*FunctionSignature) bool {
	embedders := make(map[string]bool)
	for _, collection := range manifestdata.GetManifest().Collections {
		for _, searchMethod := range collection.SearchMethods {
			embedders[searchMethod.Embedder] = true
		}
	}

	return func(f *FunctionSignature) bool {
		return !embedders[f.Name]
	}
}
