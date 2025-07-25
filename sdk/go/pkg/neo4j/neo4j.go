/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

type Neo4jOption func(*neo4jOptions)

type neo4jOptions struct {
	dbName string
}

func WithDbName(dbName string) Neo4jOption {
	return func(o *neo4jOptions) {
		o.dbName = dbName
	}
}

// The result of a Neo4j query.
type EagerResult struct {
	Keys    []string
	Records []*Record
}

// A record in a Neo4j query result.
type Record struct {
	Keys   []string
	Values []string
}

// A type constraint for retrieving values from a Neo4j record.
type RecordValue interface {
	bool | int64 | float64 | string | time.Time |
		[]byte | []any | map[string]any |
		Node | Relationship | Path | Point2D | Point3D
}

type Entity interface {
	GetElementId() string
	GetProperties() map[string]any
}

type Node struct {
	ElementId string
	Labels    []string
	Props     map[string]any
}

func (n *Node) GetElementId() string {
	return n.ElementId
}

func (n *Node) GetProperties() map[string]any {
	return n.Props
}

// A relationship in a Neo4j query result.
type Relationship struct {
	ElementId      string
	StartElementId string
	EndElementId   string
	Type           string
	Props          map[string]any
}

func (r *Relationship) GetElementId() string {
	return r.ElementId
}

func (r *Relationship) GetProperties() map[string]any {
	return r.Props
}

type Path struct {
	Nodes         []Node
	Relationships []Relationship
}

// A type constraint for retrieving property values from a Neo4j entity.
type PropertyValue interface {
	bool | int64 | float64 | string | time.Time |
		[]byte | []any | Point2D | Point3D
}

// Point2D represents a two dimensional point in a particular coordinate reference system.
type Point2D struct {
	X float64
	Y float64

	// Id of coordinate reference system.
	SpatialRefId uint32
}

// Point3D represents a three dimensional point in a particular coordinate reference system.
type Point3D struct {
	X float64
	Y float64
	Z float64

	// Id of coordinate reference system.
	SpatialRefId uint32
}

// String returns string representation of this point.
func (p Point2D) String() string {
	return fmt.Sprintf("Point{SpatialRefId=%d, X=%f, Y=%f}", p.SpatialRefId, p.X, p.Y)
}

// String returns string representation of this point.
func (p Point3D) String() string {
	return fmt.Sprintf("Point{SpatialRefId=%d, X=%f, Y=%f, Z=%f}", p.SpatialRefId, p.X, p.Y, p.Z)
}

// Executes a query or mutation on the Neo4j database.
func ExecuteQuery(connection, query string, parameters map[string]any, opts ...Neo4jOption) (*EagerResult, error) {
	dbOpts := &neo4jOptions{
		dbName: "neo4j",
	}

	for _, opt := range opts {
		opt(dbOpts)
	}

	bytes, err := utils.JsonSerialize(parameters)
	if err != nil {
		return nil, err
	}

	parametersJson := string(bytes)

	response := hostExecuteQuery(&connection, &dbOpts.dbName, &query, &parametersJson)

	return response, nil
}

// Get a value from a record at a given key and cast or decode it to a specific type.
func GetRecordValue[T RecordValue](record *Record, key string) (T, error) {
	var val T
	for i, k := range record.Keys {
		if k == key {
			err := json.Unmarshal([]byte(record.Values[i]), &val)
			if err != nil {
				return *new(T), err
			} else {
				return val, nil
			}
		}
	}
	return *new(T), fmt.Errorf("key not found in record")

}

func (r *Record) Get(key string) (string, bool) {
	for i, k := range r.Keys {
		if k == key {
			return r.Values[i], true
		}
	}
	return "", false
}

func (r *Record) AsMap() map[string]string {
	result := make(map[string]string)
	for i, k := range r.Keys {
		result[k] = r.Values[i]
	}
	return result
}

func (r *Record) JSONMarshal() ([]byte, error) {
	result := "{"
	for i, k := range r.Keys {
		keyBytes, err := utils.JsonSerialize(k)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf("%s:%s", keyBytes, r.Values[i])
		if i < len(r.Keys)-1 {
			result += ","
		}
	}
	result += "}"
	return []byte(result), nil
}

// Get a property from an entity at a given key and cast or decode it to a specific type.
func GetProperty[T PropertyValue](e Entity, key string) (T, error) {
	var val T
	rawVal, ok := e.GetProperties()[key]
	if !ok {
		return *new(T), fmt.Errorf("key not found in node")
	}
	switch any(val).(type) {
	case int64:
		float64Val, ok := rawVal.(float64)
		if !ok {
			return *new(T), fmt.Errorf("expected value to have type int64 but found type %T", rawVal)
		}
		return any(int64(float64Val)).(T), nil
	default:
		val, ok = rawVal.(T)
		if !ok {
			zeroValue := *new(T)
			return zeroValue, fmt.Errorf("expected value to have type %T but found type %T", zeroValue, rawVal)
		}
		return val, nil
	}

}
