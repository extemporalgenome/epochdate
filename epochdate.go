// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Represents dates from Jan 1 1970 - Jun 6 2149 as days since the Unix epoch.
This format requires 2 bytes (it's a uint16), in contrast to the 16 or 20 byte
representations (on 32 or 64-bit systems, respectively) used by the standard
time package.

Timezone is accounted for when applicable; when converting from standard time
format into a Date, the date relative to the time value's zone is retained.
Times at any point during a given day (relative to timezone) are normalized to
the same date.

Conversely, conversions back to standard time format may be done using the
Local, UTC, and In methods (semantically corresponding to the same-named Time
methods), but with the result normalized to midnight (the beginning of the day)
relative to that timezone.
*/
package epochdate

import (
	"errors"
	"time"
)

const (
	day     = 60 * 60 * 24
	maxUnix = (1<<16)*day - 1
)

const (
	RFC3339       = "2006-01-02"
	AmericanShort = "1-2-06"
)

var ErrOutOfRange = errors.New("The given date is out of range")

type Date uint16

func Today() Date {
	date, err := NewFromTime(time.Now())
	if err != nil {
		panic(err)
	}
	return date
}

func Parse(layout, value string) (d Date, err error) {
	t, err := time.Parse(layout, value)
	if err == nil {
		d, err = NewFromTime(t)
	}
	return
}

func NewFromTime(t time.Time) (Date, error) {
	s := t.Unix()
	_, offset := t.Zone()
	return NewFromUnix(s + int64(offset))
}

func NewFromDate(year int, month time.Month, day int) (Date, error) {
	return NewFromUnix(time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix())
}

func NewFromUnix(seconds int64) (d Date, err error) {
	if UnixInRange(seconds) {
		d = Date(seconds / day)
	} else {
		err = ErrOutOfRange
	}
	return
}

func UnixInRange(seconds int64) bool {
	return seconds >= 0 && seconds <= maxUnix
}

func (d Date) String() string {
	return d.Format(RFC3339)
}

func (d Date) Format(layout string) string {
	return d.UTC().Format(layout)
}

func (d Date) Date() (year int, month time.Month, day int) {
	return d.UTC().Date()
}

func (d Date) UTC() time.Time {
	return time.Unix(int64(d)*day, 0).UTC()
}

func (d Date) Local() time.Time {
	return d.In(time.Local)
}

func (d Date) In(loc *time.Location) time.Time {
	t := time.Unix(int64(d)*day, 0).In(loc)
	_, offset := t.Zone()
	return t.Add(time.Duration(-offset) * time.Second)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format(`"` + RFC3339 + `"`)), nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	*d, err = Parse(`"`+RFC3339+`"`, string(data))
	return
}
