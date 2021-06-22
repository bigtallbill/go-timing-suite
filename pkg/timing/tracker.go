package timing

import (
	"fmt"
	"sync"
	"time"
)

// CappedOpGroup allows the tracking of a linear series Operations over time.
// It is safe to add Ops to the group in goroutines.
type CappedOpGroup struct {
	Name  string
	Ops   Ops
	Cap   int
	Done  int
	Left  int
	Total int
	l     sync.RWMutex
}

func NewCappedOpGroup(name string, cap int, totalOps int) *CappedOpGroup {
	return &CappedOpGroup{
		Name:  name,
		Cap:   cap,
		Total: totalOps,
		Left:  totalOps,
		Ops:   Ops{},
	}
}

// AddOp adds a new Op to the group.
// it is impossible to add more Ops than CappedOpGroup.Total.
// Additionally, the Op must not be running (has valid start and end).
func (t *CappedOpGroup) AddOp(op *Op) error {
	t.l.Lock()
	defer t.l.Unlock()

	if op.Running() {
		return fmt.Errorf("cannot add running Op: %+v", *op)
	}

	if t.Left == 0 {
		return fmt.Errorf("cannot add another, all Ops complete")
	}

	t.Ops = append(t.Ops, op)
	if len(t.Ops) > t.Cap {
		t.Ops = t.Ops[1:]
	}

	t.Done++
	t.Left--

	return nil
}

// Eta returns the estimated time of arrival.
func (t *CappedOpGroup) Eta() time.Time {
	return time.Now().Add(t.TimeLeft())
}

// TimeLeft returns the time left for remaining operations to complete based on
// the average completion time of the current set of Ops.
func (t *CappedOpGroup) TimeLeft() time.Duration {
	t.l.RLock()
	defer t.l.RUnlock()

	var avg = t.Ops.Avg()

	return avg * time.Duration(t.Left)
}

// PercentComplete returns the percent of complete Ops.
func (t *CappedOpGroup) PercentComplete() float64 {
	t.l.RLock()
	defer t.l.RUnlock()

	return float64(t.Left) / float64(t.Total)
}
