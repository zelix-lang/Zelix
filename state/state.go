/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
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

type SpinnerWrapper struct {
	start   time.Time
	spinner *yacspin.Spinner
	message *string
}

var spinners = map[int]SpinnerWrapper{}

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

func WarnAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightYellow, "⚠"))
	}
}

func PassAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightGreen, "✔"))
	}
}

func FailAllSpinners() {
	for _, spinner := range spinners {
		terminateSpinner(&spinner, ansi.Colorize(ansi.BoldBrightRed, "✖"))
	}
}

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
