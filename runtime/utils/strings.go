/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

func DecodeUTF16(bytes []byte) string {

	// Make sure the buffer is valid.
	if len(bytes) == 0 || len(bytes)%2 != 0 {
		return ""
	}

	// Reinterpret []byte as []uint16 to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&bytes[0])
	words := unsafe.Slice((*uint16)(ptr), len(bytes)/2)

	// Decode UTF-16 words to a UTF-8 string.
	str := string(utf16.Decode(words))
	return str
}

func EncodeUTF16(str string) []byte {
	if len(str) == 0 {
		return []byte{}
	}
	// Encode the UTF-8 string to UTF-16 words.
	words := utf16.Encode([]rune(str))

	// Reinterpret []uint16 as []byte to avoid excess copying.
	// This works because we can presume the system is little-endian.
	ptr := unsafe.Pointer(&words[0])
	bytes := unsafe.Slice((*byte)(ptr), len(words)*2)
	return bytes
}

// SanitizeUTF8 removes invalid UTF-8 sequences from a byte slice.
// It skips over any byte that is a null byte (0) or a single-byte character
// that is not part of a valid UTF-8 sequence.
// It returns a new byte slice containing only valid UTF-8 characters.
func SanitizeUTF8(s []byte) []byte {
	// This is adapted from bytes.ToValidUTF8
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		c := s[i]
		if c == 0 {
			i++
			continue
		}
		if c < utf8.RuneSelf {
			i++
			b = append(b, c)
			continue
		}
		_, wid := utf8.DecodeRune(s[i:])
		if wid == 1 {
			i++
			continue
		}
		b = append(b, s[i:i+wid]...)
		i += wid
	}
	return b
}

func TrimStringBefore(s string, sep string) string {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return s
}
