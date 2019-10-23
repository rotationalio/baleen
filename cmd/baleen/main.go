package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/kansaslabs/baleen"
	"github.com/kansaslabs/baleen/fetch"
	"github.com/kansaslabs/baleen/store"
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

	// We have feeds! So let's make sure we can connect to S3 and create a session
	config := store.AWSCredentials{
		Region: os.Getenv("AWS_REGION"),
		Bucket: os.Getenv("KANSAS_BUCKET"),
	}

	session, err := store.GetSession(&config)
	if err != nil {
		log.Println(err)
	}
	if session == nil {
		panic("could not connect to s3")
	} else {
		log.Println("connected to s3")
	}

	// We're connected to S3 so let's iterate over our urls and fetch them
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
				// If it's not one of the above errors, print out the error and stop execution
				// TODO: Add better handling for this in case it's just a temporarily lost internet connection?
				return cli.NewExitError(err, 1)
			}
		}

		for _, item := range feed.Items {

			var year int
			var month string
			var day int

			if item.PublishedParsed == nil {
				// Some feed have the date formatted incorrectly (no day)
				// In this case, we'll just infer that it's today's year, month, day
				currentTime := time.Now()
				year = currentTime.Year()
				month = currentTime.Month().String()
				day = currentTime.Day()
			} else {
				year = item.PublishedParsed.Year()
				month = item.PublishedParsed.Month().String()
				day = item.PublishedParsed.Day()
			}
			feedID := slug.Make(feed.Title)

			// TODO: This doesn't seem to reliably retrieve text for all items, sometimes it's just a bunch of links
			content := fetch.GetContent(item.Link)

			// TODO: Hash the content to see if it exists already in the manifest & if so, skip

			doc := store.Document{
				FeedID:       feedID,
				LanguageCode: slug.Make(feed.Language),
				Year:         year,
				Month:        month,
				Day:          day,
				Title:        item.Title,
				Description:  item.Description,
				Link:         item.Link,
				Content:      content,
			}

			// Using the open session, upload the document to the bucket
			err = store.Upload(session, doc, config.Bucket)
			if err != nil {
				log.Println(err)
			}
			// TODO: If there's no error so far, add the hashed content to the manifest
		}
	}
	return nil
}
