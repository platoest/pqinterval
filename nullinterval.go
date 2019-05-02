package pqinterval

import (
	"database/sql/driver"
)

type NullInterval struct {
	Interval
	Valid bool // Valid is true if Interval is not NULL
}

func (nival *NullInterval) Scan(src interface{}) error {
	if src == nil {
		nival.Valid = false
		return nil
	}

	var (
		ival Interval
		err  error
	)
	if err = (&ival).Scan(src); err != nil {
		nival.Valid = false
		return err
	}

	nival.Interval, nival.Valid = ival, true
	return nil
}

func (nival NullInterval) Value() (driver.Value, error) {
	if !nival.Valid {
		return nil, nil
	}
	return nival.Interval, nil
}

var _ driver.Valuer = NullInterval{}
