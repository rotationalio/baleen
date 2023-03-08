package baleen

import (
	"context"
	"errors"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rotationalio/baleen/config"
	"github.com/rotationalio/baleen/events"
	"github.com/rotationalio/baleen/fetch"
	mime "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/rs/zerolog/log"
)

func (s *Baleen) AddPostFetch(conf config.PostFetchConfig) error {
	if !conf.Enabled {
		return errors.New("post fetch is not enabled")
	}

	// Add the handler to handle messages from the subscriptions topic.
	handler := s.router.AddHandler(
		"post_fetch",
		TopicFeeds,
		s.subscriber,
		TopicDocuments,
		s.publisher,
		PostFetch,
	)

	// Filter the type of messages handled
	handler.AddMiddleware(
		TypeFilter(mime.ApplicationMsgPack.MimeType(), events.TypeFeedItem),
	)

	return nil
}

func PostFetch(msg *message.Message) (_ []*message.Message, err error) {
	var event *events.FeedItem
	if event, err = events.UnmarshalFeedItem(msg); err != nil {
		return nil, err
	}

	if event.Link == "" {
		return nil, nil
	}

	log.Info().Str("feed_id", event.FeedID).Str("url", event.Link).Msg("fetching post")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	doc := &events.Document{
		FetchedAt: time.Now(),
		Active:    true,
		FeedID:    event.FeedID,
	}

	var html *fetch.HTML
	fetcher := fetch.NewHTMLFetcher(event.Link)
	if html, err = fetcher.Fetch(ctx); err != nil {
		log.Warn().Err(err).Str("url", event.Link).Str("feed_id", event.FeedID).Msg("could not fetch post")
		httperr, ok := err.(*fetch.HTTPError)
		if !ok {
			return nil, err
		}

		// If we receive an http error pass an document error event on.
		doc.Active = false
		doc.StatusCode = httperr.Code
		doc.Error = httperr.Status
	}

	if doc.Content, err = html.Extract(); err != nil {
		log.Warn().Err(err).Str("url", event.Link).Str("feed_id", event.FeedID).Msg("could not decode post")
		return nil, err
	}

	doc.Title = html.Title()
	doc.Description = html.Description()
	doc.Link = event.Link

	var out *message.Message
	if out, err = events.Marshal(doc, watermill.NewULID()); err != nil {
		return nil, err
	}

	return []*message.Message{out}, nil
}
