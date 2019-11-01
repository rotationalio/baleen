/*
Package main serves as the primary entry point for launching the Baleen
command line application.
*/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/kansaslabs/baleen"
	"github.com/kansaslabs/baleen/fetch"
	"github.com/kansaslabs/baleen/store"
	"github.com/kansaslabs/baleen/utils"
	"github.com/spaolacci/murmur3"
	"github.com/syndtr/goleveldb/leveldb"
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

	// Retrieve the manifest so that we don't re-ingest docs we already have
	var db *leveldb.DB
	db, err = store.OpenManifest("./db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// We're connected to S3 so let's iterate over our urls and fetch them
	for _, url := range urls {

		feedFetcher := fetch.NewFeedFetcher(url)
		feed, err := feedFetcher.Fetch()
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
				// If it's not one of the above errors, print out the error but don't stop execution
				// Looks like the culprit is usually either a blip in internet connection or bad XML encoding
				continue
			}
		}

		// If we failed to get a feed, just skip it
		if feed == nil {
			break
		}

		for _, item := range feed.Items {

			// Hash the title to get a key to lookup or add to the manifest
			hasher := murmur3.New64()
			hasher.Write([]byte(item.Title))
			key := strconv.FormatInt(int64(hasher.Sum64()), 10)

			// Look up the item's key, if it exists, we have the item already, so can skip
			if _, err := db.Get([]byte(key), nil); err == nil {
				continue
			} else {
				// Otherwise prepare to retrieve and store the full details and text of the item
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

				htmlFetcher := fetch.NewHTMLFetcher(item.Link)
				content, err := htmlFetcher.Fetch()
				// TODO: Better error handling
				if err != nil {
					fmt.Println(err)
				}

				// Make the doc, store it & add to the manifest
				doc := store.Document{
					FeedID:       slug.Make(feed.Title),
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

				// If there's no error so far, add to the manifest where the key is the hash and the value is a DB write timestamp
				now := time.Now().String()
				db.Put([]byte(key), []byte(now), nil)
			}
		}
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
