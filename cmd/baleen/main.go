/*
Package main serves as the primary entry point for launching the Baleen
command line application.
*/
package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/rotationalio/baleen"
	"github.com/rotationalio/baleen/config"
	"github.com/urfave/cli/v2"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	// Create a new CLI app
	app := cli.NewApp()
	app.Name = "baleen"
	app.Version = baleen.Version()
	app.Usage = "a toolkit for ingesting data from RSS feeds"

	// Define commands available to the application
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "run the baleen ingestion service",
			Action: run,
			Flags:  []cli.Flag{},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func run(c *cli.Context) (err error) {
	var svc *baleen.Baleen
	if svc, err = baleen.New(config.Config{}); err != nil {
		return cli.Exit(err, 1)
	}

	if err = svc.Run(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}
