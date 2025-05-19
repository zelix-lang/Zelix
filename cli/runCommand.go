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
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
)

// RunCommand represents the run command of the Fluent CLI
// it runs a fluent file
func RunCommand(context *cli.Command) {
	finalPath := BuildCommand(context, false)

	// Directly call the Fluent interpreter
	cmd := exec.Command("fluenti", "--path", finalPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return
	}
}
