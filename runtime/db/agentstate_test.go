package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/db"
)

func setupModusDBTest(t *testing.T) (context.Context, func()) {
	os.Setenv("MODUS_USE_MODUSDB", "true")
	ctx := context.Background()
	tmpDir := t.TempDir()
	cfg := app.NewAppConfig().WithAppPath(tmpDir)
	app.SetConfig(cfg)
	db.InitModusDb(ctx, "")
	cleanup := func() {
		db.CloseModusDb(ctx)
	}
	return ctx, cleanup
}

func TestAgentStateModusGraph(t *testing.T) {
	ctx, cleanup := setupModusDBTest(t)
	defer cleanup()

	// 1. Write
	agent := db.AgentState{
		Id:        "agent-123",
		Name:      "TestAgent",
		Status:    "active",
		Data:      "{\"foo\":\"bar\"}",
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
		DType:     []string{"AgentState"},
	}
	err := db.WriteAgentState(ctx, agent)
	require.NoError(t, err, "WriteAgentState should succeed")

	// 2. Get
	got, err := db.GetAgentState(ctx, agent.Id)
	require.NoError(t, err, "GetAgentState should succeed")
	require.NotNil(t, got, "Returned AgentState should not be nil")
	assert.Equal(t, agent.Id, got.Id, "AgentState ID should match")
	assert.Equal(t, agent.Name, got.Name, "AgentState Name should match")
	assert.Equal(t, agent.Status, got.Status, "AgentState Status should match")
	assert.Equal(t, agent.Data, got.Data, "AgentState Data should match")
	assert.Equal(t, agent.UpdatedAt, got.UpdatedAt, "AgentState UpdatedAt should match")

	// 3. Update
	err = db.UpdateAgentStatus(ctx, agent.Id, "terminated")
	require.NoError(t, err, "UpdateAgentStatus should succeed")
	got, err = db.GetAgentState(ctx, agent.Id)
	require.NoError(t, err, "GetAgentState should succeed")
	assert.Equal(t, "terminated", got.Status, "AgentState Status should match")
	assert.Equal(t, agent.UpdatedAt, got.UpdatedAt, "AgentState UpdatedAt should match")

	// 4. Query
	agents, err := db.QueryActiveAgents(ctx)
	require.NoError(t, err, "QueryActiveAgents should succeed")
	found := false
	for _, a := range agents {
		if a.Id == agent.Id {
			found = true
			break
		}
	}
	assert.False(t, found, "Agent should not be active after status update to terminated")

	// Set status back to active and query again
	err = db.UpdateAgentStatus(ctx, agent.Id, "active")
	require.NoError(t, err, "UpdateAgentStatus should succeed")
	agents, err = db.QueryActiveAgents(ctx)
	require.NoError(t, err, "QueryActiveAgents should succeed")
	found = false
	for _, a := range agents {
		if a.Id == agent.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "Agent should be active after status update to active")
}
