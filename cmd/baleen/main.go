/*
Package main serves as the primary entry point for launching the Baleen
command line application.
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/joho/godotenv"
	"github.com/rotationalio/baleen"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/events"
	"github.com/rotationalio/baleen/logger"
	"github.com/rotationalio/baleen/opml"
	"github.com/rotationalio/watermill-ensign/pkg/ensign"
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
			Before: configure,
			Action: run,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "feeds:add",
			Usage:  "add a feed subscription to baleen if its not already added",
			Before: mkpub,
			After:  rmpub,
			Action: addFeed,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "add a new subscription via its xml url",
				},
				&cli.StringFlag{
					Name:    "opml",
					Aliases: []string{"o"},
					Usage:   "add subscriptions from an OPML file (json or xml)",
				},
			},
		},
		{
			Name:   "posts:add",
			Usage:  "add posts for document processing",
			Before: mkpub,
			After:  rmpub,
			Action: addPost,
		},
		{
			Name:   "debug",
			Usage:  "subscribe to all topics to debug messages being published",
			Before: configure,
			Action: debug,
			Flags:  []cli.Flag{},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

var (
	conf      config.Config
	publisher message.Publisher
)

func configure(c *cli.Context) (err error) {
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func mkpub(c *cli.Context) (err error) {
	if err = configure(c); err != nil {
		return err
	}

	if publisher, err = baleen.CreatePublisher(conf.Publisher, watermill.NopLogger{}); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func rmpub(c *cli.Context) (err error) {
	if err = publisher.Close(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func run(c *cli.Context) (err error) {
	var svc *baleen.Baleen
	if svc, err = baleen.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = svc.Run(context.Background()); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func addFeed(c *cli.Context) (err error) {
	var nEvents int
	if c.String("url") == "" && c.String("opml") == "" {
		return cli.Exit("specify either -url or -opml to add a feed", 1)
	}

	// Handle single URL case
	if url := c.String("url"); url != "" {
		sub := &events.Subscription{
			FeedURL: url,
		}

		var msg *message.Message
		if msg, err = events.Marshal(sub, watermill.NewULID()); err != nil {
			return cli.Exit(err, 1)
		}

		if err = publisher.Publish(baleen.TopicSubscriptions, msg); err != nil {
			return cli.Exit(err, 1)
		}
		nEvents++
	}

	// Handle OPML case
	if path := c.String("opml"); path != "" {
		var outline *opml.OPML
		if outline, err = opml.Load(path); err != nil {
			return cli.Exit(err, 1)
		}

		for _, feed := range outline.Body.Outlines {
			sub := &events.Subscription{
				FeedType: feed.Type,
				Title:    feed.Title,
				FeedURL:  feed.XMLURL,
				SiteURL:  feed.HTMLURL,
			}

			var msg *message.Message
			if msg, err = events.Marshal(sub, watermill.NewULID()); err != nil {
				return cli.Exit(err, 1)
			}

			if err = publisher.Publish(baleen.TopicSubscriptions, msg); err != nil {
				return cli.Exit(err, 1)
			}
			nEvents++
		}
	}

	fmt.Printf("published %d subscription events\n", nEvents)
	return nil
}

func addPost(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.Exit("specify at least one url", 1)
	}

	var nEvents int
	for i := 0; i < c.NArg(); i++ {
		fitem := &events.FeedItem{
			Link: c.Args().Get(i),
		}

		var msg *message.Message
		if msg, err = events.Marshal(fitem, watermill.NewULID()); err != nil {
			return cli.Exit(err, 1)
		}

		if err = publisher.Publish(baleen.TopicFeeds, msg); err != nil {
			return cli.Exit(err, 1)
		}
		nEvents++
	}

	fmt.Printf("published %d feed item events\n", nEvents)
	return nil
}

func debug(c *cli.Context) (err error) {
	var subscriber message.Subscriber
	if subscriber, err = baleen.CreateSubscriber(conf.Subscriber, logger.New()); err != nil {
		return cli.Exit(err, 1)
	}
	defer subscriber.Close()

	subs, _ := subscriber.Subscribe(context.Background(), baleen.TopicSubscriptions)

	for msg := range subs {
		etype := msg.Metadata.Get(ensign.TypeNameKey)
		size := len(msg.Payload)
		log.Printf("%s - %d bytes", etype, size)

		msg.Ack()
	}

	return nil
}
