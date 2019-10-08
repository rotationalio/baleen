package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/kansaslabs/baleen"
	"github.com/kansaslabs/baleen/fetch"
	"github.com/kansaslabs/baleen/utils"
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
	var root = filepath.Join("fixtures")
	var files []string
	var urls []string

	// If the user specifies a feed via the command line, only get that one
	if c.NArg() > 0 {
		urls = append(urls, c.Args()[0])
	} else {
		// Otherwise retrieve feeds from files in the fixtures directory
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			switch filepath.Ext(file) {
			case ".opml":
				o := utils.OPML{}
				content, _ := ioutil.ReadFile(file)
				err := xml.Unmarshal(content, &o)
				if err != nil {
					panic(err)
				}
				for _, outline := range o.Body.Outlines {
					urls = append(urls, outline.XMLURL)
				}
			case ".json":
				fmt.Println("parsing from json not yet implemented")
			}
		}
	}

	// Return an error is no feed was specified in the cli and none were retrieved from fixtures
	if len(urls) == 0 {
		return cli.NewExitError("specify a feed to fetch or add feeds to fixtures directory", 1)
	}

	for _, url := range urls {
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
	}

	return nil
}
