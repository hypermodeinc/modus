/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"errors"
	"fmt"
)

type Cleaner interface {
	Len() int
	Cap() int
	Clean() error
	AddCleanup(fn func() error)
	AddCleaner(c Cleaner)
}

type cleaner struct {
	cleanupFuncs []func() error
}

func NewCleaner() Cleaner {
	return &cleaner{}
}

func NewCleanerN(capacity int) Cleaner {
	return &cleaner{
		cleanupFuncs: make([]func() error, 0, capacity),
	}
}

func (c *cleaner) Len() int {
	return len(c.cleanupFuncs)
}

func (c *cleaner) Cap() int {
	return cap(c.cleanupFuncs)
}

func (c *cleaner) AddCleanup(fn func() error) {
	c.cleanupFuncs = append(c.cleanupFuncs, fn)
}

func (c *cleaner) AddCleaner(cln Cleaner) {
	if cln == nil {
		return
	}

	if len(c.cleanupFuncs) == 0 {
		c.cleanupFuncs = cln.(*cleaner).cleanupFuncs
	} else {
		c.AddCleanup(cln.Clean)
	}
}

func (c *cleaner) Clean() error {
	numFuncs := len(c.cleanupFuncs)
	switch numFuncs {
	case 0:
		return nil
	case 1:
		if err := c.cleanupFuncs[0](); err != nil {
			return fmt.Errorf("failed to cleanup: %w", err)
		}
		return nil
	}

	errs := make([]error, 0, numFuncs)
	for i := numFuncs - 1; i >= 0; i-- {
		if fn := c.cleanupFuncs[i]; fn != nil {
			if err := fn(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("failed to cleanup: %w", errs[0])
	default:
		return fmt.Errorf("failed to cleanup: %w", errors.Join(errs...))
	}
}
