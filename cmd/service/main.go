package main

import (
	"assignments/simplebank/app"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	application := &cli.App{
		Name: "simplebank",
		Commands: []*cli.Command{
			app.CliCommand(),
		},
	}
	err := application.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
