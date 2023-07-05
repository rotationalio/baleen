package events_test

import (
	"os"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/rotationalio/baleen/events"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	generateFixtures := false
	if os.Getenv("BALEEN_TEST_GENERATE_FIXTURE") == "1" {
		generateFixtures = true
	}

	t.Run("Subscription", func(t *testing.T) {
		t.Parallel()
		sub := &events.Subscription{
			FeedID:   watermill.NewUUID(),
			Title:    "Test Subscription",
			FeedType: "RSS",
			FeedURL:  "https://example.com/rss",
			SiteURL:  "http://example.com",
		}

		msg, err := events.Marshal(sub, watermill.NewUUID())
		require.NoError(t, err, "could not marshal subscription")

		if generateFixtures {
			os.WriteFile("testdata/subscription.msgp", []byte(msg.Payload), 0644)
		}

		cmp, err := events.UnmarshalSubscription(msg)
		require.NoError(t, err, "could not unmarshal subscription")

		require.Equal(t, sub, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("FeedSync", func(t *testing.T) {
		t.Parallel()

		fsync := &events.FeedSync{
			FeedID:       watermill.NewULID(),
			ETag:         watermill.NewULID(),
			LastModified: time.Now().Add(-2313 * time.Second).Format(time.RFC3339Nano),
			Active:       true,
			StatusCode:   200,
			SyncedAt:     time.Now().Truncate(time.Microsecond),
			FeedItems:    12,
			Title:        "Test Subscription",
			Description:  "This is an example RSS subscription",
			Link:         "https://example.com",
			FeedLink:     "https://example.com/rss",
			Updated:      time.Now().Add(-2313 * time.Second).Format(time.RFC3339),
			Published:    time.Now().Add(-2313 * time.Second).Format(time.RFC3339),
			Language:     "en-US",
			Copyright:    "Creative Commons 3.0",
			Generator:    "Test Fixture",
			Categories:   []string{"test", "example"},
			FeedType:     "RSS",
			FeedVersion:  "2.0",
		}

		msg, err := events.Marshal(fsync, watermill.NewUUID())
		require.NoError(t, err, "could not marshal feed sync")

		if generateFixtures {
			os.WriteFile("testdata/feedsync.msgp", []byte(msg.Payload), 0644)
		}

		cmp, err := events.UnmarshalFeedSync(msg)
		require.NoError(t, err, "could not unmarshal feed sync")

		require.NotZero(t, cmp.SyncedAt)
		require.Equal(t, fsync, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("FeedItem", func(t *testing.T) {
		t.Parallel()

		item := &events.FeedItem{
			FeedID:      watermill.NewULID(),
			Title:       "Thoughts on Testing Examples",
			Description: "A blog post about creating effective test fixtures.",
			Content:     "It is very important to get test examples correct. This blog posts describes how to do it right.",
			Link:        "https://example.com/blog/testing-examples.html",
			Updated:     time.Now().Add(-2313 * time.Second).Format(time.RFC3339),
			Published:   time.Now().Add(-2313 * time.Second).Format(time.RFC3339),
			GUID:        watermill.NewUUID(),
			Authors:     []string{"John E. Quincy", "Mary Anne Tester"},
			Categories:  []string{"tests", "examples"},
		}

		msg, err := events.Marshal(item, watermill.NewUUID())
		require.NoError(t, err, "could not marshal feed item")

		if generateFixtures {
			os.WriteFile("testdata/feeditem.msgp", []byte(msg.Payload), 0644)
		}

		cmp, err := events.UnmarshalFeedItem(msg)
		require.NoError(t, err, "could not unmarshal feed item")

		require.Equal(t, item, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("Document", func(t *testing.T) {
		t.Parallel()

		// TODO: populate data
		doc := &events.Document{
			ETag:         watermill.NewULID(),
			LastModified: time.Now().Add(-2313 * time.Second).Format(time.RFC3339),
			Content:      []byte("<html><head><title>Thoughts on Testing Examples</title><head><body><h1>Thoughts on Testing Examples</h1><p>A blog post about creating effective test fixtures.</p></body></html>"),
			Active:       true,
			StatusCode:   200,
			FetchedAt:    time.Now().Truncate(time.Millisecond),
			FeedID:       watermill.NewULID(),
			Language:     "en-US",
			Year:         2023,
			Month:        "July",
			Day:          5,
			Title:        "Thoughts on Testing Examples",
			Description:  "A blog post about creating effective test fixtures.",
			Encoding:     "UTF-8",
			Link:         "https://example.com/blog/testing-examples.html",
		}

		msg, err := events.Marshal(doc, watermill.NewUUID())
		require.NoError(t, err, "could not marshal document")

		if generateFixtures {
			os.WriteFile("testdata/document.msgp", []byte(msg.Payload), 0644)
		}

		cmp, err := events.UnmarshalDocument(msg)
		require.NoError(t, err, "could not unmarshal document")

		require.NotZero(t, cmp.FetchedAt)
		require.Equal(t, doc, cmp, "unmarshaled and marshaled message do not match")
	})
}
