package db_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModusDb(t *testing.T) {
	// Setup test environment
	os.Setenv("MODUS_USE_MODUSDB", "true")
	ctx := context.Background()
	testAppPath := t.TempDir()
	cfg := app.NewAppConfig().WithAppPath(testAppPath)
	app.SetConfig(cfg)

	// Test initialization with default path
	t.Run("InitModusDb with default path", func(t *testing.T) {
		db.InitModusDb(ctx, "")
		client, err := db.GetClient()
		require.NoError(t, err)
		defer client.Close()
		assert.NotNil(t, client, "Client should be initialized")

		// Verify .modusdb directory was created
		modusDbPath := filepath.Join(testAppPath, ".modusdb")
		_, err = os.Stat(modusDbPath)
		assert.NoError(t, err, "ModusDB directory should exist")

		db.CloseModusDb(ctx)
	})

	// Test initialization with custom file URI
	t.Run("InitModusDb with custom file URI", func(t *testing.T) {
		customPath := filepath.Join(testAppPath, "custom")
		err := os.MkdirAll(customPath, 0755)
		require.NoError(t, err)

		db.InitModusDb(ctx, "file://"+customPath)
		client, err := db.GetClient()
		require.NoError(t, err)
		defer client.Close()

		assert.NotNil(t, client, "Client should be initialized with custom path")

		db.CloseModusDb(ctx)
	})

	// Test initialization with build path
	t.Run("InitModusDb with build path", func(t *testing.T) {
		buildPath := filepath.Join(testAppPath, "build")
		err := os.MkdirAll(buildPath, 0755)
		require.NoError(t, err)

		// Set app path to build directory
		cfg := app.NewAppConfig().WithAppPath(buildPath)
		app.SetConfig(cfg)

		db.InitModusDb(ctx, "")

		// Verify .modusdb directory was created in parent directory
		modusDbPath := filepath.Join(testAppPath, ".modusdb")
		_, err = os.Stat(modusDbPath)
		assert.NoError(t, err, "ModusDB directory should exist in parent directory")

		// Verify .gitignore was created
		gitIgnorePath := filepath.Join(testAppPath, ".gitignore")
		_, err = os.Stat(gitIgnorePath)
		assert.NoError(t, err, ".gitignore should exist")

		// Verify .gitignore contains .modusdb/
		content, err := os.ReadFile(gitIgnorePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), ".modusdb/", ".gitignore should contain .modusdb/")

		db.CloseModusDb(ctx)
	})

	// Test GetClient thread safety
	t.Run("GetClient thread safety", func(t *testing.T) {
		db.InitModusDb(ctx, "file://"+t.TempDir())
		client1, err := db.GetClient()
		require.NoError(t, err)
		defer client1.Close()
		assert.NotNil(t, client1, "First client should not be nil")

		client2, err := db.GetClient()
		require.NoError(t, err)
		defer client2.Close()
		assert.NotNil(t, client2, "Second client should not be nil")
		assert.Equal(t, client1, client2, "GetClient should return the same instance")

		db.CloseModusDb(ctx)
	})

	// Out of order and lifecycle safety
	t.Run("CloseModusDb", func(t *testing.T) {
		client, err := db.GetClient()
		require.Error(t, err)
		require.Nil(t, client, "Client should be nil if the engine is not initialized")

		db.InitModusDb(ctx, "file://"+t.TempDir())
		db.CloseModusDb(ctx)

		client, err = db.GetClient()
		require.Error(t, err)
		require.Nil(t, client, "Client should be nil if the engine is closed")
	})
}
