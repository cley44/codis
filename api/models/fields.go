package models

import (
	"database/sql/driver"
	"encoding/json"
)

// PGSlice is a custom type for handling slices stored as JSONB in Postgres.
type PGSlice[T any] []T

// https://golang.org/pkg/database/sql/driver/#Valuer implementation for Postrgres String Array
func (slice PGSlice[T]) Value() (driver.Value, error) {
	return json.Marshal(slice)
}

// https://golang.org/pkg/database/sql/#Scanner implementation for Postrgres String Array
func (slice *PGSlice[T]) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &slice)
	}

	return nil
}
