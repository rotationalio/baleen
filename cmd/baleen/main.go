package main

import (
	"os"

	"github.com/kansaslabs/baleen"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	// Create a new CLI app
	app := cli.NewApp()
	app.Name = "baleen"
	app.Version = baleen.Version(false)

	// Run the CLI app
	app.Run(os.Args)
}
