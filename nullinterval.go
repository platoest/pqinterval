package pqinterval

import (
	"database/sql/driver"
)

type NullInterval struct {
	Interval
	Valid bool // Valid is true if Interval is not NULL
}

func (nival *NullInterval) Scan(src interface{}) error {
	nival.Interval, nival.Valid = src.(Interval)
	return nil
}

func (nival NullInterval) Value() (driver.Value, error) {
	if !nival.Valid {
		return nil, nil
	}
	return nival.Interval, nil
}

var _ driver.Valuer = NullInterval{}
