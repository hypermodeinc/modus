/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"bytes"
	"fmt"
	"time"
)

// TimeFormat is the default time format template throughout the application.
// It is used for output to logs, query results, etc.
const TimeFormat = "2006-01-02T15:04:05.000Z"

// JSONTime extends the time type so that when it marshals to JSON
// the output is fixed to include milliseconds as three digits.
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 26))
	buf.WriteString("\"")
	buf.WriteString(t.String())
	buf.WriteString("\"")
	return buf.Bytes(), nil
}

func (t JSONTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

// Allow date, time, date-time formats.
// Seconds, milliseconds, and timezone are optional when including time.
// If time zone is not included, it is assumed to be UTC.
// If time is not included, it is assumed to be midnight.
// Allow space or T as separator between date and time.
var inputTimeFormats = []string{
	"2006-01-02T15:04:05.999Z07:00",
	"2006-01-02T15:04Z07:00",
	"2006-01-02T15:04:05.999",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05.999Z07:00",
	"2006-01-02 15:04Z07:00",
	"2006-01-02 15:04:05.999",
	"2006-01-02 15:04",
	"2006-01-02",
}

// ParseTime parses a string into a time.Time object.
// It tries multiple formats to parse the string.
func ParseTime(s string) (time.Time, error) {

	// TODO: This is a naive implementation.
	// We should consider separation of multiple date and time formats,
	// and possibly use a more efficient parsing library.

	for _, format := range inputTimeFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse date time string: %s", s)
}
