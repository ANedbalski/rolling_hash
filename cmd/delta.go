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

	diff, err := os.OpenFile(c.Args().Get(2), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return fmt.Errorf("canot create diff-file %s \n %w", c.Args().Get(2), err)
	}
	defer diff.Close()

	signature, err := storage.NewSignatureStorage(nil, sig).Load()
	if err != nil {
		return fmt.Errorf("incorrect signature data: %w", err)
	}

	delta, err := app.NewDelta(algo.NewAdler32(), algo.MD5).Calc(signature, f)
	if err != nil {
		return fmt.Errorf("error calculation delta: %w", err)
	}

	err = storage.NewDeltaStorage(diff).Store(delta.GetRecords())
	if err != nil {
		return fmt.Errorf("error storing delta: %w", err)
	}

	return nil
}
