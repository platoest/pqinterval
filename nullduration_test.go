package pqinterval

import (
	"testing"
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
