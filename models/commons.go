package models

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

var b = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)

// PageSize defines the number of row of a paginated query
const PageSize = 10

var nullBytes []byte = []byte("null")

// GetPaginateParams returns the correct offset and page according to the total number of rows.
func GetPaginateParams(page int64, count int64) (offset int64, newPage int64) {
	if count == 0 {
		return 0, 1
	}
	offset = (page - 1) * PageSize
	if offset < 0 {
		offset = 0
	}
	if offset >= count {
		offset = (count - 1) - ((count - 1) % PageSize)
	}
	newPage = offset/PageSize + 1
	return offset, newPage
}

type jsonError struct {
	Erreur string `json:"error"`
}

// NullTime is used for nullable time column
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Time, nt.Valid = value.(time.Time), true
	return nil
}

// Value implements the driver Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON implents the marshall interface
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid == false {
		return []byte("null"), nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON implents the unmarshal interface
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nt.Time)
	nt.Valid = (err == nil)
	return err
}

// NullBool uses alias in order to mashall and un marshall correctly
type NullBool sql.NullBool

// MarshalJSON implents the marshall interface
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid == false {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON implents the unmarshal interface
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		nb.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = (err == nil)
	return err
}

// Scan implements the Scanner interface
func (nb *NullBool) Scan(value interface{}) error {
	if value == nil {
		nb.Valid = false
		return nil
	}
	var n sql.NullBool
	if err := n.Scan(value); err != nil {
		return err
	}
	nb.Bool, nb.Valid = n.Bool, n.Valid
	return nil
}

// Value implements the driver Valuer interface
func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// NullInt64 uses alias in order to mashall and un marshall correctly
type NullInt64 sql.NullInt64

// MarshalJSON implents the marshall interface
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid == false {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int64)
}

// UnmarshalJSON implents the unmarshal interface
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// Scan implements the Scanner interface
func (ni *NullInt64) Scan(value interface{}) error {
	if value == nil {
		ni.Valid = false
		return nil
	}
	var n sql.NullInt64
	if err := n.Scan(value); err != nil {
		return err
	}
	ni.Int64, ni.Valid = n.Int64, n.Valid
	return nil
}

// Value implements the driver Valuer interface
func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

// NullString uses alias in order to mashall and un marshall correctly
type NullString sql.NullString

// MarshalJSON implents the marshall interface
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid == false {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}

// UnmarshalJSON implents the unmarshal interface
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		ns.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

// Scan implements the Scanner interface
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.Valid = false
		return nil
	}
	var n sql.NullString
	if err := n.Scan(value); err != nil {
		return err
	}
	ns.String, ns.Valid = n.String, n.Valid
	return nil
}

// Value implements the driver Valuer interface
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// TrimSpace removes leading and trailing spaces if value is not null
func (ns *NullString) TrimSpace() NullString {
	if !ns.Valid {
		return *ns
	}
	return NullString{Valid: true, String: strings.TrimSpace(ns.String)}
}

// NullFloat64 uses alias in order to mashall and un marshall correctly
type NullFloat64 sql.NullFloat64

// MarshalJSON implents the marshall interface
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid == false {
		return []byte("null"), nil
	}

	return json.Marshal(nf.Float64)
}

// UnmarshalJSON implents the unmarshal interface
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		nf.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

// Scan implements the Scanner interface
func (nf *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		nf.Valid = false
		return nil
	}
	var n sql.NullFloat64
	if err := n.Scan(value); err != nil {
		return err
	}
	nf.Float64, nf.Valid = n.Float64, n.Valid
	return nil
}

// Value implements the driver Valuer interface
func (nf NullFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}
