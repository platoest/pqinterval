package pqinterval

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationValue(t *testing.T) {
	i := new(Duration)
	_ = i.Scan("3 years 182 days 01:22:33.456789")

	val, err := i.Value()
	assert.Nil(t, err, "Duration.Value() error")
	assert.EqualValues(
		t,
		"3 years 182 days 1 hours 22 minutes 33 seconds 456 milliseconds 789 microseconds",
		val,
		"Duration.Value() result")

	j := time.Duration(30) * time.Minute
	k := Duration(j)
	val, err = k.Value()
	assert.Nil(t, err, "Duration.Value() error")
	assert.EqualValues(
		t,
		"30 minutes",
		val,
		"Duration.Value() compatibility with time.Duration")
}

func TestZeroDuration(t *testing.T) {
	i := new(Duration)
	assert.EqualValues(t, time.Duration(0), *i, "Duration.Scan() result")

	val, err := i.Value()
	assert.Nil(t, err, "Duration.Value() error")
	assert.EqualValues(t, "0 microseconds", val, "Duration.Value() result")

	assert.NoError(t, i.Scan("00:00:00"), "Duration.Scan() error")
	assert.EqualValues(t, time.Duration(0), *i, "Duration.Scan() result")

	val, err = i.Value()
	assert.Nil(t, err, "Duration.Value() error")
	assert.EqualValues(t, "0 microseconds", val, "Duration.Value() result")
}

func TestDuration_MarshalJSON(t *testing.T) {
	orig := "20m30s"
	d, err := time.ParseDuration(orig)
	if err != nil {
		t.Fatal(err)
	}
	pqd := Duration(d)
	b, err := pqd.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if got, want := string(b), fmt.Sprintf(`"%v"`, orig); got != want {
		t.Errorf("bad marshal: got %v, want %v", got, want)
	}
}

func TestDuration_UnmarshalJSON_string(t *testing.T) {
	input := []byte(`"20m30s"`)
	var d Duration
	err := (&d).UnmarshalJSON(input)
	if err != nil {
		t.Error(err)
	}
	if got, want := d, Duration(1230000000000); got != want {
		t.Errorf("bad unmarshal: got %v, want %v", got, want)
	}
}
