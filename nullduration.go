package pqinterval

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
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

// MarshalJSON implements json.Marshaler
func (nd NullDuration) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(time.Duration(nd.Duration).String())
}

// UnmarshalJSON implements json.Unmarshaler
func (nd *NullDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		nd.Valid = true
		nd.Duration = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			nd.Valid = false
			return err
		}

		nd.Valid = true
		nd.Duration = Duration(tmp)
		return nil
	case nil:
		nd.Valid = false
		return nil
	default:
		return errors.New("invalid duration")
	}
}
