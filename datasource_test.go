//go:build integration
// +build integration

package main

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

	"hmruntime/config"
	"hmruntime/httpserver"
	"hmruntime/services"

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
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(graphQLResponse.Errors) > 0 {
		return nil, errors.New(graphQLResponse.Errors[0].Message)
	}

	return graphQLResponse.Data, nil
}

// updateManifest updates the manifest to provided content and returns
// a function, calling the returned function restores the original manifest.
func updateManifest(t *testing.T, jsonManifest []byte) func() {
	manifestFilePath := path.Join(testPluginsPath, "hypermode.json")
	originalManifest, err := os.ReadFile(manifestFilePath)
	assert.Nil(t, err)

	assert.Nil(t, os.WriteFile(manifestFilePath, jsonManifest, os.ModePerm))
	return func() {
		assert.Nil(t, os.WriteFile(manifestFilePath, originalManifest, os.ModePerm))
	}
}

func TestMain(m *testing.M) {
	// setup config
	config.StoragePath = testPluginsPath
	config.RefreshInterval = refreshPluginInterval
	config.Port = httpListenPort

	// setup runtime services
	services.Start(context.Background())
	defer services.Stop(context.Background())

	// start HTTP server
	ctx, stop := context.WithCancel(context.Background())
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
	m.Run()

	// stop the HTTP server
	stop()
	<-done
}

func TestPostgresqlNoConnection(t *testing.T) {
	// wait here to make sure the plugin is loaded
	time.Sleep(waitRefreshPluginInterval)

	query := `query QueryPeople {queryPeople}`
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error setting up a new tx: failed to connect")
}

func TestPostgresqlNoHost(t *testing.T) {
	jsonManifest := []byte(`
{
  "hosts": {
    "postgres": {
      "type": "postgresql",
      "connString": "postgresql://postgres:password@localhost:5432/data?sslmode=disable"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name does not exist
	query := `query QueryPeople {queryPeople}`
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "postgresql host [neon] not found")
}

func TestPostgresqlNoPostgresqlHost(t *testing.T) {
	jsonManifest := []byte(`
{
  "hosts": {
    "neon": {
      "type": "http",
      "connString": "http://localhost:5432/data"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name has the wrong host type
	query := `query QueryPeople {queryPeople}`
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "host [neon] is not a postgresql host")
}

func TestPostgresqlWrongConnString(t *testing.T) {
	jsonManifest := []byte(`
{
  "hosts": {
    "neon": {
      "type": "postgresql",
      "connString": "http://localhost:5432/data"
    }
  }
}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when connection string is wrong
	query := `query QueryPeople {queryPeople}`
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to parse as keyword/value (invalid keyword/value)")
}

func TestPostgresqlNoConnString(t *testing.T) {
	jsonManifest := []byte(`
	{
	  "hosts": {
		"neon": {
		  "type": "postgresql"
		}
	  }
	}`)
	restoreFn := updateManifest(t, jsonManifest)
	defer restoreFn()

	// wait for new manifest to be loaded
	time.Sleep(waitRefreshPluginInterval)

	// when host name has no connection string
	query := `query QueryPeople {queryPeople}`
	response, err := runGraphqlQuery(graphQLRequest{Query: query})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "postgresql host [neon] has empty connString")
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
	io.Copy(os.Stdout, imagePullResp)

	ports := nat.PortSet{"5432": {}}
	pgEnv := []string{"POSTGRES_PASSWORD=password", "POSTGRES_DB=data"}
	cconf := &container.Config{Image: "docker.io/postgres:16-alpine", ExposedPorts: ports, Env: pgEnv}
	portBindings := nat.PortMap(map[nat.Port][]nat.PortBinding{"5432": {{HostPort: "5432"}}})
	hconf := &container.HostConfig{PortBindings: portBindings, AutoRemove: true}
	resp, err := ps.dcli.ContainerCreate(ctx, cconf, hconf, nil, nil, "postgres")
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
		time.Sleep(time.Second)

		pgLog, err := getLogs()
		if err == nil && 2 == strings.Count(pgLog, "database system is ready to accept connections") {
			break
		}

		counter++
		if counter == 3 {
			return fmt.Errorf("database did not come up")
		}
	}

	// create people table
	connStr := "postgresql://postgres:password@localhost:5432/data?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
create table people(
    id serial primary key,
    name varchar not null,
    age int not null,
    balance money not null,
    privkey bytea not null,
    created_at timestamp with time zone not null,
    male boolean,
    home point
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
query AddPerson {
    addPerson(
        name: "gulsaba"
        age: 26
        balance: 678943.34
        privkey: "this is a secret key"
        male: false
        lat: 26.9124
        lon: 75.7873
    )
}`
	_, err := runGraphqlQuery(graphQLRequest{Query: query})
	ps.Assert().Nil(err)
}

func (ps *postgresqlSuite) TestPostgresqlWrongTypeInsert() {
	// try inserting data with wrong type, column: balance
	query := `
query AddPerson {
    addPerson(
        name: "gulsaba"
        age: 26
        balance: "678943"
        privkey: "this is a secret key"
        male: false
        lat: 26.9124
        lon: 75.7873
    )
}`
	_, err := runGraphqlQuery(graphQLRequest{Query: query})
	ps.Assert().NotNil(err)
	ps.Assert().Contains(err.Error(), "input value is not a float")
}

func TestPostgresqlSuite(t *testing.T) {
	suite.Run(t, new(postgresqlSuite))
}
