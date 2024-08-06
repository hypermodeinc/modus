package collections

import (
	"context"
	"fmt"
	"math"
	"sort"

	wasm "github.com/tetratelabs/wazero/api"

	collection_utils "hmruntime/collections/utils"
	"hmruntime/functions"
	"hmruntime/manifestdata"
	"hmruntime/utils"
	"hmruntime/wasmhost"
)

func UpsertToCollection(ctx context.Context, collectionName string, keys []string, texts []string, labels [][]string) (*collectionMutationResult, error) {

	// Get the collectionName data from the manifest
	collectionData := manifestdata.GetManifest().Collections[collectionName]

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		keys = make([]string, len(texts))
		for i := range keys {
			keys[i] = utils.GenerateUUIDv7()
		}
	}

	if len(labels) != 0 && len(labels) != len(texts) {
		return nil, fmt.Errorf("mismatch in number of labels and texts: %d != %d", len(labels), len(texts))
	}

	err = collection.InsertTexts(ctx, keys, texts, labels)
	if err != nil {
		return nil, err
	}

	// compute embeddings for each search method, and insert into vector index
	for searchMethodName, searchMethod := range collectionData.SearchMethods {
		vectorIndex, err := collection.GetVectorIndex(ctx, searchMethodName)
		if err != nil {
			return nil, err
		}

		embedder := searchMethod.Embedder

		info, err := functions.GetFunctionInfo(embedder)
		if err != nil {
			return nil, err
		}

		err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
		if err != nil {
			return nil, err
		}

		executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
		if err != nil {
			return nil, err
		}

		result := executionInfo.Result

		textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
		if err != nil {
			return nil, err
		}

		if len(textVecs) != len(texts) {
			return nil, fmt.Errorf("mismatch in number of embeddings generated by embedder %s", embedder)
		}

		ids := make([]int64, len(keys))
		for i := range textVecs {
			key := keys[i]

			id, err := collection.GetExternalId(ctx, key)
			if err != nil {
				return nil, err
			}
			ids[i] = id
		}

		err = vectorIndex.InsertVectors(ctx, ids, textVecs)
		if err != nil {
			return nil, err
		}
	}

	return &collectionMutationResult{
		Collection: collectionName,
		Operation:  "upsert",
		Status:     "success",
		Keys:       keys,
		Error:      "",
	}, nil
}

func DeleteFromCollection(ctx context.Context, collectionName string, key string) (*collectionMutationResult, error) {
	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	textId, err := collection.GetExternalId(ctx, key)
	if err != nil {
		return nil, err
	}
	for _, vectorIndex := range collection.GetVectorIndexMap() {
		err = vectorIndex.DeleteVector(ctx, textId, key)
		if err != nil {
			return nil, err
		}
	}
	err = collection.DeleteText(ctx, key)
	if err != nil {
		return nil, err
	}

	keys := []string{key}

	return &collectionMutationResult{
		Collection: collectionName,
		Operation:  "delete",
		Status:     "success",
		Keys:       keys,
		Error:      "",
	}, nil

}

func getEmbedder(ctx context.Context, collectionName string, searchMethod string) (string, error) {
	manifestColl, ok := manifestdata.GetManifest().Collections[collectionName]
	if !ok {
		return "", fmt.Errorf("collection %s not found in manifest", collectionName)
	}
	manifestSearchMethod, ok := manifestColl.SearchMethods[searchMethod]
	if !ok {
		return "", fmt.Errorf("search method %s not found in collection %s", searchMethod, collectionName)
	}
	embedder := manifestSearchMethod.Embedder
	if embedder == "" {
		return "", fmt.Errorf("embedder not found in search method %s of collection %s", searchMethod, collectionName)
	}
	return embedder, nil
}

func SearchCollection(ctx context.Context, collectionName string, searchMethod string,
	text string, limit int32, returnText bool) (*collectionSearchResult, error) {

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	embedder, err := getEmbedder(ctx, collectionName, searchMethod)
	if err != nil {
		return nil, err
	}

	info, err := functions.GetFunctionInfo(embedder)
	if err != nil {
		return nil, err
	}

	err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
	if err != nil {
		return nil, err
	}

	texts := []string{text}

	executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
	if err != nil {
		return nil, err
	}

	result := executionInfo.Result

	textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		return nil, err
	}

	if len(textVecs) == 0 {
		return nil, fmt.Errorf("no embeddings generated by embedder %s", embedder)
	}

	objects, err := vectorIndex.Search(ctx, textVecs[0], int(limit), nil)
	if err != nil {
		return nil, err
	}

	output := &collectionSearchResult{
		Collection:   collectionName,
		SearchMethod: searchMethod,
		Status:       "success",
		Objects:      make([]*collectionSearchResultObject, len(objects)),
	}

	for i, object := range objects {
		if returnText {
			text, err := collection.GetText(ctx, object.GetIndex())
			if err != nil {
				return nil, err
			}
			output.Objects[i] = &collectionSearchResultObject{
				Key:      object.GetIndex(),
				Text:     text,
				Distance: object.GetValue(),
				Score:    1 - object.GetValue(),
			}
		} else {
			output.Objects[i] = &collectionSearchResultObject{
				Key:      object.GetIndex(),
				Distance: object.GetValue(),
				Score:    1 - object.GetValue(),
			}
		}
	}

	return output, nil
}

func NnClassify(ctx context.Context, collectionName string, searchMethod string, text string) (*collectionClassificationResult, error) {

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	embedder, err := getEmbedder(ctx, collectionName, searchMethod)
	if err != nil {
		return nil, err
	}

	info, err := functions.GetFunctionInfo(embedder)
	if err != nil {
		return nil, err
	}

	err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
	if err != nil {
		return nil, err
	}

	texts := []string{text}

	executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
	if err != nil {
		return nil, err
	}

	result := executionInfo.Result

	textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		return nil, err
	}

	if len(textVecs) == 0 {
		return nil, fmt.Errorf("no embeddings generated by embedder %s", embedder)
	}

	lenTexts, err := collection.Len(ctx)
	if err != nil {
		return nil, err
	}

	nns, err := vectorIndex.Search(ctx, textVecs[0], int(math.Log10(float64(lenTexts)))*int(math.Log10(float64(lenTexts))), nil)
	if err != nil {
		return nil, err
	}

	// remove elements with score out of first standard deviation

	// calculate mean
	var sum float64
	for _, nn := range nns {
		sum += nn.GetValue()
	}
	mean := sum / float64(len(nns))

	// calculate standard deviation
	var variance float64
	for _, nn := range nns {
		variance += math.Pow(float64(nn.GetValue())-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(nns)))

	// remove elements with score out of first standard deviation and return the most frequent label
	labelCounts := make(map[string]int)

	res := &collectionClassificationResult{
		Collection:   collectionName,
		LabelsResult: make([]*collectionClassificationLabelObject, 0),
		SearchMethod: searchMethod,
		Status:       "success",
		Cluster:      make([]*collectionClassificationResultObject, 0),
	}

	totalLabels := 0

	for _, nn := range nns {
		if math.Abs(nn.GetValue()-mean) <= 2*stdDev {
			labels, err := collection.GetLabels(ctx, nn.GetIndex())
			if err != nil {
				return nil, err
			}
			for _, label := range labels {
				labelCounts[label]++
				totalLabels++
			}

			res.Cluster = append(res.Cluster, &collectionClassificationResultObject{
				Key:      nn.GetIndex(),
				Labels:   labels,
				Score:    1 - nn.GetValue(),
				Distance: nn.GetValue(),
			})
		}
	}

	// Create a slice of pairs
	labelsResult := make([]*collectionClassificationLabelObject, 0, len(labelCounts))
	for label, count := range labelCounts {
		labelsResult = append(labelsResult, &collectionClassificationLabelObject{
			Label:      label,
			Confidence: float64(count) / float64(totalLabels),
		})
	}

	// Sort the pairs by count in descending order
	sort.Slice(labelsResult, func(i, j int) bool {
		return labelsResult[i].Confidence > labelsResult[j].Confidence
	})

	res.LabelsResult = labelsResult

	return res, nil
}

func ComputeDistance(ctx context.Context, collectionName string, searchMethod string, id1 string, id2 string) (*collectionSearchResultObject, error) {

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	vec1, err := vectorIndex.GetVector(ctx, id1)
	if err != nil {
		return nil, err
	}

	vec2, err := vectorIndex.GetVector(ctx, id2)
	if err != nil {
		return nil, err
	}

	distance, err := collection_utils.CosineDistance(vec1, vec2)
	if err != nil {
		return nil, err
	}

	return &collectionSearchResultObject{
		Key:      "",
		Text:     "",
		Distance: distance,
		Score:    1 - distance,
	}, nil
}

func RecomputeSearchMethod(ctx context.Context, mod wasm.Module, collectionName string, searchMethod string) (*searchMethodMutationResult, error) {

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	embedder, err := getEmbedder(ctx, collectionName, searchMethod)
	if err != nil {
		return nil, err
	}

	info, err := functions.GetFunctionInfo(embedder)
	if err != nil {
		return nil, err
	}

	err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
	if err != nil {
		return nil, err
	}

	err = ProcessTextMapWithModule(ctx, mod, collection, embedder, vectorIndex)
	if err != nil {
		return nil, err
	}

	return &searchMethodMutationResult{
		Collection: collectionName,
		Operation:  "recompute",
		Status:     "success",
		Error:      "",
	}, nil

}

func GetTextFromCollection(ctx context.Context, collectionName string, key string) (string, error) {
	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return "", err
	}

	text, err := collection.GetText(ctx, key)
	if err != nil {
		return "", err
	}

	return text, nil
}

func GetTextsFromCollection(ctx context.Context, collectionName string) (map[string]string, error) {

	collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	textMap, err := collection.GetTextMap(ctx)
	if err != nil {
		return nil, err
	}

	return textMap, nil
}
