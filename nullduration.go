package pqinterval

import (
	"database/sql/driver"
)

type NullDuration struct {
	Duration
	Valid bool
}

func (nd *NullDuration) Scan(src interface{}) error {
	if src == nil {
		nd.Duration, nd.Valid = Duration(0), false
		return nil
	}

	nival := NullInterval{}
	err := (&nival).Scan(src)
	if err != nil {
		return err
	}

	result, err := nival.Duration()
	if err != nil {
		return err
	}

	nd.Duration, nd.Valid = Duration(result), true
	return nil
}

func (nd NullDuration) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}

	return nd.Duration, nil
}

var _ driver.Valuer = NullDuration{}
