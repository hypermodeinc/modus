/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/collections"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/experimental"
)

const (
	modelName  = "minilm"
	myProducts = "myProducts"
)

var searchMethods = []string{"searchMethod1", "searchMethod2"}

func Embed(text []string) ([][]float32, error) {
	model, err := models.GetModel[experimental.EmbeddingsModel](modelName)
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(text...)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	return output.Predictions, nil

}

func AddProduct(description string) ([]string, error) {
	res, err := collections.Upsert(myProducts, nil, description, nil)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func AddProducts(descriptions []string) ([]string, error) {
	res, err := collections.UpsertBatch(myProducts, nil, descriptions, nil)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func DeleteProduct(key string) (string, error) {
	res, err := collections.Remove(myProducts, key)
	if err != nil {
		return "", err
	}
	return res.Status, nil
}

func GetProduct(key string) (string, error) {
	return collections.GetText(myProducts, key)
}

func GetProducts() (map[string]string, error) {
	return collections.GetTexts(myProducts)
}

func GetProductVector(key string) ([]float32, error) {
	return collections.GetVector(myProducts, searchMethods[0], key)
}

func GetNamespaces() ([]string, error) {
	return collections.GetNamespaces(myProducts)
}

func GetLabels(key string) ([]string, error) {
	return collections.GetLabels(myProducts, key)
}

func RecomputeSearchMethods() (string, error) {
	for _, method := range searchMethods {
		_, err := collections.RecomputeSearchMethod(myProducts, method)
		if err != nil {
			return "", err
		}
	}
	return "success", nil
}

func ComputeDistanceBetweenProducts(key1, key2 string) (float64, error) {
	res, err := collections.ComputeDistance(myProducts, searchMethods[0], key1, key2)
	if err != nil {
		return 0, err
	}
	return res.Distance, nil
}

func SearchProducts(productDescription string, maxItems int) (*collections.CollectionSearchResult, error) {
	return collections.Search(myProducts, searchMethods[0], productDescription, collections.WithLimit(maxItems), collections.WithReturnText(true))
}

func SearchProductsWithKey(productKey string, maxItems int) (*collections.CollectionSearchResult, error) {
	vec, err := GetProductVector(productKey)
	if err != nil {
		return nil, err
	}
	return collections.SearchByVector(myProducts, searchMethods[0], vec, collections.WithLimit(maxItems), collections.WithReturnText(true))
}
