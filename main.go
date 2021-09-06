package main

import (
	"github.com/pinealctx/renault/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var app = cli.NewApp()
	app.Name = "Renault"
	app.Usage = "A useful tools."
	app.Version = "0.0.1"
	app.Commands = cli.Commands{cmd.InitCommand}
	var err = app.Run(os.Args[:])
	if err != nil {
		log.Fatalln(err)
	}
}
