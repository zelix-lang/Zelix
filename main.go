/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package main

import (
	"context"
	cli2 "fluent/cli"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	app := &cli.Command{
		Name:  "fluent",
		Usage: "A blazingly fast programming language",
		Commands: []*cli.Command{
			{
				Name:    "benchmark",
				Usage:   "Runs the benchmarking tool",
				Aliases: []string{"bE"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "times",
						Value:    1000,
						Usage:    "How many times should the code be tested",
						Aliases:  []string{"t"},
						Required: false,
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					cli2.BenchmarkCommand(cmd)
					return nil
				},
			},
			{
				Name:    "build",
				Usage:   "Compiles your project into an executable",
				Aliases: []string{"b"},
				Flags:   []cli.Flag{},
				Action: func(_ context.Context, cmd *cli.Command) error {
					cli2.BuildCommand(cmd)
					return nil
				},
			},
			{
				Name:    "run",
				Usage:   "Runs a fluent file",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "timer",
						Value:   false,
						Usage:   "Prints the time taken for each step",
						Aliases: []string{"t"},
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					cli2.RunCommand(cmd)
					return nil
				},
			},
			{
				Name:    "check",
				Usage:   "Checks a fluent file",
				Aliases: []string{"c"},
				Action: func(_ context.Context, cmd *cli.Command) error {
					cli2.CheckCommand(cmd)
					return nil
				},
			},
			{
				Name:    "license",
				Usage:   "Prints the license",
				Aliases: []string{"l"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "full",
						Value:   false,
						Usage:   "Prints the full license",
						Aliases: []string{"f"},
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					cli2.LicenseCommand(cmd)
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
