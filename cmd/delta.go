package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	deltaCommand = &cli.Command{
		Name:    "delta",
		Aliases: []string{"d"},
		Usage:   "delta signature-file new-file delta-file",
		Action:  deltaAction,
	}
)

func deltaAction(c *cli.Context) error {
	if c.Args().Len() != 3 {
		return fmt.Errorf("incorrect number of arguments")
	}

	if c.Args().Get(0) == "" {
		return fmt.Errorf("signature-file cannot be empty")
	}

	if c.Args().Get(1) == "" {
		return fmt.Errorf("new-file cannot be empty")
	}

	if c.Args().Get(2) == "" {
		return fmt.Errorf("delta-file cannot be empty")
	}

	sig, err := os.Open(c.Args().Get(0))
	if err != nil {
		return fmt.Errorf("canot open signature-file %s \n %w", c.Args().Get(0), err)
	}
	defer sig.Close()

	f, err := os.Open(c.Args().Get(1))
	if err != nil {
		return fmt.Errorf("canot open new-file %s \n %w", c.Args().Get(1), err)
	}
	defer f.Close()

	diff, err := os.Open(c.Args().Get(2))
	if err != nil {
		return fmt.Errorf("canot create diff-file %s \n %w", c.Args().Get(2), err)
	}
	defer diff.Close()

	return nil
}
