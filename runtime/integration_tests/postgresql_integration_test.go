//go:build integration

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/httpserver"
	"github.com/hypermodeinc/modus/runtime/services"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	testPluginsPath           = "./testdata"
	refreshPluginInterval     = 100 * time.Millisecond
	waitRefreshPluginInterval = refreshPluginInterval * 2
	httpListenPort            = 8765

	requestTimeout = time.Minute

	healthURL  = "http://localhost:8765/health"
	graphqlURL = "http://localhost:8765/graphql"

	// Note: most DB query failures will fail with this error message in the graphql response.
	// The log messages will be more detailed, but those are intentionally not exposed to the client.
	expectedError = "error calling function"
)

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func runGraphqlQuery(greq graphQLRequest) ([]byte, error) {
	payload, err := json.Marshal(greq)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	resp, err := http.Post(graphqlURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var graphQLResponse struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}
	if err := json.Unmarshal(respBody, &graphQLResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(graphQLResponse.Errors) > 0 {
		return nil, errors.New(graphQLResponse.Errors[0].Message)
	}

	return graphQLResponse.Data, nil
}

// updateManifest updates the manifest to provided content and returns
// a function, calling the returned function restores the original manifest.
func updateManifest(t *testing.T, jsonManifest []byte) func() {
	manifestFilePath := path.Join(testPluginsPath, "modus.json")
	originalManifest, err := os.ReadFile(manifestFilePath)
	assert.Nil(t, err)

	assert.Nil(t, os.WriteFile(manifestFilePath, jsonManifest, os.ModePerm))
	return func() {
		assert.Nil(t, os.WriteFile(manifestFilePath, originalManifest, os.ModePerm))
	}
}

func TestMain(m *testing.M) {
	// setup config
	cfg := app.NewAppConfig().
		WithAppPath(testPluginsPath).
		WithRefreshInterval(refreshPluginInterval).
		WithPort(httpListenPort)
	app.SetConfig(cfg)

	// Create the main background context
	ctx := context.Background()

	// setup runtime services
	services.Start(ctx)
	defer services.Stop(ctx)

	// start HTTP server
	ctx, stop := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		httpserver.Start(ctx, true)
		close(done)
	}()

	// wait for HTTP server to be up
	counter := 0
	for {
		time.Sleep(time.Millisecond * 100)

		resp, err := http.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		counter++
		if counter == 3 {
			panic("error starting HTTP server")
		}
	}

	// run tests
	exitCode := m.Run()

	// stop the HTTP server
	stop()
	<-done

	os.Exit(exitCode)
}

func TestPostgresqlNoConnection(t *testing.T) {
	// wait here to make sure the plugin is loaded
	time.Sleep(waitRefreshPluginInterval)

	query := "{ allPeople { id name age } }"
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedError)
}

func TestPostgresqlNoHost(t *testing.T) {
	jsonManifest := []byte(`
{
  "endpoints": {
    "default": {
      "type": "graphql",
      "path": "/graphql",
      "auth": "bearer-token"
    }
  },
  "connections": {
    "postgres": {
      "type": "postgresql",
      "connString": "postgresql://postgres:password@localhost:5499/data?sslmode=disable"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name does not exist
	query := "{ allPeople { id name age } }"
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedError)
}

func TestPostgresqlNoPostgresqlHost(t *testing.T) {
	jsonManifest := []byte(`
{
  "endpoints": {
    "default": {
      "type": "graphql",
      "path": "/graphql",
      "auth": "bearer-token"
    }
  },
  "connections": {
    "my-database": {
      "type": "http",
      "connString": "http://localhost:5499/data"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name has the wrong host type
	query := "{ allPeople { id name age } }"
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedError)
}

func TestPostgresqlWrongConnString(t *testing.T) {
	jsonManifest := []byte(`
{
  "endpoints": {
    "default": {
      "type": "graphql",
      "path": "/graphql",
      "auth": "bearer-token"
    }
  },
  "connections": {
    "my-database": {
      "type": "postgresql",
      "connString": "http://localhost:5499/data"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when connection string is wrong
	query := "{ allPeople { id name age } }"
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedError)
}

func TestPostgresqlNoConnString(t *testing.T) {
	jsonManifest := []byte(`
{
  "endpoints": {
    "default": {
      "type": "graphql",
      "path": "/graphql",
      "auth": "bearer-token"
    }
  },
  "connections": {
    "my-database": {
      "type": "postgresql"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name has no connection string
	query := "{ allPeople { id name age } }"
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedError)
}

type postgresqlSuite struct {
	suite.Suite
	dcli   *docker.Client
	contID string
}

func (ps *postgresqlSuite) setupPostgresContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	imagePullResp, err := ps.dcli.ImagePull(ctx, "docker.io/postgres:16-alpine", image.PullOptions{})
	if err != nil {
		return fmt.Errorf("error pulling docker image: %w", err)
	}
	defer imagePullResp.Close()
	if _, err := io.Copy(os.Stdout, imagePullResp); err != nil {
		return fmt.Errorf("error copying image pull response: %w", err)
	}

	pgEnv := []string{"POSTGRES_PASSWORD=password", "POSTGRES_DB=data"}
	cconf := &container.Config{Image: "docker.io/postgres:16-alpine", ExposedPorts: nat.PortSet{"5432": {}}, Env: pgEnv}
	portBindings := nat.PortMap(map[nat.Port][]nat.PortBinding{"5432": {{HostPort: "5499"}}})
	hconf := &container.HostConfig{PortBindings: portBindings, AutoRemove: true}
	resp, err := ps.dcli.ContainerCreate(ctx, cconf, hconf, nil, nil, "postgres_modus_integration_test")
	if err != nil {
		return err
	}
	ps.contID = resp.ID

	if err := ps.dcli.ContainerStart(ctx, ps.contID, container.StartOptions{}); err != nil {
		return err
	}

	getLogs := func() (string, error) {
		opts := container.LogsOptions{ShowStdout: true, ShowStderr: true, Details: true}
		ro, err := ps.dcli.ContainerLogs(ctx, ps.contID, opts)
		if err != nil {
			return "", fmt.Errorf("error collecting logs for %v: %w", ps.contID, err)
		}
		defer ro.Close()

		data, err := io.ReadAll(ro)
		if err != nil {
			return "", fmt.Errorf("error in reading logs from stream: %w", err)
		}
		return string(data), nil
	}

	// wait for postgres to be up
	counter := 0
	for {
		pgLog, err := getLogs()
		if err == nil && strings.Count(pgLog, "database system is ready to accept connections") == 2 {
			break
		}

		counter++
		if counter == 5 {
			return fmt.Errorf("database did not come up")
		}

		time.Sleep(time.Second * 2)
	}

	// create people table
	connStr := "postgresql://postgres:password@localhost:5499/data?sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	sql := `
CREATE TABLE people (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  age INT NOT NULL,
  home POINT NULL
)`
	if _, err := conn.Exec(ctx, sql); err != nil {
		return err
	}

	return nil
}

func (ps *postgresqlSuite) SetupSuite() {
	dcli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	ps.dcli = dcli

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	if _, err := dcli.Ping(ctx); err != nil {
		panic(err)
	}

	if err := ps.setupPostgresContainer(); err != nil {
		panic(err)
	}
}

func (ps *postgresqlSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	ro := container.RemoveOptions{RemoveVolumes: true, Force: true}
	if err := ps.dcli.ContainerRemove(ctx, ps.contID, ro); err != nil {
		panic(err)
	}
}

func (ps *postgresqlSuite) TestPostgresqlBasicOps() {
	query := `
mutation AddPerson {
    addPerson(name: "test", age: 21) {
        id
        name
        age
    }
}`
	_, err := runGraphqlQuery(graphQLRequest{Query: query})
	ps.Assert().Nil(err)
}

func (ps *postgresqlSuite) TestPostgresqlWrongTypeInsert() {
	// try inserting data with wrong type, column: age
	query := `
mutation AddPerson {
    addPerson(name: "test", age: "abc") {
        id
        name
        age
    }
}`
	_, err := runGraphqlQuery(graphQLRequest{Query: query})

	// Note, the expected error is related to query input validation, not function execution.
	ps.Assert().NotNil(err)
	ps.Assert().EqualError(err, `Int cannot represent non-integer value: "abc"`)
}

func TestPostgresqlSuite(t *testing.T) {
	suite.Run(t, new(postgresqlSuite))
}
