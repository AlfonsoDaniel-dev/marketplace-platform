package helpers

import (
	"database/sql"
	"time"
)

func TimeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{
		Time: t,
	}
	if !null.Time.IsZero() {
		null.Valid = true
	}

	return null
}

func IntToNull(number int64) sql.NullInt64 {
	null := sql.NullInt64{
		Int64: number,
	}

	if null.Int64 != 0 {
		null.Valid = true
	}

	return null
}

func StringToNull(char string) sql.NullString {
	null := sql.NullString{
		String: char,
	}

	if null.String != "" {
		null.Valid = true
	}

	return null
}
