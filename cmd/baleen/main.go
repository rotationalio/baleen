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

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gosimple/slug"
	"github.com/rotationalio/baleen"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/fetch"
	"github.com/rotationalio/baleen/opml"
	"github.com/rotationalio/baleen/publish"
	"github.com/rotationalio/baleen/store"
	"github.com/spaolacci/murmur3"
	"github.com/urfave/cli/v2"
)

func main() {
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
	var files []string
	var urls []string
	var conf config.Config

	if conf, err = config.New(); err != nil {
		return err
	}

	// If the user specifies a feed via the command line, only get that one
	if c.NArg() > 0 {
		urls = append(urls, c.Args().First())
	} else {
		// Otherwise retrieve feeds from files in the fixtures directory
		err := filepath.Walk(conf.FixturesDir, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			switch filepath.Ext(file) {
			case ".opml":
				o := opml.OPML{}
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
		return cli.Exit("specify a feed to fetch or add feeds to fixtures directory", 1)
	}

	var session *session.Session
	if conf.AWS.Enabled {
		// We have feeds! So let's make sure we can connect to S3 and create a session
		creds := store.AWSCredentials{
			Region: conf.AWS.Region,
			Bucket: conf.AWS.Bucket,
		}

		session, err = store.GetSession(&creds)
		if err != nil {
			log.Println(err)
		}
		if session == nil {
			panic("could not connect to s3")
		} else {
			log.Println("connected to s3")
		}
	}

	var publisher *publish.KafkaPublisher
	if conf.Kafka.Enabled {
		if publisher, err = publish.New(conf.Kafka); err != nil {
			panic(err)
		}
	}

	// Retrieve the manifest so that we don't re-ingest docs we already have
	db := store.MustOpen(conf.DBPath)
	defer db.Close()

	// We're connected to S3 so let's iterate over our urls and fetch them
	for _, url := range urls {

		feedFetcher := fetch.NewFeedFetcher(url)
		feed, err := feedFetcher.Fetch()
		if err != nil {
			switch he := err.(type) {
			case fetch.HTTPError:
				switch {
				case he.NotModified():
					fmt.Printf("no new items in the specified feed!: %s\n", url)
				case he.Forbidden():
					fmt.Printf("unable to access feed (forbidden): %s\n", url)
				case he.NotFound():
					fmt.Printf("the url you supplied was not valid: %s\n", url)
				default:
					fmt.Printf("unrecognized HTTP error %d\n", he.Code)
				}
			default:
				// If it's not an HTTP error, print out the error but don't stop execution
				// Looks like the culprit is usually either a blip in internet connection or bad XML encoding
				fmt.Printf("unrecognized fetch error: %s\n", err.Error())
			}

			if publisher != nil {
				// Write feed error
				feed := &store.Feed{
					URL:    url,
					Active: false,
					Error:  err.Error(),
				}
				if err = publisher.WriteFeed(feed); err != nil {
					fmt.Printf("unable to compose Kafka feed message: %s\n", err.Error())
				}
			}
			continue
		}

		// If we failed to get a feed, just skip it
		if feed == nil {
			break
		}

		if publisher != nil {
			// Write feed active
			feed := &store.Feed{
				URL:    url,
				Active: true,
			}
			if err = publisher.WriteFeed(feed); err != nil {
				fmt.Printf("unable to compose Kafka feed message: %s\n", err.Error())
			}
		}

		for _, item := range feed.Items {

			// Hash the title to get a key to lookup or add to the manifest
			hasher := murmur3.New64()
			hasher.Write([]byte(item.Title))
			key := strconv.FormatInt(int64(hasher.Sum64()), 10)

			// Look up the item's key, if it exists, we have the item already, so can skip
			if _, err := db.Get([]byte(key), nil); err == nil && !conf.Testing {
				continue
			} else {
				// TODO: Detect encoding so that we can set the Encoding on the Document that gets written to S3

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
				html, err := htmlFetcher.Fetch()

				// TODO: Better error handling
				if err != nil {
					fmt.Println(err)
				}

				var languageCode string
				if feed.Language != "" {
					languageCode = slug.Make(feed.Language)
				} else {
					languageCode = "unknown"
				}

				// Make the doc, store it & add to the manifest
				doc := store.Document{
					FeedID:       slug.Make(feed.Title),
					LanguageCode: languageCode,
					Year:         year,
					Month:        month,
					Day:          day,
					Title:        item.Title,
					Description:  item.Description,
					Link:         item.Link,
					Content:      html,
				}

				if conf.AWS.Enabled {
					// Using the open session, upload the document to the bucket
					err = store.Upload(session, doc, conf.AWS.Bucket)
					if err != nil {
						log.Println(err)
					}
				}

				if publisher != nil {
					if err = publisher.WriteDocument(&doc); err != nil {
						fmt.Printf("unable to compose Kafka document message: %s\n", err.Error())
					}
				}

				// If there's no error so far, add to the manifest where the key is the hash and the value is a DB write timestamp
				now := time.Now().String()
				db.Put([]byte(key), []byte(now), nil)
			}
		}
	}

	if publisher != nil {
		if err = publisher.PublishMessages(); err != nil {
			fmt.Printf("unable to publish some Kafka messages: %s\n", err.Error())
		}
	}

	return nil
}
