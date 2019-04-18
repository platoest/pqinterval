package pqinterval

import (
	"testing"
)

func TestScanNullInterval(t *testing.T) {
	var nival NullInterval
	ival := New(1, 2, 3, 4, 5, 6)
	err := nival.Scan(ival)
	if err != nil {
		t.Fatal(err)
	}
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
