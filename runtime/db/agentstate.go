/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/hypermodeinc/modusgraph"
	mg_utils "github.com/hypermodeinc/modusgraph/api/apiutils"
	"github.com/jackc/pgx/v5"
)

var ErrAgentNotFound = fmt.Errorf("agent not found")

type AgentState struct {
	Gid       uint64 `json:"gid,omitempty"`
	Id        string `json:"id" db:"constraint=unique"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Data      string `json:"data,omitempty"`
	UpdatedAt string `json:"updated"`
}

func WriteAgentState(ctx context.Context, state AgentState) error {
	if useModusDB() {
		return writeAgentStateToModusDB(ctx, state)
	} else {
		return writeAgentStateToPostgresDB(ctx, state)
	}
}

func UpdateAgentStatus(ctx context.Context, id string, status string) error {
	if useModusDB() {
		return updateAgentStatusInModusDB(ctx, id, status)
	} else {
		return updateAgentStatusInPostgresDB(ctx, id, status)
	}
}

func GetAgentState(ctx context.Context, id string) (*AgentState, error) {
	if useModusDB() {
		return getAgentStateFromModusDB(ctx, id)
	} else {
		return getAgentStateFromPostgresDB(ctx, id)
	}
}

func QueryActiveAgents(ctx context.Context) ([]AgentState, error) {
	if useModusDB() {
		return queryActiveAgentsFromModusDB(ctx)
	} else {
		return queryActiveAgentsFromPostgresDB(ctx)
	}
}

func writeAgentStateToModusDB(ctx context.Context, state AgentState) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	gid, _, _, err := modusgraph.Upsert(ctx, globalModusDbEngine, state)
	state.Gid = gid

	return err
}

func updateAgentStatusInModusDB(ctx context.Context, id string, status string) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	// TODO: this should just be an update in a single operation

	state, err := getAgentStateFromModusDB(ctx, id)
	if err != nil {
		return err
	}

	state.Status = status
	state.UpdatedAt = time.Now().UTC().Format(utils.TimeFormat)

	return writeAgentStateToModusDB(ctx, *state)
}

func getAgentStateFromModusDB(ctx context.Context, id string) (*AgentState, error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	_, result, err := modusgraph.Get[AgentState](ctx, globalModusDbEngine, modusgraph.ConstrainedField{
		Key:   "id",
		Value: id,
	})
	if errors.Is(err, mg_utils.ErrNoObjFound) {
		return nil, ErrAgentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query agent state: %w", err)
	}

	return &result, nil
}

func queryActiveAgentsFromModusDB(ctx context.Context) ([]AgentState, error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	_, results, err := modusgraph.Query[AgentState](ctx, globalModusDbEngine, modusgraph.QueryParams{
		Filter: &modusgraph.Filter{
			Not: &modusgraph.Filter{
				Field: "status",
				String: modusgraph.StringPredicate{
					Equals: "terminated",
				},
			},
		},
		// TODO: Sorting gives a dgraph error. Why?
		// Sorting: &modusgraph.Sorting{
		// 	OrderDescField: "updated",
		// 	OrderDescFirst: true,
		// },
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query agent state: %w", err)
	}

	return results, nil
}

func writeAgentStateToPostgresDB(ctx context.Context, state AgentState) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	const query = "INSERT INTO agents (id, name, status, data, updated) VALUES ($1, $2, $3, $4, $5) " +
		"ON CONFLICT (id) DO UPDATE SET name = $2, status = $3, data = $4, updated = $5"

	err := WithTx(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, query, state.Id, state.Name, state.Status, state.Data, state.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to write agent state: %w", err)
		}
		return nil
	})

	return err
}

func updateAgentStatusInPostgresDB(ctx context.Context, id string, status string) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	const query = "UPDATE agents SET status = $2, updated = $3 WHERE id = $1"
	now := time.Now().UTC()

	err := WithTx(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, query, id, status, now)
		if err != nil {
			return fmt.Errorf("failed to update agent status: %w", err)
		}
		return nil
	})

	return err
}

func getAgentStateFromPostgresDB(ctx context.Context, id string) (*AgentState, error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	const query = "SELECT id, name, status, data, updated FROM agents WHERE id = $1"

	var a AgentState
	var ts time.Time
	err := WithTx(ctx, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, query, id)
		if err := row.Scan(&a.Id, &a.Name, &a.Status, &a.Data, &ts); err != nil {
			if err == pgx.ErrNoRows {
				return ErrAgentNotFound
			}
			return fmt.Errorf("failed to get agent state: %w", err)
		}
		a.UpdatedAt = ts.UTC().Format(utils.TimeFormat)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func queryActiveAgentsFromPostgresDB(ctx context.Context) ([]AgentState, error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	const query = "SELECT id, name, status, data, updated FROM agents " +
		"WHERE status != 'terminated' ORDER BY updated DESC"

	results := make([]AgentState, 0)
	err := WithTx(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var a AgentState
			var ts time.Time
			if err := rows.Scan(&a.Id, &a.Name, &a.Status, &a.Data, &ts); err != nil {
				return err
			}
			a.UpdatedAt = ts.UTC().Format(utils.TimeFormat)
			results = append(results, a)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})

	return results, err
}
