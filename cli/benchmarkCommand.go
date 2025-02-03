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
	"fluent/analyzer"
	"fluent/ansi"
	"fluent/filecode/converter"
	"fluent/logger"
	"fluent/util"
	"fmt"
	"github.com/theckman/yacspin"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// BenchmarkCommand benchmarks the current implementation of the Fluent programming language.
func BenchmarkCommand(context *cli.Command) {
	ShowHeaderMessage()

	// Get the times to run the benchmark
	times := context.Int("times")

	// Ensure that the times are at least 1
	if times < 1 {
		logger.Warn("The times to run the benchmark must be at least 1")
		logger.Info("Defaulting to 1")
		times = 1
	}

	// Ensure the benchmark directory exists
	if !util.DirExists("benchmark") {
		err := os.Mkdir("benchmark", 0755)

		if err != nil {
			logger.Error("Could not create the benchmark directory")
			logger.Help("Ensure that you have the necessary permissions")
			return
		}
	}

	// Get the cwd
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Could not get the current working directory")
		logger.Help("Ensure that you have the necessary permissions")
		return
	}

	exampleDir := path.Join(cwd, "example")

	// Ensure the "example" directory has files
	files := util.ReadDir(exampleDir)
	if len(files) == 0 {
		logger.Error("The example directory does not exist or is empty")
		logger.Help("Ensure that the example directory exists")
		return
	}

	fmt.Print(
		fmt.Sprintf(
			"%s=> Running the benchmark%s %s%d%s%s times on the example directory%s\n\n",
			ansi.BoldBrightBlue,
			ansi.Reset,
			ansi.BoldBrightPurple,
			times,
			ansi.Reset,
			ansi.BoldBrightBlue,
			ansi.Reset,
		),
	)

	// Use a strings.builder to build the output
	var output strings.Builder

	timestamp := time.Now().Format("2006_January_02 15_04_05")

	// Append a header to the output
	output.WriteString("The Fluent Programming Language Benchmark\n")
	output.WriteString("Time: ")
	output.WriteString(timestamp)
	output.WriteString("\n")
	output.WriteString("========================================\n")
	output.WriteString("Running on: ")
	output.WriteString(runtime.GOOS)
	output.WriteString(" ")
	output.WriteString(runtime.GOARCH)
	output.WriteString("\n")
	output.WriteString("========================================\n")
	output.WriteString("Times: ")
	output.WriteString(fmt.Sprintf("%d", times))
	output.WriteString("\n")
	output.WriteString("========================================\n\n")

	// Build the file name in the format: year-month_name-day-hour-minute-second-bench.txt
	fileName := fmt.Sprintf(
		"benchmark/%s-bench.txt",
		timestamp,
	)

	// Create a new spinner
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[69],
		Suffix:          "",
		SuffixAutoColon: true,
		Message:         "Benchmarking",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)

	if err != nil {
		logger.Error("Could not create a spinner")
		logger.Help("Ensure that you have the necessary permissions")
		return
	}

	err = spinner.Start()

	if err != nil {
		logger.Error("Could not start the spinner")
		logger.Help("This is probably an issue with your terminal")
		return
	}

	// Run the benchmark
	for _, file := range files {
		name := file.Name()

		fileHeader := fmt.Sprintf(
			"== %s ==",
			name,
		)

		// Append the file header to the output
		output.WriteString(fileHeader)
		output.WriteString("\n")

		// Calculate the absolute path of the file
		filePath := path.Join(exampleDir, name)

		var i int64
		for i = 1; i <= times; i++ {
			// Run the check command
			start := time.Now()

			output.WriteString("Time ")
			output.WriteString(fmt.Sprintf("%d", i))
			output.WriteString(": ")

			// Convert the code to file codes
			fileCodes := converter.ConvertToFileCode(filePath, true)

			// Analyze the project's codebase
			analyzer.AnalyzeCode(fileCodes, filePath, true)

			// Calculate the elapsed time
			elapsed := time.Since(start)

			// Write the elapsed time
			output.WriteString(elapsed.String())
			output.WriteString("\n")
		}

		output.WriteString("== End of ")
		output.WriteString(name)
		output.WriteString(" ==\n\n")
	}

	// Create a file with write-only permissions
	file, err := os.Create(fileName)

	if err != nil || file == nil {
		logger.Error("Could not write the benchmark output to a file")
		logger.Help("Ensure that you have the necessary permissions")
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Could not close the file")
			logger.Help("Ensure that you have the necessary permissions")
		}
	}(file) // Ensure the file is closed

	// Append the output to a file
	_, err = file.WriteString(output.String())

	if err != nil {
		logger.Error("Could not write the benchmark output to a file")
		logger.Help("Ensure that you have the necessary permissions")
		return
	}

	// Stop the spinner
	err = spinner.Stop()

	if err != nil {
		logger.Error("Could not write stop the spinner")
		logger.Help("This is most likely an internal issue")
		logger.Help("If it keeps happening, feel free to open an issue")
		return
	}

	fmt.Print(
		fmt.Sprintf(
			"%s=> The benchmark results have been written to:\n %s%s\n",
			ansi.BoldBrightGreen,
			fileName,
			ansi.Reset,
		),
	)

}
