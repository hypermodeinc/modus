/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"errors"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewCleaner(t *testing.T) {
	c := utils.NewCleaner()

	assert.NotNil(t, c)
	assert.Equal(t, 0, c.Len())
	assert.Equal(t, 0, c.Cap())
}

func TestNewCleanerN(t *testing.T) {
	c := utils.NewCleanerN(10)

	assert.NotNil(t, c)
	assert.Equal(t, 0, c.Len())
	assert.Equal(t, 10, c.Cap())
}

func TestCleaner_AddCleanup(t *testing.T) {
	c := utils.NewCleaner()
	c.AddCleanup(func() error { return nil })
	c.AddCleanup(func() error { return nil })

	assert.Equal(t, 2, c.Len())
}

func TestCleaner_AddCleaner(t *testing.T) {
	c1 := utils.NewCleaner()
	c1.AddCleanup(func() error { return nil })
	c1.AddCleanup(func() error { return nil })

	c2 := utils.NewCleaner()
	c2.AddCleanup(func() error { return nil })

	c1.AddCleaner(c2)
	assert.Equal(t, 3, c1.Len())

	c1.AddCleaner(nil)
	assert.Equal(t, 3, c1.Len())

	c3 := utils.NewCleaner()
	c3.AddCleaner(c2)
	assert.Equal(t, 1, c3.Len())
}

func TestCleaner_Clean(t *testing.T) {
	c := utils.NewCleaner()
	assert.Nil(t, c.Clean())

	a := false
	c.AddCleanup(func() error { a = true; return nil })
	assert.Nil(t, c.Clean())
	assert.True(t, a)

	a = false
	b := false
	c.AddCleanup(func() error { b = true; return nil })
	assert.Nil(t, c.Clean())
	assert.True(t, a)
	assert.True(t, b)

	c = utils.NewCleaner()
	c.AddCleanup(func() error { return errors.New("error") })
	assert.NotNil(t, c.Clean())

	c = utils.NewCleaner()
	c.AddCleanup(func() error { return nil })
	c.AddCleanup(func() error { return errors.New("error") })
	assert.NotNil(t, c.Clean())

	c = utils.NewCleaner()
	c.AddCleanup(func() error { return errors.New("error") })
	c.AddCleanup(func() error { return errors.New("error") })
	assert.NotNil(t, c.Clean())
}
