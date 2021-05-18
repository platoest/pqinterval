package pqinterval

import (
	"strconv"
	"testing"
	"time"
)

func TestScanNullDuration(t *testing.T) {
	s := "3 years 182 days 01:22:33.456789"
	ival := New(3, 182, 1, 22, 33, 456789)
	d, _ := ival.Duration()

	nd := new(NullDuration)
	err := nd.Scan(s)
	if err != nil {
		t.Fatal(err)
	}

	if !nd.Valid {
		t.Errorf("valid duration: got %v, want %v", nd.Valid, true)
	}
	if got, want := nd.Duration, Duration(d); got != want {
		t.Errorf("time value mismatch: got %v, want %v", got, want)
	}
}

func TestScanNullDurationLargeHours(t *testing.T) {
	s := "163:08:24.636"
	ival := New(0, 0, 163, 8, 24, 636000)
	d, _ := ival.Duration()

	nd := new(NullDuration)
	err := nd.Scan(s)
	if err != nil {
		t.Fatal(err)
	}

	if !nd.Valid {
		t.Errorf("valid duration: got %v, want %v", nd.Valid, true)
	}
	if got, want := nd.Duration, Duration(d); got != want {
		t.Errorf("time value mismatch: got %v, want %v", got, want)
	}
}

func TestScanNilNullDuration(t *testing.T) {
	var nd NullDuration
	err := nd.Scan(nil)
	if err != nil {
		t.Error(err)
	}

	if nd.Valid {
		t.Errorf("invalid duration valid: got %v, want %v", nd.Valid, false)
	}
}

func TestNullDuration_MarshalJSON(t *testing.T) {
	orig := "20m30s"
	d, err := time.ParseDuration(orig)
	if err != nil {
		t.Fatal(err)
	}
	pqd := NewNullDuration(Duration(d), true)
	b, err := pqd.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if got, want := string(b), strconv.Itoa(1230000); got != want {
		t.Errorf("bad marshal: got %v, want %v", got, want)
	}
}

func TestNullDuration_UnmarshalJSON_string(t *testing.T) {
	input := []byte(`"20m30s"`)
	var nd NullDuration
	err := (&nd).UnmarshalJSON(input)
	if err != nil {
		t.Error(err)
	}
	if got, want := nd.Duration, Duration(1230000000000); got != want {
		t.Errorf("bad unmarshal: got %v, want %v", got, want)
	}
	if got, want := nd.Valid, true; got != want {
		t.Errorf("invalid NullDuration: got %v, want %v", got, want)
	}
}

func TestNullDuration_UnmarshalJSON_millis(t *testing.T) {
	input := []byte(`1230000`)
	var nd NullDuration
	err := (&nd).UnmarshalJSON(input)
	if err != nil {
		t.Error(err)
	}
	if got, want := nd.Duration, Duration(1230000000000); got != want {
		t.Errorf("bad unmarshal: got %v, want %v", got, want)
	}
	if got, want := nd.Valid, true; got != want {
		t.Errorf("invalid NullDuration: got %v, want %v", got, want)
	}
}

func TestNullDuration_UnmarshalJSON_null(t *testing.T) {
	input := []byte(`null`)
	var nd NullDuration
	err := (&nd).UnmarshalJSON(input)
	if err != nil {
		t.Error(err)
	}
	if got, want := nd.Valid, false; got != want {
		t.Errorf("bad NullDuration validity: got %v, want %v", got, want)
	}
	if got, want := nd.Duration, Duration(0); got != want {
		t.Errorf("bad unmarshal: got %v, want %v", got, want)
	}
}
