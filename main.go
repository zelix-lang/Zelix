package main

import (
	"github.com/urfave/cli/v2"
	"os"
	cli2 "surf/cli"
)

func main() {
	/*file := os.Getenv("SURF_DUMMY_FILE")

	// Read the file
	input, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	start := time.Now()

	// Lex the input
	tokens := lexer.Lex(string(input), file)

	fmt.Println("Lexed in:", time.Since(start))

	start = time.Now()
	// Parse the tokens into a FileCode
	fileCode := ast.Parse(tokens)
	fmt.Println("Wrapped to AST in:", time.Since(start))

	start = time.Now()
	// Analyze the file code
	checker.AnalyzeFileCode(*fileCode, file)
	fmt.Println("Analyzed in:", time.Since(start))

	start = time.Now()
	// Interpret the file code
	core.Interpret(fileCode, file)
	fmt.Println("Interpreted in:", time.Since(start))*/
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
