package app

import "github.com/urfave/cli/v2"

func CliCommand() *cli.Command {
	return &cli.Command{
		Name:  "app",
		Usage: "Run the service",
		Subcommands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start application",
				Action: func(c *cli.Context) error {
					return Start()
				},
			},
		},
	}
}
