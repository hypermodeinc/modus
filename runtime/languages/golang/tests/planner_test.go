/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang_test

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/hypermodeinc/modus/runtime/langsupport"
)

func TestGetHandler_int(t *testing.T) {
	typ := "int"
	rt := reflect.TypeFor[int]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 1 {
		t.Fatalf("expected 1 handler, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 4 {
		t.Errorf("expected type size 4, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}
}

func TestGetHandler_intPtr(t *testing.T) {
	typ := "*int"
	rt := reflect.TypeFor[*int]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 2 {
		t.Fatalf("expected 2 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 4 {
		t.Errorf("expected type size 4, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "int"
	rtInner := reflect.TypeFor[int]()

	innerInfo := innerHandlers[0].TypeInfo()
	if innerInfo.Name() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.Name())
	}
	if innerInfo.Size() != 4 {
		t.Errorf("expected inner type size 4, got %d", innerInfo.Size())
	}
	if innerInfo.ReflectedType() != rtInner {
		t.Errorf("expected inner reflected type %v, got %v", rtInner, innerInfo.ReflectedType())
	}
}

func TestGetHandler_string(t *testing.T) {
	typ := "string"
	rt := reflect.TypeFor[string]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 1 {
		t.Fatalf("expected 1 handler, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 8 {
		t.Errorf("expected type size 8, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}
}

func TestGetHandler_stringPtr(t *testing.T) {
	typ := "*string"
	rt := reflect.TypeFor[*string]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 2 {
		t.Fatalf("expected 2 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 4 {
		t.Errorf("expected type size 4, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].TypeInfo()
	if innerInfo.Name() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.Name())
	}
	if innerInfo.Size() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.Size())
	}
	if innerInfo.ReflectedType() != rtInner {
		t.Errorf("expected inner reflected type %v, got %v", rtInner, innerInfo.ReflectedType())
	}
}

func TestGetHandler_stringSlice(t *testing.T) {
	typ := "[]string"
	rt := reflect.TypeFor[[]string]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 2 {
		t.Fatalf("expected 2 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 12 {
		t.Errorf("expected type size 12, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].TypeInfo()
	if innerInfo.Name() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.Name())
	}
	if innerInfo.Size() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.Size())
	}
	if innerInfo.ReflectedType() != rtInner {
		t.Errorf("expected inner reflected type %v, got %v", rtInner, innerInfo.ReflectedType())
	}
}

func TestGetHandler_stringArray(t *testing.T) {
	typ := "[2]string"
	rt := reflect.TypeFor[[2]string]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 2 {
		t.Fatalf("expected 2 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 16 {
		t.Errorf("expected type size 16, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].TypeInfo()
	if innerInfo.Name() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.Name())
	}
	if innerInfo.Size() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.Size())
	}
	if innerInfo.ReflectedType() != rtInner {
		t.Errorf("expected inner reflected type %v, got %v", rtInner, innerInfo.ReflectedType())
	}
}

func TestGetHandler_time(t *testing.T) {
	typ := "time.Time"
	rt := reflect.TypeFor[time.Time]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 1 {
		t.Fatalf("expected 1 handler, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 20 {
		t.Errorf("expected type size 20, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}
}

func TestGetHandler_duration(t *testing.T) {
	typ := "time.Duration"
	rt := reflect.TypeFor[time.Duration]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 1 {
		t.Fatalf("expected 1 handler, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 8 {
		t.Errorf("expected type size 8, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}
}

func TestGetHandler_map(t *testing.T) {
	typ := "map[string]string"
	rt := reflect.TypeFor[map[string]string]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 3 {
		t.Fatalf("expected 3 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 4 {
		t.Errorf("expected type size 4, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 2 {
		t.Fatalf("expected 2 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "[]string"
	rtInner0 := reflect.TypeFor[[]string]()

	innerInfo0 := innerHandlers[0].TypeInfo()
	if innerInfo0.Name() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.Name())
	}
	if innerInfo0.Size() != 12 {
		t.Errorf("expected inner type size 12, got %d", innerInfo0.Size())
	}
	if innerInfo0.ReflectedType() != rtInner0 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner0, innerInfo0.ReflectedType())
	}

	typInner1 := "[]string"
	rtInner1 := reflect.TypeFor[[]string]()

	innerInfo1 := innerHandlers[1].TypeInfo()
	if innerInfo1.Name() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.Name())
	}
	if innerInfo1.Size() != 12 {
		t.Errorf("expected inner type size 12, got %d", innerInfo1.Size())
	}
	if innerInfo1.ReflectedType() != rtInner1 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner1, innerInfo1.ReflectedType())
	}
}

func TestGetHandler_struct(t *testing.T) {
	typ := "testdata.TestStruct3"
	rt := reflect.TypeFor[TestStruct3]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 4 {
		t.Fatalf("expected 4 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.Size() != 16 {
		t.Errorf("expected type size 16, got %d", info.Size())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 3 {
		t.Fatalf("expected 3 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "bool"
	rtInner0 := reflect.TypeFor[bool]()

	innerInfo0 := innerHandlers[0].TypeInfo()
	if innerInfo0.Name() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.Name())
	}
	if innerInfo0.Size() != 1 {
		t.Errorf("expected inner type size 1, got %d", innerInfo0.Size())
	}
	if innerInfo0.ReflectedType() != rtInner0 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner0, innerInfo0.ReflectedType())
	}

	typInner1 := "int"
	rtInner1 := reflect.TypeFor[int]()

	innerInfo1 := innerHandlers[1].TypeInfo()
	if innerInfo1.Name() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.Name())
	}
	if innerInfo1.Size() != 4 {
		t.Errorf("expected inner type size 4, got %d", innerInfo1.Size())
	}
	if innerInfo1.ReflectedType() != rtInner1 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner1, innerInfo1.ReflectedType())
	}

	typInner2 := "string"
	rtInner2 := reflect.TypeFor[string]()

	innerInfo2 := innerHandlers[2].TypeInfo()
	if innerInfo2.Name() != typInner2 {
		t.Errorf("expected inner type name %q, got %q", typInner2, innerInfo2.Name())
	}
	if innerInfo2.Size() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo2.Size())
	}
	if innerInfo2.ReflectedType() != rtInner2 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner2, innerInfo2.ReflectedType())
	}
}

func TestGetHandler_recursiveStruct(t *testing.T) {
	typ := "testdata.TestRecursiveStruct"
	rt := reflect.TypeFor[TestRecursiveStruct]()

	planner := fixture.NewPlanner()
	handler, err := planner.GetHandler(fixture.Context, typ)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalHandlers := len(planner.AllHandlers())
	if totalHandlers != 3 {
		t.Fatalf("expected 3 handlers, got %d", totalHandlers)
	}

	info := handler.TypeInfo()
	if info.Name() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.Name())
	}
	if info.ReflectedType() != rt {
		t.Errorf("expected reflected type %v, got %v", rt, info.ReflectedType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 2 {
		t.Fatalf("expected 2 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "bool"
	rtInner0 := reflect.TypeFor[bool]()

	innerInfo0 := innerHandlers[0].TypeInfo()
	if innerInfo0.Name() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.Name())
	}
	if innerInfo0.ReflectedType() != rtInner0 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner0, innerInfo0.ReflectedType())
	}

	typInner1 := "*testdata.TestRecursiveStruct"
	rtInner1 := reflect.TypeFor[*TestRecursiveStruct]()

	innerInfo1 := innerHandlers[1].TypeInfo()
	if innerInfo1.Name() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.Name())
	}
	if innerInfo1.ReflectedType() != rtInner1 {
		t.Errorf("expected inner reflected type %v, got %v", rtInner1, innerInfo1.ReflectedType())
	}
}

var rtTypeHandler = reflect.TypeFor[langsupport.TypeHandler]()

func getInnerHandlers(handler langsupport.TypeHandler) []langsupport.TypeHandler {
	var results []langsupport.TypeHandler
	rvHandler := reflect.ValueOf(handler).Elem()
	for i := 0; i < rvHandler.NumField(); i++ {
		rf := rvHandler.Field(i)
		field := reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
		if field.Type().Implements(rtTypeHandler) {
			results = append(results, field.Interface().(langsupport.TypeHandler))
		} else if field.Kind() == reflect.Slice && field.Type().Elem().Implements(rtTypeHandler) {
			for j := 0; j < field.Len(); j++ {
				results = append(results, field.Index(j).Interface().(langsupport.TypeHandler))
			}
		}
	}
	return results
}
