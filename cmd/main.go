package cmd

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Usage: "compare files",
		Commands: []*cli.Command{
			runCommand,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "rh.yml",
				Usage:   "path to config file",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Application failed: %v", err)
	}

}
