package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ModelSpec struct {
	ID        string `json:"id"`
	ModelType string `json:"modelType"`
	Endpoint  string `json:"endpoint"`
}

func GetModelEndpoint(url, id, jwt string) (string, error) {

	payload := []byte(fmt.Sprintf(`{"query":"query GetModelSpec {\n getModelSpec(id: \"%v\") {\n id\n modelType\n endpoint\n }\n}","variables":{}}`, id))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", jwt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Create an instance of the ModelSpec struct
	var spec ModelSpec

	// Unmarshal the JSON data into the ModelSpec struct
	err = json.Unmarshal(body, &spec)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	if spec.ID != id {
		return "", fmt.Errorf("error: ID does not match")
	}

	if spec.ModelType != "classifier" {
		return "", fmt.Errorf("error: model type does not match")
	}

	return spec.Endpoint, nil

}
