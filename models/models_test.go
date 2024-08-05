package models

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hmruntime/manifestdata"
	"hmruntime/secrets"

	"github.com/hypermodeAI/manifest"
	"github.com/stretchr/testify/assert"
)

const (
	testModelName = "test"
	testHostName  = "mock"
)

// TestMain runs in the main goroutine and can do whatever setup and teardown is necessary around a call to m.Run
func TestMain(m *testing.M) {
	secrets.Initialize(context.Background())
	manifestdata.SetManifest(&manifest.HypermodeManifest{
		Models: map[string]manifest.ModelInfo{
			testModelName: {
				Name:        testModelName,
				SourceModel: "",
				Provider:    "",
				Host:        testHostName,
			},
		},
		Hosts: map[string]manifest.HostInfo{
			testHostName: manifest.HTTPHostInfo{
				Name:     testHostName,
				Endpoint: "",
			},
		},
	})

	m.Run()
}

func TestGetModels(t *testing.T) {
	tests := []struct {
		desc      string
		valid     bool
		modelName string
	}{
		{
			desc:      "valid model",
			valid:     true,
			modelName: testModelName,
		},
		{
			desc:      "invalid model",
			valid:     false,
			modelName: "invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			model, err := GetModel(tc.modelName)
			if tc.valid {
				assert.NoError(t, err)
				assert.Equal(t, testModelName, model.Name)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPostExternalModelEndpoint(t *testing.T) {

	// Create an http handler that simply echoes the input, expecting a JSON object
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input map[string]any
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)

		output := input

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(output)
	})

	// Create a mock server with the handler to act as the external model endpoint
	tsrv := httptest.NewServer(handler)
	defer tsrv.Close()

	h := manifestdata.GetManifest().Hosts[testHostName].(manifest.HTTPHostInfo)
	h.Endpoint = tsrv.URL
	manifestdata.GetManifest().Hosts[testHostName] = h

	sentenceMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	testModel := manifestdata.GetManifest().Models[testModelName]
	resp, err := PostToModelEndpoint[map[string]string](context.Background(), testModel, sentenceMap)
	assert.NoError(t, err)

	// Expected response is the same as the input sentence map,
	// as the mock server just echoes the inputs
	assert.Equal(t, sentenceMap, resp)
}
