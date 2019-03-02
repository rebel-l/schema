// Package bar provides a fake progress bar to be compatible with the Progressor interface of package schema
package bar

// BlackHole represents a progress bar with no output
type BlackHole struct{}

// Increment does nothing
func (b *BlackHole) Increment() int { return 0 }

// FinishPrint does nothing
func (b *BlackHole) FinishPrint(msg string) {}
