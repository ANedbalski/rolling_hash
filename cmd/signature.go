package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"rollingHash/app"
)

var (
	signatureCommand = &cli.Command{
		Name:    "signature",
		Aliases: []string{"s"},
		Usage:   "signature old-file signature-file",
		Action:  signatureAction,
	}
)

func signatureAction(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return fmt.Errorf("wrong number of arguments")
	}

	if c.Args().Get(0) == "" {
		return fmt.Errorf("old-file cannot be empty")
	}

	if c.Args().Get(1) == "" {
		return fmt.Errorf("signature-file cannot be empty")
	}

	in, err := os.Open(c.Args().Get(0))
	if err != nil {
		return fmt.Errorf("cannot open old-file %s \n %w", c.Args().Get(0), err)
	}
	defer in.Close()

	out, err := os.Open(c.Args().Get(1))
	if err != nil {
		return fmt.Errorf("caanot create signature-file %s \n %w", c.Args().Get(1), err)
	}
	defer out.Close()

	_, err = app.NewSignature(in, out)

	return err
}
