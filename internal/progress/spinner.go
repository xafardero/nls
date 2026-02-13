package progress

import (
	"github.com/schollz/progressbar/v3"
)

// Spinner provides visual feedback during scanning using a progress spinner.
// It implements the Reporter interface using the progressbar library.
type Spinner struct {
	bar *progressbar.ProgressBar
}

// NewSpinner creates a new Spinner progress reporter.
// The spinner displays an animated indicator with the provided message.
func NewSpinner() *Spinner {
	return &Spinner{}
}

// Start begins displaying the progress spinner with the given message.
func (s *Spinner) Start(message string) {
	s.bar = progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(message),
		progressbar.OptionSpinnerType(14),
	)
}

// Update advances the spinner animation.
// Should be called periodically to animate the spinner.
func (s *Spinner) Update() {
	if s.bar != nil {
		_ = s.bar.Add(1)
	}
}

// Finish stops the spinner and clears the display.
func (s *Spinner) Finish() {
	if s.bar != nil {
		_ = s.bar.Finish()
	}
}
