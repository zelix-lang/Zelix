/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package state

import (
	"fluent/ansi"
	"github.com/theckman/yacspin"
	"time"
)

type State int

const (
	Lexing State = iota
	Parsing
	Processing
	Analyzing
	Compiling
)

// SpinnerWrapper is a struct that holds information about a spinner,
// including its start time, the spinner instance, and a message.
type SpinnerWrapper struct {
	start   time.Time        // The time when the spinner was started
	spinner *yacspin.Spinner // The spinner instance
	message *string          // The message associated with the spinner
}

// A map of spinners, where the key is the spinner ID and the value is the SpinnerWrapper.
var spinners = map[int]SpinnerWrapper{}

// terminateSpinner stops the spinner, updates its message with the elapsed time,
// and removes it from the spinners map.
// Parameters:
// - spinner: A pointer to the SpinnerWrapper to be terminated.
// - character: A string character to prepend to the stop message.
func terminateSpinner(spinner *SpinnerWrapper, character string) {
	// Calculate the time taken
	elapsed := time.Since(spinner.start)

	// Modify the message
	spinner.spinner.StopMessage(
		character + *spinner.message +
			ansi.Colorize(ansi.BrightBlack, " ("+elapsed.String()+")"),
	)

	// Stop the spinner
	err := spinner.spinner.Stop()

	if err != nil {
		panic(err)
	}

	// Remove the spinner
	delete(spinners, 0)
}

// WarnAllSpinners terminates all spinners with a warning message.
func WarnAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightYellow, "⚠"))
	}
}

// PassAllSpinners terminates all spinners with a success message.
func PassAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightGreen, "✔"))
	}
}

// FailAllSpinners terminates all spinners with a failure message.
func FailAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightRed, "✖"))
	}
}

// Emit creates and starts a new spinner based on the given state and text.
// It returns the ID of the created spinner or -1 if an error occurs.
//
// Parameters:
// - state: The current state of the process (e.g., Lexing, Parsing).
// - text: The text message to be displayed with the spinner.
//
// Returns:
// - int: The ID of the created spinner or -1 if an error occurs.
func Emit(state State, text string) int {
	// Add the message together accordingly
	spinnerText := " "
	var color string

	switch state {
	case Lexing:
		color = ansi.BoldBrightPurple
		spinnerText += "Lexing"
	case Parsing:
		color = ansi.BoldBrightBlue
		spinnerText += "Parsing"
	case Processing:
		color = ansi.BoldBrightYellow
		spinnerText += "Processing"
	case Analyzing:
		color = ansi.BoldBrightGreen
		spinnerText += "Analyzing"
	case Compiling:
		color = ansi.BoldBrightGreen
		spinnerText += "Compiling"
	}

	spinnerText += " " + text
	spinnerText = ansi.Colorize(color, spinnerText)

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[69],
		Suffix:          "",
		SuffixAutoColon: true,
		Message:         spinnerText,
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)

	// Handle the error
	if err != nil {
		return -1
	}

	err = spinner.Start()

	if err != nil {
		return -1
	}

	// Create a new wrapper and push it
	wrapper := SpinnerWrapper{
		start:   time.Now(),
		spinner: spinner,
		message: &spinnerText,
	}

	id := len(spinners)
	spinners[id] = wrapper

	return id
}
