package sqlite

import (
	"database/sql"
	"time"
)

// SQLite stores timestamps as TEXT. These helpers handle scanning.

// Common SQLite datetime formats to try when parsing.
var timeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	time.RFC3339,
	time.RFC3339Nano,
}

// parseTime parses a SQLite TEXT timestamp string into time.Time.
func parseTime(s string) (time.Time, error) {
	for _, fmt := range timeFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &time.ParseError{Value: s, Message: "no matching time format"}
}

// nullableTimeStr is a helper for scanning nullable time columns from SQLite.
type nullableTimeStr struct {
	sql.NullString
}

// Time returns the parsed time or nil if NULL.
func (n nullableTimeStr) Time() *time.Time {
	if !n.Valid || n.String == "" {
		return nil
	}
	t, err := parseTime(n.String)
	if err != nil {
		return nil
	}
	return &t
}
