package cmd

import "github.com/urfave/cli/v2"

var (
	runCommand = &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "run app",
		Action:  runAction,
	}
)

func runAction(c *cli.Context) error {
	_ = NewContext()
	return nil
}
