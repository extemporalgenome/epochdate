# epochdate

epochdate is a small library for storing contemporary dates (without
time-of-day information) in 2 bytes per date, instead of the 16-20 bytes that
are used by a time.Time value from the Go standard library.

epochdate.Date is a uint16 value representing the number of days since
Jan 1 1970. The maximum representable date is Jun 6 2149. Arithmetical
operations react predictably. 

	import (
		"github.com/extemporalgenome/epochdate"
		"time"
	)
	
	var (
		today     = epochdate.Today()
		yesterday = today - 1
		tomorrow  = today + 1
		fortnight = today + 7*2
	)

## Converting to Date values

	assert := func(condition bool) { if !condition { panic("assertion failed") } }

	d1, _ := epochdate.Parse(epochdate.RFC3339, "2012-03-10")
	d2, _ := epochdate.NewFromDate(2012, 3, 10)
	assert(d1 == d2)

	times := []string{
		// Any time of the same day forms an equivalent epochdate value
		"2012-03-10T00:00:00Z",
		"2012-03-10T23:59:59Z",

		// And this is relative to location -- even locations that are
		// 26 hours apart (politics) will result in the same epochdate value
		// (in each location, the date at the given instant *is* 2012-03-10)
		"2012-03-10T00:00:00-12:00",
		"2012-03-10T00:00:00+14:00",
	}

	for _, str := range times {
		t, _ := time.Parse(time.RFC3339, str)
		d, _ := epochdate.NewFromTime(t)
		assert(d == d1)
	}

	t1, _ := time.Parse(time.RFC3339, "2012-03-10T00:00:00-12:00")
	t2, _ := time.Parse(time.RFC3339, "2012-03-10T00:00:00+14:00")

	d1, _ := epochdate.NewFromUnix(t1.Unix())
	d2, _ := epochdate.NewFromUnix(t2.Unix())

	// t1.Unix() returns seconds since the Unix epoch relative to UTC!
	// Because t1 and t2 represent very different time instants
	// (26 hours apart), don't expect passing those timestamps to NewFromUnix
	// to result in the same date; unless you know what you're doing, use the
	// other "New" functions instead, which normalize timezones.
	assert(d1 != d2)

## Converting from Date values

	var (
		d, _                = epochdate.Parse(epochdate.RFC3339, "2012-03-10")
		t, _                = time.Parse(time.RFC3339, "2012-03-10T00:00:00Z")
		year1, month1, day1 = d.Date()
		year2, month2, day2 = t.Date()
	)
	assert(year1 == year2 && month1 == month2 && day1 == day2)
	assert(t.Equal(d.UTC()) == true)

	// epochdate's Local and In methods return times adjusted to midnight on
	// that date in the given timezone.
	assert(t.Equal(d.Local()) == false)

	// So if you want the local time at the instant that it was midnight on
	// that date relative to UTC, do d.UTC().Local(). Conversely, the UTC time
	// at the instant that it was midnight on the given date in the local
	// timezone is expressed as: d.Local().UTC()

	// Date values print nicely
	assert(d.String() == "2012-03-10")
