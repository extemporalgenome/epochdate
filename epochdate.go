// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package epochdate provides date handling in a compact format.
//
// Represents dates from Jan 1 1970 - Jun 6 2149 as days since the Unix epoch.
// This format requires 2 bytes (it's a uint16), in contrast to the 16 or 20
// byte representations (on 32 or 64-bit systems, respectively) used by the
// standard time package.
//
// Timezone is accounted for when applicable; when converting from standard
// time format into a Date, the date relative to the time value's zone is
// retained.  Times at any point during a given day (relative to timezone) are
// normalized to the same date.
//
// Conversely, conversions back to standard time format may be done using the
// Local, UTC, and In methods (semantically corresponding to the same-named
// Time methods), but with the result normalized to midnight (the beginning of
// the day) relative to that timezone.
//
// All functions and methods with the same names as those found in the stdlib
// time package have identical semantics in epochdate, with the exception that
// epochdate truncates time-of-day information.
package epochdate

import (
	"bytes"
	"errors"
	"time"
)

const (
	day      = 60 * 60 * 24
	nsPerSec = 1e9
	maxUnix  = (1<<16)*day - 1
)

// Format constants, for use with Parse and the Date.Format method.
const (
	RFC3339        = "2006-01-02"
	AmericanShort  = "1-2-06"
	AmericanCommon = "01-02-06"
)

// ErrOutOfRange is returned if the input date is not a representable Date.
var ErrOutOfRange = errors.New("epochdate: dates must be in the range [1970-01-01,2149-06-06]")

// Date stores the number of days since Jan 1, 1970. The last representable
// date is June 6, 2149.
type Date uint16

// Today returns the local date at this instant. If the local date does not
// fall within the representable range, then then zero value will be returned
// (1970-01-01).
func Today() Date {
	date, err := NewFromTime(time.Now())
	if err != nil {
		return 0
	}
	return date
}

// TodayUTC returns the date at this instant, relative to UTC. If the UTC
// date does not fall within the representable range, then then zero value
// will be returned (1970-01-01).
func TodayUTC() Date {
	date, _ := NewFromTime(time.Now().UTC())
	return date
}

// Parse follows the same semantics as time.Parse, but ignores time-of-day
// information and returns a Date value.
func Parse(layout, value string) (d Date, err error) {
	t, err := time.Parse(layout, value)
	if err == nil {
		d, err = NewFromTime(t)
	}
	return
}

// NewFromTime returns a Date equivalent to NewFromDate(t.Date()),
// where t is a time.Time object.
func NewFromTime(t time.Time) (Date, error) {
	s := t.Unix()
	_, offset := t.Zone()
	return NewFromUnix(s + int64(offset))
}

// NewFromDate returns a Date value corresponding to the supplied
// year, month, and day.
func NewFromDate(year int, month time.Month, day int) (Date, error) {
	return NewFromUnix(time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix())
}

// NewFromUnix creates a Date from a Unix timestamp, relative to any location
// Specifically, if you pass in t.Unix(), where t is a time.Time value with a
// non-UTC zone, you may receive an unexpected Date. Unless this behavior is
// specifically desired (returning the date in one location at the given time
// instant in another location), it's best to use epochdate.NewFromTime(t),
// which normalizes the resulting Date value by adjusting for zone offsets.
func NewFromUnix(seconds int64) (d Date, err error) {
	if UnixInRange(seconds) {
		d = Date(seconds / day)
	} else {
		err = ErrOutOfRange
	}
	return
}

// UnixInRange is true if the provided Unix timestamp is in Date's
// representable range. The timestamp is interpreted according to the semantics
// used by NewFromUnix. You probably won't need to use this, since this will
// only return false if NewFromUnix returns an error of ErrOutOfRange.
func UnixInRange(seconds int64) bool {
	return seconds >= 0 && seconds <= maxUnix
}

// Returns an RFC3339/ISO-8601 date string, of the form "2006-01-02".
func (d Date) String() string {
	return d.Format(RFC3339)
}

// Unix returns the number of seconds elapsed since Jan 1 1970 UTC, from the
// start of the given date value. In this case, the date is considered to be
// a UTC date, rather than a location-independent date.
func (d Date) Unix() int64 {
	return int64(d) * day
}

// UnixNano is semantically identical to the Unix method, except that it
// returns elapsed nanoseconds.
func (d Date) UnixNano() int64 {
	return int64(d) * day * nsPerSec
}

// Format is identical to time.Time.Format, except that any time-of-day format
// specifiers that are used will be equivalent to "00:00:00Z".
func (d Date) Format(layout string) string {
	return d.UTC().Format(layout)
}

// Date is semantically identical to the behavior of t.Date(), where t is a
// time.Time value.
func (d Date) Date() (year int, month time.Month, day int) {
	return d.UTC().Date()
}

// UTC returns a UTC Time object set to 00:00:00 on the given date.
func (d Date) UTC() time.Time {
	return time.Unix(int64(d)*day, 0).UTC()
}

// Local returns a local Time object set to 00:00:00 on the given date.
func (d Date) Local() time.Time {
	return d.In(time.Local)
}

// In returns a location-relative Time object set to 00:00:00 on the given date.
func (d Date) In(loc *time.Location) time.Time {
	t := time.Unix(int64(d)*day, 0).In(loc)
	_, offset := t.Zone()
	return t.Add(time.Duration(-offset) * time.Second)
}

// MarshalText implements encoding.TextMarshaler.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.Format(RFC3339)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = Parse(RFC3339, string(data))
	return err
}

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format(`"` + RFC3339 + `"`)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		return nil
	}
	data = bytes.Trim(data, `"`)
	return d.UnmarshalText(data)
}

var jsonNull = []byte(`null`)
