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

package cli

import (
	"fluent/ansi"
	"fluent/util"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"regexp"
	"strings"
)

const downRightArrow = "↳"
const arrow = "➜"
const cross = "✖"
const check = "✔"

// printQuestion prints a formatted question to the console.
// Parameters:
// - question: The question text to be printed.
// - qColor: The color code for the question text.
// - symbol: The symbol to be printed before the question.
// - color: The color code for the symbol.
func printQuestion(question, qColor string, symbol string, color string) {
	fmt.Print(
		color,
		symbol,
		ansi.Reset,
		" ",
		qColor,
		question,
		" ",
		ansi.BrightBlack,
		arrow,
		ansi.Reset,
		" ",
	)
}

// handleErr checks if an error is not nil and panics if it is.
// Parameters:
// - err: The error to be checked.
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// clearFallback clears the fallback text from the console.
// Parameters:
// - buffer: A pointer to a slice of runes representing the current input buffer.
// - fallback: A pointer to a string representing the fallback text.
// - addingNew: A boolean indicating whether new input is being added.
func clearFallback(buffer *[]rune, fallback *string, addingNew bool) {
	if len(*buffer) > 0 && !addingNew {
		return
	}

	for i := 0; i < len(*fallback)+2; i++ {
		fmt.Print("\b \b")
	}
}

// ask prompts the user with a question and reads the input from the keyboard.
// Parameters:
// - question: The question text to be printed.
// - initialColor: The color code for the question text.
// - fallback: The fallback text to be used if the user does not provide input.
// - pattern: A regular expression pattern to validate the input against (if usePattern is true).
// Returns:
// - The user's input as a string, or the fallback text if no input is provided.
func ask(
	question,
	initialColor,
	fallback string,
	pattern *regexp.Regexp,
) string {
	usePattern := pattern != nil
	// Start listening for keyboard events
	handleErr(keyboard.Open())

	// Ask the user indefinitely until we get valid input
	for {
		// Create a buffer for the response
		buffer := make([]rune, 0)

		// Print the question
		printQuestion(question, initialColor, downRightArrow, ansi.BoldBrightYellow)
		fmt.Print(ansi.BrightBlack, "(", fallback, ")", ansi.Reset)

		// Read the keys pressed by the user
		for {
			char, key, err := keyboard.GetKey()

			if err != nil {
				handleErr(keyboard.Close())
				panic(err)
			}

			// Check for exit
			if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
				clearFallback(&buffer, &fallback, true)
				fmt.Print(ansi.Colorize(ansi.BoldBrightRed, "Process interrupted by the user"))
				fmt.Println()
				handleErr(keyboard.Close())
				os.Exit(0)
			}

			// Check if we have to process the response
			if key == keyboard.KeyEnter {
				handleErr(keyboard.Close())
				clearFallback(&buffer, &fallback, len(buffer) == 0)
				fmt.Print("\r")
				// Convert the buffer to a string
				str := string(buffer)

				// Convert to fallback if the buffer is empty
				if len(buffer) == 0 {
					str = fallback
				}

				// Check if we have to use a pattern
				if usePattern {
					// Check if the string matches the pattern
					if !pattern.MatchString(str) {
						printQuestion(question, ansi.BoldBrightRed, cross, ansi.BoldBrightRed)
						fmt.Printf("%s(%s)%s", ansi.BrightBlack, str, ansi.Reset)
						fmt.Println()
						continue
					}
				}

				printQuestion(question, ansi.BoldBrightGreen, check, ansi.BoldBrightGreen)
				fmt.Printf("%s(%s)%s", ansi.BrightBlack, str, ansi.Reset)
				fmt.Println()
				return strings.ReplaceAll(str, "\"", "\\\"")
			}

			// Handle the backspace key
			if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
				if len(buffer) > 0 {
					buffer = buffer[:len(buffer)-1]
					fmt.Print("\b \b")

					if len(buffer) == 0 {
						fmt.Print(ansi.BrightBlack, "(", fallback, ")", ansi.Reset)
					}
				}

				continue
			}

			// Handle the space key
			if key == keyboard.KeySpace {
				clearFallback(&buffer, &fallback, len(buffer) == 0)
				fmt.Printf(" ")
				buffer = append(buffer, ' ')
				continue
			}

			clearFallback(&buffer, &fallback, len(buffer) == 0)
			fmt.Printf("%c", char)
			buffer = append(buffer, char)
		}
	}

	return ""
}

func InitCommand(context *cli.Command) {
	ShowHeaderMessage()
	fmt.Print(ansi.BrightBlack)
	fmt.Println("This command will guide you through the")
	fmt.Println("process of creating a new Fluent project.")
	fmt.Println("If you wish to exit the process at any time,")
	fmt.Println("Press Ctrl+C or Esc.")
	fmt.Println(ansi.Reset)

	// Get the current working directory
	target, err := os.Getwd()
	handleErr(err)

	// Get the selected path (if any)
	preferredPath := context.Args().First()

	// Update the target path if a preferred path is provided
	if preferredPath != "" {
		target = path.Join(target, preferredPath)
	}

	// Check if the target is a directory
	if !util.DirExists(target) {
		fmt.Print(ansi.BoldBrightRed, "The target is not a directory.", ansi.Reset)
		fmt.Println()
		os.Exit(1)
	}

	packagePath := path.Join(target, "package.fluent")
	// Check if a package.fluent already exists in this directory
	if _, err := os.Stat(packagePath); err == nil {
		fmt.Print(ansi.BoldBrightRed, "A package.fluent file already exists in this directory.", ansi.Reset)
		fmt.Println()
		os.Exit(1)
	}

	name := ask("Name", ansi.BoldBrightBlue, "Project", nil)
	desc := ask("Description", ansi.BoldBrightBlue, "A Fluent project", nil)
	author := ask("Author", ansi.BoldBrightBlue, "John Doe", nil)
	license := ask("License", ansi.BoldBrightBlue, "MIT", regexp.MustCompile("^[a-zA-Z-.0-9]+$"))
	entry := ask("Entry file", ansi.BoldBrightBlue, "main.fluent", fluentExtensionRegex)
	fmt.Println()

	// Write the package.fluent file
	file, err := os.Create(packagePath)
	handleErr(err)

	// Make sure to close the file
	defer func(file *os.File) {
		handleErr(file.Close())
	}(file)

	// Write all strings to a builder and write in one go
	builder := strings.Builder{}

	builder.WriteString("// The Fluent Programming Language\n")
	builder.WriteString("// This file was generated by the Fluent CLI\n")
	builder.WriteString("// Feel free to modify it as needed\n")
	builder.WriteString("// -----------------------------------------------------\n\n")
	builder.WriteString("Name = \"")
	builder.WriteString(name)
	builder.WriteString("\"\n")
	builder.WriteString("Description = \"")
	builder.WriteString(desc)
	builder.WriteString("\"\n")
	builder.WriteString("Author = \"")
	builder.WriteString(author)
	builder.WriteString("\"\n")
	builder.WriteString("License = \"")
	builder.WriteString(license)
	builder.WriteString("\"\n")
	builder.WriteString("Entry = \"")
	builder.WriteString(entry)
	builder.WriteString("\"\n")

	// Save the file
	_, err = file.WriteString(builder.String())
	handleErr(err)

	// Print a success message
	fmt.Print(ansi.BoldBrightGreen, "The package.fluent file was successfully created.", ansi.Reset)
	fmt.Println()
	fmt.Print(ansi.BoldBrightPurple, "What now?.", ansi.Reset)
	fmt.Println()
	fmt.Print(ansi.BoldBrightBlue, "- Run 'fluent r' to run your code", ansi.Reset)
	fmt.Println()
	fmt.Print(ansi.BoldBrightBlue, "- Run 'fluent c' to compile your code", ansi.Reset)
	fmt.Println()
	fmt.Print(ansi.BoldBrightBlue, "- Happy coding!", ansi.Reset)
	fmt.Println()
}
