/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package cli

import (
	"fluent/logger"
	"fmt"
	"github.com/urfave/cli/v3"
)

// BenchmarkCommand benchmarks the current implementation of the Fluent programming language.
func BenchmarkCommand(context *cli.Command) {
	ShowHeaderMessage()

	// Get the times to run the benchmark
	times := context.Int("times")

	fmt.Println(times)
	// Ensure that the times are at least 1
	if times < 1 {
		logger.Warn("The times to run the benchmark must be at least 1")
		logger.Info("Defaulting to 1")
		times = 1
	}

	// TODO!
}
