// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epochdate

import (
	"encoding"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var (
	_ encoding.TextMarshaler   = Date(0)
	_ encoding.TextUnmarshaler = new(Date)
	_ json.Marshaler           = Date(0)
	_ json.Unmarshaler         = new(Date)
)

type triple struct {
	year  int
	month time.Month
	day   int
}

type equiv struct {
	date Date
	unix int64
	trip triple
	str  string
}

var equivs = []equiv{
	{0, 0, triple{1970, 1, 1}, "1970-01-01"},
	{0, day - 1, triple{1970, 1, 1}, "1970-01-01"},
	{366, 366 * day, triple{1971, 1, 2}, "1971-01-02"},
	{366, 367*day - 1, triple{1971, 1, 2}, "1971-01-02"},
	{65535, 65535 * day, triple{2149, 6, 6}, "2149-06-06"},
	{65535, 65536*day - 1, triple{2149, 6, 6}, "2149-06-06"},
}

var extrema = []struct {
	unix  int64
	valid bool
}{
	{-1, false},
	{0, true},
	{65536*day - 1, true},
	{65536 * day, false},
}

func TestEquivalences(t *testing.T) {
	for _, e := range equivs {
		if unix, err := NewFromUnix(e.unix); err != nil {
			t.Fatal(err)
		} else if trip, err := NewFromDate(e.trip.year, e.trip.month, e.trip.day); err != nil {
			t.Fatal(err)
		} else if str, err := Parse(RFC3339, e.str); err != nil {
			t.Fatal(err)
		} else if e.date != unix || e.date != trip || e.date != str {
			t.Fatal("Unexpected non-equivalence:", e.date, unix, trip, str)
		}
	}
}

func TestExtrema(t *testing.T) {
	var desc string
	for _, e := range extrema {
		if UnixInRange(e.unix) != e.valid {
			if e.valid {
				desc = "valid"
			} else {
				desc = "invalid"
			}
			t.Fatal("Unix timestamp", e.unix, "should be", desc)
		}
	}
}

func TestTimezoneIrrelevance(t *testing.T) {
	const hour = 60 * 60
	min := time.FixedZone("min", -12*hour)
	max := time.FixedZone("max", +14*hour)
	t1 := time.Date(2149, 06, 06, 0, 0, 0, 0, min)
	t2 := time.Date(2149, 06, 06, 0, 0, 0, 0, max)
	var (
		d1, d2 Date
		err    error
	)
	if d1, err = NewFromTime(t1); err != nil {
		t.Fatal(err)
	}
	if d2, err = NewFromTime(t2); err != nil {
		t.Fatal(err)
	}
	if d1 != d2 {
		t.Fatal("Expected", t1, "and", t2, "to result in same date; got", d1, "and", d2)
	}
}

func TestDateToTime(t *testing.T) {
	var date Date
	local := date.Local()
	utc := date.UTC()
	prefix := "1970-01-01T00:00:00"
	if f := local.Format(time.RFC3339); !strings.HasPrefix(f, prefix) {
		t.Fatalf("Expected local time to %q; got %q\n", prefix, f)
	} else if f := utc.Format(time.RFC3339); !strings.HasPrefix(f, prefix) {
		t.Fatalf("Expected universal time to %q; got %q\n", prefix, f)
	}
}

func TestUnix(t *testing.T) {
	var d Date = 1
	const (
		dayInSecs     = 60 * 60 * 24
		dayInNanosecs = dayInSecs * 1e9
	)
	if s := d.Unix(); s != dayInSecs {
		t.Error("Expected Date(1).Unix() to return", dayInSecs, "but got", s)
	}
	if ns := d.UnixNano(); ns != dayInNanosecs {
		t.Error("Expected Date(1).UnixNano() to return", dayInNanosecs, "but got", ns)
	}
}

func TestEncDec(t *testing.T) {
	const (
		unquoted = "1970-01-02"
		quoted   = `"` + unquoted + `"`
		n        = 1
	)
	d := Date(n)
	b, err := d.MarshalText()
	if err != nil {
		t.Error("Unexpected MarshalText error:", err)
	} else if string(b) != unquoted {
		t.Errorf("Expected Date(%d).MarshalText() to return %#q but got %#q", n, unquoted, b)
	}
	b, err = d.MarshalJSON()
	if err != nil {
		t.Error("Unexpected MarshalJSON error:", err)
	} else if string(b) != quoted {
		t.Errorf("Expected Date(%d).MarshalJSON() to return %#q but got %#q", n, quoted, b)
	}
	err = d.UnmarshalText([]byte(unquoted))
	if err != nil {
		t.Error("Unexpected UnmarshalText error:", err)
	} else if d != n {
		t.Errorf("Expected Date(%d).UnmarshalText() to return %#q but got %#q", n, unquoted, b)
	}
	err = d.UnmarshalJSON([]byte(quoted))
	if err != nil {
		t.Error("Unexpected UnmarshalJSON error:", err)
	} else if d != n {
		t.Errorf("Expected Date(%d).UnmarshalJSON() to return %#q but got %#q", n, quoted, b)
	}
}

func TestDate_UnmarshalJSON_null(t *testing.T) {
	data := []byte("null")
	input := Date(123)
	date := input
	want := input

	err := json.Unmarshal(data, &date)
	if err != nil {
		t.Errorf("json.Unmarshal(%q) = %q, want nil", data, err)
	}

	if date != want {
		t.Errorf("json.Unmarshal(%q, %v) -> %v, want %v",
			data, input, date, want)
	}
}
