package main

import (
	cli2 "fluent/cli"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "fluent",
		Usage: "A blazingly fast programming language",
		Commands: []*cli.Command{
			{
				Name:    "benchmark",
				Usage:   "Runs the benchmarking tool",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "times",
						Value:    1000,
						Usage:    "How many times should the code be tested",
						Aliases:  []string{"t"},
						Required: false,
					},
				},
				Action: func(context *cli.Context) error {
					cli2.BenchmarkCommand(context)
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
				Action: func(context *cli.Context) error {
					cli2.RunCommand(context)
					return nil
				},
			},
			{
				Name:    "check",
				Usage:   "Checks a fluent file",
				Aliases: []string{"c"},
				Action: func(context *cli.Context) error {
					cli2.CheckCommand(context)
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
				Action: func(context *cli.Context) error {
					cli2.LicenseCommand(context)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}
