package models

import (
	"context"
	"encoding/json"
	"fmt"
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
	testAuthHeader   = ""
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
				AuthHeader: testAuthHeader,
			},
		},
	}

	m.Run()
}

func TestGetModels(t *testing.T) {
	model, err := GetModel(testModelName, testModelTask)
	assert.NoError(t, err)
	assert.Equal(t, testModelName, model.Name)
}

func TestPostExternalModelEndpoint(t *testing.T) {
	l, err := net.Listen("tcp", testHostEndpoint)
	assert.NoError(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b requestBody
		err := json.NewDecoder(r.Body).Decode(&b)
		assert.NoError(t, err)
		fmt.Println(b.Instances)

		resp := mockModelResponse{
			Predictions: b.Instances,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	tsrv := httptest.NewUnstartedServer(handler)
	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
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

	assert.Equal(t, sentenceMap, resp)
}
