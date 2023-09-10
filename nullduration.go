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

func NewNullDuration(d Duration, valid bool) NullDuration {
	return NullDuration{
		Duration: d,
		Valid:    valid,
	}
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

	return nd.Duration.Value()
}

var _ driver.Valuer = NullDuration{}

// Milliseconds returns the number of milliseconds in the duration
func (nd NullDuration) Milliseconds() float64 {
	if !nd.Valid {
		return 0
	}
	td := time.Duration(nd.Duration)
	msec := td / time.Millisecond
	nsec := td % time.Millisecond
	return float64(msec) + float64(nsec)/1e6
}

// MarshalJSON implements json.Marshaler. IMPORTANT NOTE: We are serializing
// NullDuration as an *integer* number of milliseconds. For our use case, this
// is sufficient: we don't need any more precision than millisecond.
func (nd NullDuration) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nd.Duration)
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
		nd.Duration = Duration(time.Duration(value * 1e6))
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
