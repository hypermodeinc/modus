package legacymodels

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hypermodeinc/modus/pkg/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"github.com/stretchr/testify/assert"
)

const (
	testModelName = "test"
	testHostName  = "mock"
)

type requestBody struct {
	Instances []string `json:"instances"`
}
type mockModelResponse struct {
	Predictions []string `json:"predictions"`
}

// TestMain runs in the main goroutine and can do whatever setup and teardown is necessary around a call to m.Run
func TestMain(m *testing.M) {
	secrets.Initialize(context.Background())
	manifestdata.SetManifest(&manifest.Manifest{
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

	os.Exit(m.Run())
}

func TestPostExternalModelEndpoint(t *testing.T) {

	// Create an http handler that simply echoes the input strings
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b requestBody
		err := json.NewDecoder(r.Body).Decode(&b)
		assert.NoError(t, err)

		resp := mockModelResponse{
			Predictions: b.Instances,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
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
	resp, err := postToModelEndpoint[string](context.Background(), &testModel, sentenceMap)
	assert.NoError(t, err)

	// Expected response is the same as the input sentence map,
	// as the mock server just echoes the inputs
	assert.Equal(t, sentenceMap, resp)
}
