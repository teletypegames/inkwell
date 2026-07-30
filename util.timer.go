package inkwell

// Timer accumulates seconds. Tick advances it; Reset zeroes it.
type Timer struct {
	Elapsed float64
}

func (t *Timer) Tick(dt float64) { t.Elapsed += dt }
func (t *Timer) Reset()           { t.Elapsed = 0 }
func (t *Timer) Done(after float64) bool { return t.Elapsed >= after }
