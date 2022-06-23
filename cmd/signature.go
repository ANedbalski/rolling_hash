package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"rollingHash/app"
	"rollingHash/app/algo"
	"rollingHash/app/storage"
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

	out, err := os.OpenFile(c.Args().Get(1), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return fmt.Errorf("caanot create signature-file %s \n %w", c.Args().Get(1), err)
	}
	defer out.Close()

	sig := app.NewSignature(algo.NewAdler32(), algo.MD5, 512)
	sig.Calc(in)
	if err != nil {
		return fmt.Errorf("signature-file cannot be empty")
	}

	err = storage.NewSignatureStorage(out, nil).Store(sig)
	if err != nil {
		return fmt.Errorf("signature-file cannot be empty")
	}

	return nil
}
