package pqinterval

import (
	"testing"
)

func TestScanNullDuration(t *testing.T) {
	nd := new(NullDuration)
	ival := New(1, 2, 3, 4, 5, 6)
	d, _ := ival.Duration()
	_ = nd.Scan(ival)

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
