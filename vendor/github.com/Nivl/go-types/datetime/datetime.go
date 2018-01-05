package datetime

import (
	"database/sql/driver"
	"time"
)

// ISO8601 is a time.Time layout for the ISO8601 format
const ISO8601 = "2006-01-02T15:04:05-0700"

// DateTime represents a time.Time that uses ISO8601 for json input/output
// instead of RFC3339
type DateTime struct {
	time.Time
}

// Now returns the current UTC time.
func Now() *DateTime {
	return &DateTime{Time: time.Now().UTC()}
}

// Value returns a value that the database can handle
// https://golang.org/pkg/database/sql/driver/#Valuer
func (t *DateTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.UTC().Format(ISO8601), nil
}

// Scan assigns a value from a database driver
// https://golang.org/pkg/database/sql/#Scanner
func (t *DateTime) Scan(value interface{}) error {
	if value != nil {
		t.Time = value.(time.Time)
		t.Time = t.Time.UTC()
	}
	return nil
}

// MarshalJSON returns a valid json representation of the struct
// https://golang.org/pkg/encoding/json/#Marshaler
func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.UTC().Format(ISO8601) + `"`), nil
}

// UnmarshalJSON tries to parse a json data into a valid struct
// https://golang.org/pkg/encoding/json/#Unmarshaler
func (t *DateTime) UnmarshalJSON(data []byte) error {
	var err error
	t.Time, err = time.Parse(`"`+ISO8601+`"`, string(data))
	if err != nil {
		return err
	}
	t.Time = t.Time.UTC()
	return nil
}

// Equal check if the given date is equal to the current one
func (t *DateTime) Equal(u *DateTime) bool {
	return t.Time.Equal(u.Time)
}

// AddDate returns the time corresponding to adding the given number of years, months, and days to t.
func (t *DateTime) AddDate(years int, months int, days int) *DateTime {
	return &DateTime{Time: t.Time.AddDate(years, months, days)}
}
