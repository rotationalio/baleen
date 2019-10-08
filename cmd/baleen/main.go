/*
Package main serves as the primary entry point for launching the Baleen
command line application.
*/
package main

import (
	"fmt"
	"os"

	"github.com/kansaslabs/baleen"
	"github.com/kansaslabs/baleen/fetch"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	// Create a new CLI app
	app := cli.NewApp()
	app.Name = "baleen"
	app.Version = baleen.Version(false)
	app.Usage = "a toolkit for ingesting data from RSS feeds"

	// Define commands available to the application
	app.Commands = []cli.Command{
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
	if c.NArg() == 0 {
		return cli.NewExitError("specify a feed to fetch", 1)
	}

	url := c.Args()[0]
	fetcher := fetch.New(url)
	feed, err := fetcher.Fetch()

	if err != nil {
		switch he := err.(type) {
		case fetch.HTTPError:
			if he.NotFound() {
				fmt.Println("the url you supplied was not valid")
			}
			if he.NotModified() {
				fmt.Println("no new items in the feed!")
			}
			return cli.NewExitError(err, 1)
		default:
			return cli.NewExitError(err, 1)
		}
	}

	fmt.Println(feed.Title)
	fmt.Println(feed.Description + "\n")

	for _, item := range feed.Items {
		fmt.Println(item.Title)
		fmt.Println(item.Description + "\n")
	}

	return nil
}

/*
Author:  Benjamin Bengfort
Author:  Rebecca Bilbro
Created: Thu Apr 25 18:32:19 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: main.go [68a2562] benjamin@bengfort.com $
*/
