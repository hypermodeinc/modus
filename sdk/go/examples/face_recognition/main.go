/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/modusdb"
)

func AlterSchema() string {
	schema := `
	base64: string @index(exact) .
	base64_v: float32vector @index(hnsw(exponent: "5", metric: "euclidean")) .
  
	type Image {
		base64
		base64_v
	}
	`
	err := modusdb.AlterSchema(schema)
	if err != nil {
		return err.Error()
	}

	return "Success"
}

func FetchImageEmbedding(base64 string) (*ImageEmbedding, error) {
	req := http.NewRequest("http://localhost:5000/encode", &http.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"image": base64,
		},
	})

	res, err := http.Fetch(req)
	if err != nil {
		return nil, err
	}

	var embedding *ImageEmbedding

	res.JSON(embedding)

	return embedding, nil
}

func AddImageAsJSON(base64 string) (*map[string]uint64, error) {

	embedding, err := FetchImageEmbedding(base64)
	if err != nil {
		return nil, err
	}
	image := Image{
		Uid:     "_:image1",
		Base64:  base64,
		Base64V: embedding.Encoding,
		DType:   []string{"Image"},
	}

	data, err := json.Marshal(image)
	if err != nil {
		return nil, err
	}

	response, err := modusdb.Mutate(&modusdb.MutationRequest{
		Mutations: []*modusdb.Mutation{
			{
				SetJson: string(data),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func FindSimilarImage(base64 string) (*Image, error) {
	embedding, err := FetchImageEmbedding(base64)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
	{
		images(func: similar_to(base64_v, 1, "%v")) {
				uid
				base64
				base64_v
				dgraph.type
		}
	}
	`, embedding.Encoding)

	response, err := modusdb.Query(&query)
	if err != nil {
		return nil, err
	}

	var images ImagesData
	if err := json.Unmarshal([]byte(response.Json), &images); err != nil {
		return nil, err
	}

	if len(images.Images) == 0 {
		return nil, nil
	}

	image := images.Images[0]

	return image, nil
}
