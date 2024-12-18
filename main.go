package main

import (
	"github.com/urfave/cli/v2"
	"os"
	cli2 "surf/cli"
)

func main() {
	app := &cli.App{
		Name:  "surf",
		Usage: "A blazingly fast programming language",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Usage:   "Runs a surf file",
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
				Usage:   "Checks a surf file",
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
