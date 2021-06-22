package timing

import "time"

//==============================================================================
// Singular
//==============================================================================

type Op struct {
	Name    string
	Runtime Range
}

// Start fills the Op's start time with the current time.
func (o *Op) Start() {
	o.Runtime.Start = time.Now()
}

// End fills the Op's end time with the current time.
func (o *Op) End() {
	o.Runtime.End = time.Now()
}

// Running is true if either start or end are time.Time zero value.
func (o *Op) Running() bool {
	return o.Runtime.UnBounded()
}

func NewOp(name string) *Op {
	return &Op{Name: name}
}

//==============================================================================
// Group
//==============================================================================

type Ops []*Op

// Avg calculates the avg duration of all ops.
func (ops Ops) Avg() time.Duration {
	var sum time.Duration
	for _, op := range ops {
		sum += op.Runtime.Duration()
	}

	return sum / time.Duration(len(ops))
}

// Longest finds the Op with the longest duration.
func (ops Ops) Longest() *Op {
	var longest *Op
	for _, op := range ops {
		if longest == nil {
			longest = op
			continue
		}

		if op.Runtime.Duration() > longest.Runtime.Duration() {
			longest = op
		}
	}

	return longest
}

// Shortest finds the Op with the shortest duration.
func (ops Ops) Shortest() *Op {
	var shortest *Op
	for _, op := range ops {
		if shortest == nil {
			shortest = op
			continue
		}

		if op.Runtime.Duration() < shortest.Runtime.Duration() {
			shortest = op
		}
	}

	return shortest
}

// StartedFirst returns the Op with the earliest start time.
func (ops Ops) StartedFirst() *Op {
	var first *Op
	for _, op := range ops {
		if first == nil {
			first = op
			continue
		}

		if op.Runtime.Start.Before(first.Runtime.Start) {
			first = op
		}
	}

	return first
}

// EndedLast returns the Op with the latest end time.
func (ops Ops) EndedLast() *Op {
	var last *Op
	for _, op := range ops {
		if last == nil {
			last = op
			continue
		}

		if op.Runtime.End.Before(last.Runtime.End) {
			last = op
		}
	}

	return last
}

// FullRange calculates the Range of the earliest start and latest end time.
func (ops Ops) FullRange() Range {
	return Range{
		Start: ops.StartedFirst().Runtime.Start,
		End:   ops.EndedLast().Runtime.End,
	}
}
