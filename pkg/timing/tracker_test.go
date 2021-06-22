package timing

import (
	"testing"
	"time"
)

func TestCappedGroup(t *testing.T) {
	t.Run("eta calculation", func(t *testing.T) {
		now := time.Now()
		cg := NewCappedOpGroup("Group1", 5, 10)

		for i := 0; i < 5; i++ {
			_ = cg.AddOp(&Op{
				Runtime: Range{
					Start: now,
					End:   now.Add(time.Millisecond),
				},
			})
		}

		timeLeft := cg.TimeLeft()

		if actual, expected := timeLeft, time.Millisecond*5; expected != actual {
			t.Errorf("expected actual to be %d but got %d", actual, expected)
		}

		eta := cg.Eta()

		etaDiffToActual := eta.Sub(now.Add(timeLeft)).Milliseconds()
		if etaDiffToActual > 1 {
			t.Errorf("expected %s but got %s", now.Add(timeLeft), eta)
		}
	})

	t.Run("cap stays capped", func(t *testing.T) {
		now := time.Now()
		cg := NewCappedOpGroup("Group1", 5, 10)

		for i := 0; i < 10; i++ {
			_ = cg.AddOp(&Op{
				Runtime: Range{
					Start: now,
					End:   now.Add(time.Millisecond * time.Duration(i)),
				},
			})
		}

		if actual, expected := len(cg.Ops), 5; actual != expected {
			t.Errorf("expected %d Ops but found %d", expected, actual)
		}

		if actual, expected := cg.Ops.Longest().Runtime.Duration(), time.Millisecond*9; actual != expected {
			t.Errorf("expected longest to be %d but found %d", expected, actual)
		}
	})

	t.Run("cannot exceed total Ops", func(t *testing.T) {
		now := time.Now()
		cg := NewCappedOpGroup("Group1", 5, 10)

		for i := 0; i < 10; i++ {
			_ = cg.AddOp(&Op{
				Runtime: Range{
					Start: now,
					End:   now.Add(time.Millisecond * time.Duration(i)),
				},
			})
		}

		// try add one more
		err := cg.AddOp(&Op{
			Runtime: Range{
				Start: now,
				End:   now.Add(time.Millisecond * 11),
			},
		})

		if err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("Percent calculation", func(t *testing.T) {
		now := time.Now()
		cg := NewCappedOpGroup("Group1", 5, 20)

		for i := 0; i < 10; i++ {
			_ = cg.AddOp(&Op{
				Runtime: Range{
					Start: now,
					End:   now.Add(time.Millisecond * time.Duration(i)),
				},
			})
		}

		if actual, expected := cg.PercentComplete(), .5; actual != expected {
			t.Errorf("expected %f but found %f", expected, actual)
		}
	})

	t.Run("cant add running op", func(t *testing.T) {
		now := time.Now()
		cg := NewCappedOpGroup("Group1", 5, 20)

		err := cg.AddOp(&Op{
			Runtime: Range{
				Start: now,
			},
		})

		if err == nil {
			t.Error("expected an error")
		}
	})
}
