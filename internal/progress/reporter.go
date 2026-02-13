// Package progress provides interfaces and implementations for progress reporting
// during long-running operations like network scanning.
package progress

// Reporter defines the interface for progress reporting during scans.
// Implementations can provide visual feedback (spinners, progress bars)
// or silent operation for non-interactive contexts.
type Reporter interface {
	// Start begins progress reporting with the given message
	Start(message string)

	// Update signals progress is being made (e.g., animate spinner)
	Update()

	// Finish stops progress reporting and cleans up
	Finish()
}

// NoOp is a Reporter that does nothing.
// Useful for testing or when progress reporting is disabled.
type NoOp struct{}

// Start is a no-op implementation.
func (NoOp) Start(string) {}

// Update is a no-op implementation.
func (NoOp) Update() {}

// Finish is a no-op implementation.
func (NoOp) Finish() {}
