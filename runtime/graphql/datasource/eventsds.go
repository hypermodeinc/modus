/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package datasource

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"

	"github.com/hypermodeinc/modus/runtime/actors"
	"github.com/hypermodeinc/modus/runtime/utils"
)

type eventsDataSource struct{}

func (s *eventsDataSource) Start(rc *resolve.Context, input []byte, updater resolve.SubscriptionUpdater) error {
	err := startEventsDataSource(rc, input, updater)
	if err != nil {
		// We should be able to return the error, but there appears to be a bug in the GraphQL-Go-Tools library
		// where it doesn't handle the error correctly and causes a panic.
		// TODO: investigate this further.
		errMsg, _ := utils.JsonSerialize(err.Error())
		updater.Update(fmt.Appendf(nil, `{"errors":[{"message":%s}]}`, errMsg))
		updater.Close(resolve.SubscriptionCloseKindDownstreamServiceError)
	}
	return nil
}

func startEventsDataSource(rc *resolve.Context, input []byte, updater resolve.SubscriptionUpdater) error {
	var ci callInfo
	if err := utils.JsonDeserialize(input, &ci); err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	// Current there is only one hardcoded subscription for agent events.
	if subName := ci.FieldInfo.Name; subName != "agentEvent" {
		return fmt.Errorf("unknown subscription: %s", subName)
	}

	var agentId string
	if value, ok := ci.Parameters["agentId"]; !ok {
		return fmt.Errorf("missing required parameter 'agentId'")
	} else if id, ok := value.(string); !ok {
		return fmt.Errorf("invalid type for 'agentId', expected string, got %T", value)
	} else {
		agentId = id
	}

	fieldName := ci.FieldInfo.AliasOrName()
	return actors.SubscribeForAgentEvents(
		rc.Context(),
		agentId,
		func(data []byte) {
			updater.Update(fmt.Appendf(nil, `{"data":{"%s":%s}}`, fieldName, data))
		},
		updater.Complete,
	)
}

func (s *eventsDataSource) UniqueRequestID(rc *resolve.Context, input []byte, xxh *xxhash.Digest) error {
	_, err := xxh.Write(input)
	return err
}
