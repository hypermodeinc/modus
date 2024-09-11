/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"hypruntime/langsupport"
	"reflect"
	"testing"
	"time"
	"unsafe"
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 4 {
		t.Errorf("expected type size 4, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 4 {
		t.Errorf("expected type size 4, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "int"
	rtInner := reflect.TypeFor[int]()

	innerInfo := innerHandlers[0].Info()
	if innerInfo.TypeName() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.TypeName())
	}
	if innerInfo.TypeSize() != 4 {
		t.Errorf("expected inner type size 4, got %d", innerInfo.TypeSize())
	}
	if innerInfo.RuntimeType() != rtInner {
		t.Errorf("expected inner runtime type %v, got %v", rtInner, innerInfo.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 8 {
		t.Errorf("expected type size 8, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 4 {
		t.Errorf("expected type size 4, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].Info()
	if innerInfo.TypeName() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.TypeName())
	}
	if innerInfo.TypeSize() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.TypeSize())
	}
	if innerInfo.RuntimeType() != rtInner {
		t.Errorf("expected inner runtime type %v, got %v", rtInner, innerInfo.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 12 {
		t.Errorf("expected type size 12, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].Info()
	if innerInfo.TypeName() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.TypeName())
	}
	if innerInfo.TypeSize() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.TypeSize())
	}
	if innerInfo.RuntimeType() != rtInner {
		t.Errorf("expected inner runtime type %v, got %v", rtInner, innerInfo.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 16 {
		t.Errorf("expected type size 16, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 1 {
		t.Fatalf("expected 1 inner handler, got %d", len(innerHandlers))
	}

	typInner := "string"
	rtInner := reflect.TypeFor[string]()

	innerInfo := innerHandlers[0].Info()
	if innerInfo.TypeName() != typInner {
		t.Errorf("expected inner type name %q, got %q", typInner, innerInfo.TypeName())
	}
	if innerInfo.TypeSize() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo.TypeSize())
	}
	if innerInfo.RuntimeType() != rtInner {
		t.Errorf("expected inner runtime type %v, got %v", rtInner, innerInfo.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 20 {
		t.Errorf("expected type size 20, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 8 {
		t.Errorf("expected type size 8, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 4 {
		t.Errorf("expected type size 4, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 2 {
		t.Fatalf("expected 2 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "[]string"
	rtInner0 := reflect.TypeFor[[]string]()

	innerInfo0 := innerHandlers[0].Info()
	if innerInfo0.TypeName() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.TypeName())
	}
	if innerInfo0.TypeSize() != 12 {
		t.Errorf("expected inner type size 12, got %d", innerInfo0.TypeSize())
	}
	if innerInfo0.RuntimeType() != rtInner0 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner0, innerInfo0.RuntimeType())
	}

	typInner1 := "[]string"
	rtInner1 := reflect.TypeFor[[]string]()

	innerInfo1 := innerHandlers[1].Info()
	if innerInfo1.TypeName() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.TypeName())
	}
	if innerInfo1.TypeSize() != 12 {
		t.Errorf("expected inner type size 12, got %d", innerInfo1.TypeSize())
	}
	if innerInfo1.RuntimeType() != rtInner1 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner1, innerInfo1.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.TypeSize() != 16 {
		t.Errorf("expected type size 16, got %d", info.TypeSize())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 3 {
		t.Fatalf("expected 3 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "bool"
	rtInner0 := reflect.TypeFor[bool]()

	innerInfo0 := innerHandlers[0].Info()
	if innerInfo0.TypeName() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.TypeName())
	}
	if innerInfo0.TypeSize() != 1 {
		t.Errorf("expected inner type size 1, got %d", innerInfo0.TypeSize())
	}
	if innerInfo0.RuntimeType() != rtInner0 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner0, innerInfo0.RuntimeType())
	}

	typInner1 := "int"
	rtInner1 := reflect.TypeFor[int]()

	innerInfo1 := innerHandlers[1].Info()
	if innerInfo1.TypeName() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.TypeName())
	}
	if innerInfo1.TypeSize() != 4 {
		t.Errorf("expected inner type size 4, got %d", innerInfo1.TypeSize())
	}
	if innerInfo1.RuntimeType() != rtInner1 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner1, innerInfo1.RuntimeType())
	}

	typInner2 := "string"
	rtInner2 := reflect.TypeFor[string]()

	innerInfo2 := innerHandlers[2].Info()
	if innerInfo2.TypeName() != typInner2 {
		t.Errorf("expected inner type name %q, got %q", typInner2, innerInfo2.TypeName())
	}
	if innerInfo2.TypeSize() != 8 {
		t.Errorf("expected inner type size 8, got %d", innerInfo2.TypeSize())
	}
	if innerInfo2.RuntimeType() != rtInner2 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner2, innerInfo2.RuntimeType())
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

	info := handler.Info()
	if info.TypeName() != typ {
		t.Errorf("expected type name %q, got %q", typ, info.TypeName())
	}
	if info.RuntimeType() != rt {
		t.Errorf("expected runtime type %v, got %v", rt, info.RuntimeType())
	}

	innerHandlers := getInnerHandlers(handler)
	if len(innerHandlers) != 2 {
		t.Fatalf("expected 2 inner handlers, got %d", len(innerHandlers))
	}

	typInner0 := "bool"
	rtInner0 := reflect.TypeFor[bool]()

	innerInfo0 := innerHandlers[0].Info()
	if innerInfo0.TypeName() != typInner0 {
		t.Errorf("expected inner type name %q, got %q", typInner0, innerInfo0.TypeName())
	}
	if innerInfo0.RuntimeType() != rtInner0 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner0, innerInfo0.RuntimeType())
	}

	typInner1 := "*testdata.TestRecursiveStruct"
	rtInner1 := reflect.TypeFor[*TestRecursiveStruct]()

	innerInfo1 := innerHandlers[1].Info()
	if innerInfo1.TypeName() != typInner1 {
		t.Errorf("expected inner type name %q, got %q", typInner1, innerInfo1.TypeName())
	}
	if innerInfo1.RuntimeType() != rtInner1 {
		t.Errorf("expected inner runtime type %v, got %v", rtInner1, innerInfo1.RuntimeType())
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
