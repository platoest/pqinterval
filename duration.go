package pqinterval

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"
	"time"
)

// Duration is a time.Duration alias that supports the following additional interfaces:
// - driver.Valuer
// - sql.Scanner
// - json.Marshaler
// - json.Unmarshaler
type Duration time.Duration

// ErrTooBig is returned by Interval.Duration and Duration.Scan if the
// interval would overflow a time.Duration.
var ErrTooBig = errors.New("interval overflows time.Duration")

// Duration converts an Interval into a time.Duration with the same
// semantics as `EXTRACT(EPOCH from <interval>)` in PostgreSQL.
func (ival Interval) Duration() (time.Duration, error) {
	dur := int64(ival.Years())

	if dur > math.MaxInt64/nsPerYr || dur < math.MinInt64/nsPerYr {
		return 0, ErrTooBig
	}
	dur *= hrsPerYr
	dur += int64(ival.hrs)

	if dur > math.MaxInt64/int64(time.Hour) || dur < math.MinInt64/int64(time.Hour) {
		return 0, ErrTooBig
	}
	dur *= int64(time.Hour)

	us := ival.Microseconds() * int64(time.Microsecond)
	if dur > 0 {
		if math.MaxInt64-dur < us {
			return 0, ErrTooBig
		}
	} else {
		if math.MinInt64-dur > us {
			return 0, ErrTooBig
		}
	}
	dur += us

	return time.Duration(dur), nil
}

// Scan implements sql.Scanner.
func (d *Duration) Scan(src interface{}) error {
	ival := Interval{}
	err := (&ival).Scan(src)
	if err != nil {
		return err
	}

	result, err := ival.Duration()
	if err != nil {
		return err
	}

	*d = Duration(result)
	return nil
}

// Value implements driver.Valuer.
func (d Duration) Value() (driver.Value, error) {
	var years, months, days, hours, minutes, seconds, milliseconds, microseconds, nanoseconds int64
	nanoseconds = int64(d / Duration(time.Nanosecond))
	years, nanoseconds = divmod(nanoseconds, int64(time.Hour*hrsPerYr))
	days, nanoseconds = divmod(nanoseconds, int64(time.Hour*24))
	hours, nanoseconds = divmod(nanoseconds, int64(time.Hour))
	minutes, nanoseconds = divmod(nanoseconds, int64(time.Minute))
	seconds, nanoseconds = divmod(nanoseconds, int64(time.Second))
	milliseconds, nanoseconds = divmod(nanoseconds, int64(time.Millisecond))
	microseconds, _ = divmod(nanoseconds, int64(time.Microsecond))
	return formatInput(years, months, days, hours, minutes, seconds, milliseconds, microseconds), nil
}

func (d Duration) Milliseconds() float64 {
	td := time.Duration(d)
	msec := td / time.Millisecond
	nsec := td % time.Millisecond
	return float64(msec) + float64(nsec)/1e6
}

// MarshalJSON implements json.Marshaler. IMPORTANT NOTE: We are serializing
// Duration as an *integer* number of milliseconds. For our use case, this is
// sufficient: we don't need any more precision than millisecond.
func (d Duration) MarshalJSON() ([]byte, error) {
	ms := int64(math.Round(d.Milliseconds()))
	return json.Marshal(ms)
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value * 1e6))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
