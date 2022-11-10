package baleen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mmcdole/gofeed"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/events"
	"github.com/rotationalio/baleen/fetch"
	mime "github.com/rotationalio/ensign/pkg/mimetype/v1beta1"
	"github.com/rs/zerolog/log"
)

func (s *Baleen) AddFeedSync(conf config.FeedSyncConfig, publisher message.Publisher) (err error) {
	var fsync *FeedSync
	if fsync, err = NewFeedSync(conf, publisher); err != nil {
		return err
	}

	// Add the handler to handle messages from the subscriptions topic.
	handler := s.router.AddHandler(
		"feed_sync",
		TopicSubscriptions,
		s.subscriber,
		TopicFeeds,
		s.publisher,
		fsync.Handle,
	)

	// Filter the type of messages handled
	handler.AddMiddleware(
		TypeFilter(mime.ApplicationMsgPack.MimeType(), events.TypeSubscription),
	)

	// Add the plugin to start the fsync routine when the router is run.
	s.router.AddPlugin(fsync.Start)
	return nil
}

func NewFeedSync(conf config.FeedSyncConfig, publisher message.Publisher) (*FeedSync, error) {
	if !conf.Enabled {
		return nil, errors.New("feed sync is not enabled")
	}

	return &FeedSync{
		conf:      conf,
		manifest:  make(Manifest),
		stop:      make(chan struct{}),
		publisher: publisher,
	}, nil
}

type FeedSync struct {
	conf      config.FeedSyncConfig
	publisher message.Publisher
	manifest  Manifest
	stop      chan struct{}
}

func (f *FeedSync) Handle(msg *message.Message) (_ []*message.Message, err error) {
	// Parse the subscription
	var info *events.Subscription
	if info, err = events.UnmarshalSubscription(msg); err != nil {
		return nil, err
	}

	// If there is no URL then just ignore the event
	if info.FeedURL == "" {
		return nil, nil
	}

	// Create or update the feed in the manifest
	feed := f.manifest.Add(info)

	// Synchronize the feed right now
	return feed.Sync()
}

func (f *FeedSync) Start(r *message.Router) error {
	if f.conf.Interval < time.Second {
		return errors.New("interval must be 1s or greater")
	}

	go func() {
		// Setup the feed sync background routine
		ticker := time.NewTicker(f.conf.Interval)

		// Wait until the router starts running to start the feed sync process.
		<-r.Running()
		log.Info().Dur("interval", f.conf.Interval).Msg("feed_sync interval is running")
		defer log.Info().Msg("feed_sync interval has stopped")

		for {
			// TODO: when the next version of watermill comes out, also select on handler.Stopped()
			select {
			case <-f.stop:
				return
			case <-ticker.C:
			}

			log.Info().Int("nfeeds", len(f.manifest)).Msg("synchronizing feeds")

			// Handle subscriptions
			for _, feed := range f.manifest {
				msgs, err := feed.Sync()
				if err != nil {
					log.Error().Err(err).Str("feed_id", feed.info.FeedID).Str("url", feed.info.FeedURL).Msg("could not synchronize feed")
					continue
				}

				if err = f.publisher.Publish(TopicFeeds, msgs...); err != nil {
					log.Error().Err(err).Int("num", len(msgs)).Str("feed_id", feed.info.FeedID).Str("url", feed.info.FeedURL).Msg("could not publish feed messages")
					continue
				}
			}
		}
	}()
	return nil
}

func (f *FeedSync) Stop() {
	close(f.stop)
}

type Manifest map[string]*Feed

type Feed struct {
	info    *events.Subscription
	fetcher *fetch.FeedFetcher
}

// Add or update the feed to the manifest
func (m Manifest) Add(info *events.Subscription) *Feed {
	// Update the feed with the new info
	if feed, ok := m[info.FeedURL]; ok {
		if feed.info.FeedID == "" || (info.FeedID != "" && feed.info.FeedID != info.FeedID) {
			feed.info.FeedID = info.FeedID
		}

		if feed.info.FeedType == "" || (info.FeedType != "" && feed.info.FeedType != info.FeedType) {
			feed.info.FeedType = info.FeedType
		}

		if feed.info.SiteURL == "" || (info.SiteURL != "" && feed.info.SiteURL != info.SiteURL) {
			feed.info.SiteURL = info.SiteURL
		}

		return feed
	}

	// Create the Feed and return it
	if info.FeedID == "" {
		info.FeedID = watermill.NewShortUUID()
	}

	feed := &Feed{
		info:    info,
		fetcher: fetch.NewFeedFetcher(info.FeedURL),
	}
	m[info.FeedURL] = feed
	return feed
}

// Sync the feed and return the FeedItem events to publish
func (f *Feed) Sync() (msgs []*message.Message, err error) {
	log.Info().Str("feed_id", f.info.FeedID).Str("url", f.info.FeedURL).Msg("synchronizing feed")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	var rss *gofeed.Feed
	if rss, err = f.fetcher.Fetch(ctx); err != nil {
		if httperr, ok := err.(*fetch.HTTPError); ok {
			// If it is an http error emit an fsync event
			fsync := &events.FeedSync{
				FeedID:     f.info.FeedID,
				Active:     false,
				Error:      httperr.Status,
				StatusCode: httperr.Code,
				SyncedAt:   time.Now(),
				Title:      f.info.Title,
				Link:       f.info.FeedURL,
				FeedType:   f.info.FeedType,
			}

			var msg *message.Message
			if msg, err = events.Marshal(fsync, watermill.NewULID()); err != nil {
				return nil, err
			}

			return []*message.Message{msg}, nil
		}

		return nil, err
	}

	msgs = make([]*message.Message, 0, len(rss.Items)+1)

	fsync := &events.FeedSync{
		FeedID:       f.info.FeedID,
		ETag:         f.fetcher.ETag(),
		LastModified: f.fetcher.Modified(),
		Active:       true,
		SyncedAt:     time.Now(),
		FeedItems:    int64(len(rss.Items)),
		Title:        rss.Title,
		Description:  rss.Description,
		Link:         rss.Link,
		Links:        rss.Links,
		FeedLink:     rss.FeedLink,
		Updated:      rss.Updated,
		Published:    rss.Published,
		Language:     rss.Language,
		Copyright:    rss.Copyright,
		Generator:    rss.Generator,
		Categories:   rss.Categories,
		FeedType:     rss.FeedType,
		FeedVersion:  rss.FeedVersion,
	}

	var msg *message.Message
	if msg, err = events.Marshal(fsync, watermill.NewULID()); err != nil {
		return nil, err
	}
	msgs = append(msgs, msg)

	// Handle each feed item
	for _, item := range rss.Items {
		fitem := &events.FeedItem{
			FeedID:      f.info.FeedID,
			Title:       item.Title,
			Description: item.Description,
			Content:     item.Content,
			Link:        item.Link,
			Updated:     item.Updated,
			Published:   item.Published,
			GUID:        item.GUID,
			Categories:  item.Categories,
		}

		if item.Image != nil {
			fitem.Image = item.Image.URL
		}

		fitem.Authors = make([]string, 0, len(item.Authors))
		for _, author := range item.Authors {
			var name string
			switch {
			case author.Name != "" && author.Email != "":
				name = fmt.Sprintf("%s <%s>", author.Name, author.Email)
			case author.Name != "":
				name = author.Name
			case author.Email != "":
				name = author.Email
			}

			if name != "" {
				fitem.Authors = append(fitem.Authors, name)
			}
		}

		fitem.Enclosures = make([]string, 0, len(item.Enclosures))
		for _, enclosure := range item.Enclosures {
			fitem.Enclosures = append(fitem.Enclosures, enclosure.URL)
		}

		var msg *message.Message
		if msg, err = events.Marshal(fitem, watermill.NewULID()); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
