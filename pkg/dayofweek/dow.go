package dayofweek

import (
	"time"
)

// Bitmask is a bitmask type for day of week
type Bitmask uint8

// All flags for day of week
const (
	Sunday Bitmask = 1 << iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// set the flag on the bitmask
func set(b, flag Bitmask) Bitmask { return b | flag }

// clear the flag on the bitmask
func clear(b, flag Bitmask) Bitmask { return b &^ flag }

// flip the flag on the bitmask
func flip(b, flag Bitmask) Bitmask { return b ^ flag }

// check the flag in the bitmask
func check(b, flag Bitmask) bool { return b&flag != 0 }

// Enable a day of week in the bitmask
func Enable(day int, bitmask *Bitmask) {
	var d Bitmask
	switch time.Weekday(day) {
	case time.Sunday:
		d = Sunday
	case time.Monday:
		d = Monday
	case time.Tuesday:
		d = Tuesday
	case time.Wednesday:
		d = Wednesday
	case time.Thursday:
		d = Thursday
	case time.Friday:
		d = Friday
	case time.Saturday:
		d = Saturday
	}
	*bitmask = set(*bitmask, d)
}

// Enabled checks if the provided day of week
// is enabled in the provided bitmask
func Enabled(day time.Weekday, bitmask Bitmask) bool {
	var d Bitmask
	switch day {
	case time.Sunday:
		d = Sunday
	case time.Monday:
		d = Monday
	case time.Tuesday:
		d = Tuesday
	case time.Wednesday:
		d = Wednesday
	case time.Thursday:
		d = Thursday
	case time.Friday:
		d = Friday
	case time.Saturday:
		d = Saturday
	}
	return check(bitmask, d)
}

// New Bitmask pointer
func New() *Bitmask {
	return new(Bitmask)
}

// Enable a day of week in the bitmask
func (m *Bitmask) Enable(day time.Weekday) {
	Enable(int(day), m)
}

// EnableAll day of week in the bitmask
func (m *Bitmask) EnableAll() {
	*m = 0b01111111
}

// Enabled checks if the provided day of week
// is enabled in the bitmask
func (m Bitmask) Enabled(day time.Weekday) bool {
	return Enabled(day, m)
}

// AllEnabled returns true when all weekdays
// are enabled in the bitmask
func (m Bitmask) AllEnabled() bool {
	return m == 0b00000000 || m == 0b01111111
}
