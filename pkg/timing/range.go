package timing

import "time"

const (
	Desc Direction = iota - 1
	Equal
	Asc
)

type Range struct {
	Start time.Time
	End   time.Time
}

// UnBounded returns true when any of the two dates are zero values.
func (r *Range) UnBounded() bool {
	return r.Start.IsZero() || r.End.IsZero()
}

// Direction returns whether the range is Asc, Desc or Equal.
func (r *Range) Direction() Direction {
	if r.Start.Equal(r.End) {
		return Equal
	}

	if r.Start.Before(r.End) {
		return Asc
	} else {
		return Desc
	}
}

// Contains returns true when the given time is within the range
// inclusive of the start and end.
func (r *Range) Contains(time time.Time) bool {
	return r.ContainsExclusive(time) || time.Equal(r.End) || time.Equal(r.Start)
}

// ContainsExclusive returns true when the given time is within the range
// will not return true if the given time matches the start or end.
func (r *Range) ContainsExclusive(time time.Time) bool {
	return time.After(r.Start) && time.Before(r.End)
}

// Duration calculates the duration between Start and End.
// This can be negative if End is before Start.
func (r *Range) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// DurationMagnitude calculates the duration between Start and End.
// This is always positive or zero.
func (r *Range) DurationMagnitude() time.Duration {
	switch r.Direction() {
	case Asc:
		return r.End.Sub(r.Start)
	case Desc:
		return r.Start.Sub(r.End)
	default:
		return 0
	}
}

type Direction int
