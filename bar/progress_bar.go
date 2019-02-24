package bar

// BlackHole represents a progress bar with no output
type BlackHole struct{}

// Increment does nothing
func (b *BlackHole) Increment() int { return 0 }

// FinishPrint does nothing
func (b *BlackHole) FinishPrint(msg string) {}
