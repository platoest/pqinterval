package pqinterval

import (
	"testing"
)

func TestScanNullInterval(t *testing.T) {
	nival := new(NullInterval)
	err := nival.Scan("3 years 182 days 01:22:33.456789")
	if err != nil {
		t.Fatal(err)
	}

	ival := New(3, 182, 1, 22, 33, 456789)

	if !nival.Valid {
		t.Fatalf("valid interval: got %v, want %v", nival.Valid, true)
	}
	if got, want := nival.Interval, ival; got != want {
		t.Errorf("time value mismatch: got %v, want %v", got, want)
	}
}

func TestScanNilNullInterval(t *testing.T) {
	var nival NullInterval
	err := nival.Scan(nil)
	if err != nil {
		t.Error(err)
	}
	if nival.Valid {
		t.Errorf("invalid interval: got %v, want %v", nival.Valid, false)
	}
}
