package cli

import (
	"bufio"
	"fluent/ansi"
	"fluent/logger"
	"fluent/util"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// cliQuestion is a helper function to ask the user a Y/N question
func cliQuestion(message string) bool {
	// Loop until the user inserts something valid
	for {
		// Create a new scanner
		scanner := bufio.NewScanner(os.Stdin)

		// Print the message without a newline
		fmt.Print(
			ansi.Colorize(
				"green_bright",
				message,
			),
		)

		// Print the options
		fmt.Print(
			ansi.Colorize(
				"yellow_bright",
				" [Y/n] ",
			),
		)

		// Scan the input
		if scanner.Scan() {
			input := strings.ToLower(scanner.Text())

			if input == "y" || input == "yes" {
				return true
			} else if input == "n" || input == "no" {
				return false
			} else {
				fmt.Println(
					ansi.Colorize(
						"red_bright",
						"Invalid input. Try again.",
					),
				)
			}
		}
	}
}

// setupDirectory sets up the benchmark and example directories
func setupDirectory() {
	if util.FileExists("benchmark") {
		if !util.IsDir("benchmark") {
			deleteFile := cliQuestion("A file named benchmark already exists, do you want to delete it?")

			if !deleteFile {
				logger.Error("Benchmarking aborted.")
				return
			}

			// Delete the file
			err := os.Remove("benchmark")
			if err != nil {
				logger.Error("Error deleting the file")
				return
			}
		}

		// Ask the user if they want to delete the directory
		deleteDir := cliQuestion("A directory named benchmark already exists, do you want to delete it?")

		if deleteDir {
			// Delete the directory
			err := os.RemoveAll("benchmark")
			if err != nil {
				logger.Error("Error deleting the directory")
				return
			}

			// Create the directory
			err = os.Mkdir("benchmark", 0755)

			if err != nil {
				logger.Error("Error creating the benchmark directory")
				return
			}
		}
	} else {
		// Create the directory
		err := os.Mkdir("benchmark", 0755)

		if err != nil {
			logger.Error("Error creating the benchmark directory")
			return
		}
	}

	// Make sure the example directory exists
	if !util.FileExists("example") || !util.IsDir("example") {
		logger.Error("Example directory does not exist, please make sure you didn't accidentally delete it")
		return
	}

	// Make sure the example directory is not empty
	files, err := os.ReadDir("example")

	if err != nil {
		logger.Error("Error reading the example directory")
		return
	}

	if len(files) == 0 {
		logger.Error("Example directory is empty, please make sure you have examples to benchmark")
		return
	}
}

// BenchmarkCommand runs the benchmark tool
func BenchmarkCommand(context *cli.Context) {
	// Print the welcome header
	fmt.Println(
		ansi.Colorize(
			"blue_bright",
			"The Fluent Programming Language",
		),
	)

	fmt.Println(
		ansi.Colorize(
			"black_bright",
			"Welcome to the benchmarking tool",
		),
	)

	// See if there is already a benchmark directory
	setupDirectory()

	// Print an empty line before the rest of the output
	// to differentiate it from the welcome message
	fmt.Println()

	// Get the times flag
	times := context.Int("times")
	timesStr := strconv.Itoa(times)

	// Make a slice holding all the durations
	durations := make(map[string][]time.Duration)

	// Iterate over the files to check
	files, err := os.ReadDir("example")

	if err != nil {
		logger.Error("Error reading the example directory")
		return
	}

	// The result to be written into a file
	result := strings.Builder{}

	// Add some metadata
	result.WriteString("The Fluent Programming Language\n")
	result.WriteString("--------------------------------\n")
	result.WriteString("\n")
	result.WriteString("Started at: " + time.Now().String() + "\n")
	result.WriteString("Times: " + timesStr + "\n")
	result.WriteString("Running on " + runtime.GOOS + " " + runtime.GOARCH + "\n")
	result.WriteString("\n")

	for _, file := range files {
		if !file.Type().IsRegular() {
			logger.Warning("Skipping non-regular file: ", file.Name())
			continue
		}

		// Print the file name
		fmt.Println(
			ansi.Colorize(
				"green_bright",
				"=== "+file.Name()+" ===",
			),
		)

		// Write the contents of the file for reference
		fileContents, err := os.ReadFile("example/" + file.Name())

		if err != nil {
			logger.Error("Error reading the file")
			return
		}

		// Also write to the result
		result.WriteString("=== " + file.Name() + " ===\n")

		// Write the contents of the file
		result.WriteString(string(fileContents) + "\n")

		results := make([]time.Duration, times)
		for i := 1; i <= times; i++ {
			start := time.Now()

			fmt.Print(
				ansi.Colorize(
					"magenta_bright",
					"Time "+fmt.Sprint(i)+"/"+timesStr,
				),
				// Add a "\r" to overwrite the line
				"\r",
			)

			// Spawn the task
			cmd := exec.Command("go", "run", "main.go", "c", "example/"+file.Name())
			err := cmd.Run()

			if err != nil {
				logger.Error("Error invoking the main file, make sure you are in the project's root directory")
				return
			}

			// Append the duration to the slice
			taken := time.Since(start)

			// Write to the result
			result.WriteString("Time " + fmt.Sprint(i) + ": " + taken.String() + "\n")
			results[i-1] = taken
		}

		// Add the results into the map
		durations[file.Name()] = results

		// Print a newline to finish the line
		fmt.Println()

		// Print the end
		fmt.Println(
			ansi.Colorize(
				"green_bright",
				"=== END OF "+file.Name()+" ===",
			),
		)

		// Calculate the average
		var sum time.Duration

		for _, duration := range results {
			sum += duration
		}

		average := sum / time.Duration(times)

		// Write to the result
		result.WriteString("Average: " + average.String() + "\n")

		// Also write to the result
		result.WriteString("=== END OF " + file.Name() + " ===\n")
	}

	// Add when the benchmarking ended
	result.WriteString("\n")

	now := time.Now()
	result.WriteString("Ended at: " + now.String() + "\n")

	// Write the result to a file
	timeStamp := strconv.Itoa(now.Year()) + "-" + now.Month().String() + "-" + strconv.Itoa(now.Day())
	timeStamp += "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute()) + "-" + strconv.Itoa(now.Second())
	resultFile, err := os.Create("benchmark/" + timeStamp + "-bench.txt")

	if err != nil {
		logger.Error("Error creating the result file")
		fmt.Println(err)
		return
	}

	_, err = resultFile.WriteString(result.String())

	if err != nil {
		logger.Error("Error writing the result to the file")
		return
	}

	// Notify the user that the file has been written
	fmt.Println(
		ansi.Colorize(
			"green_bright",
			"The result has been written to benchmark/"+timeStamp+"-bench.txt",
		),
	)
}
