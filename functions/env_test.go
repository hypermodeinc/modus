/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"hmruntime/testutils"
	"testing"
	"time"
)

func Test_DateNow(t *testing.T) {

	startTs := time.Now().UnixMilli()

	f := testutils.NewWasmTestFixture()
	defer f.Close()

	fn := f.Module.ExportedFunction("now")
	res, err := fn.Call(f.Context)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 {
		t.Errorf("expected 1 result, got %d", len(res))
	}

	resultTs := int64(res[0])

	if resultTs <= startTs {
		t.Errorf("expected timestamp greater than %d, got %d", startTs, resultTs)
	}
}
