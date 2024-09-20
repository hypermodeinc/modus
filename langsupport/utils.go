/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

// AlignOffset returns the smallest y >= x such that y % a == 0.
func AlignOffset(x, a uint32) uint32 {
	return (x + a - 1) &^ (a - 1)
}
