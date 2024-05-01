package models

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"hmruntime/manifest"
)

const (
	testModelName    = "test"
	testModelTask    = "echo"
	testHostName     = "mock"
	testHostEndpoint = "127.0.0.1:8080"
)

type requestBody struct {
	Instances []string `json:"instances"`
}
type mockModelResponse struct {
	Predictions []string `json:"predictions"`
}

// TestMain runs in the main goroutine and can do whatever setup and teardown is necessary around a call to m.Run
func TestMain(m *testing.M) {
	manifest.HypermodeData = manifest.HypermodeManifest{
		Models: []manifest.Model{
			{
				Name:        testModelName,
				Task:        testModelTask,
				SourceModel: "",
				Provider:    "",
				Host:        testHostName,
			},
		},
		Hosts: []manifest.Host{
			{
				Name:       testHostName,
				Endpoint:   "http://" + testHostEndpoint,
				AuthHeader: "",
			},
		},
	}

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
			model, err := GetModel(tc.modelName, testModelTask)
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
	// Create a mock server to act as the external model endpoint
	l, err := net.Listen("tcp", testHostEndpoint)
	assert.NoError(t, err)
	defer l.Close()
	// Create a handler that simply echoes the input strings
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b requestBody
		err := json.NewDecoder(r.Body).Decode(&b)
		assert.NoError(t, err)

		resp := mockModelResponse{
			Predictions: b.Instances,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	tsrv := httptest.NewUnstartedServer(handler)
	// NewUnstartedServer creates a listener. Close that listener and replace with l
	tsrv.Listener.Close()
	tsrv.Listener = l
	tsrv.Start()
	defer tsrv.Close()

	sentenceMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	testModel := manifest.HypermodeData.Models[0]
	resp, err := PostToModelEndpoint[string](context.Background(), sentenceMap, testModel)
	assert.NoError(t, err)

	// Expected response is the same as the input sentence map,
	// as the mock server just echoes the inputs
	assert.Equal(t, sentenceMap, resp)
}
