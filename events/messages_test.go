package events_test

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/rotationalio/baleen/events"
	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
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

		cmp, err := events.UnmarshalSubscription(msg)
		require.NoError(t, err, "could not unmarshal subscription")

		require.Equal(t, sub, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("FeedSync", func(t *testing.T) {
		t.Parallel()

		// TODO: populate data
		fsync := &events.FeedSync{
			SyncedAt: time.Now(),
		}

		msg, err := events.Marshal(fsync, watermill.NewUUID())
		require.NoError(t, err, "could not marshal feed sync")

		cmp, err := events.UnmarshalFeedSync(msg)
		require.NoError(t, err, "could not unmarshal feed sync")

		require.NotZero(t, cmp.SyncedAt)
		// TODO: deal with timestamp comparisons
		// require.Equal(t, fsync, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("FeedItem", func(t *testing.T) {
		t.Parallel()

		// TODO: populate data
		item := &events.FeedItem{}

		msg, err := events.Marshal(item, watermill.NewUUID())
		require.NoError(t, err, "could not marshal feed item")

		cmp, err := events.UnmarshalFeedItem(msg)
		require.NoError(t, err, "could not unmarshal feed item")

		require.Equal(t, item, cmp, "unmarshaled and marshaled message do not match")
	})

	t.Run("Document", func(t *testing.T) {
		t.Parallel()

		// TODO: populate data
		doc := &events.Document{
			FetchedAt: time.Now(),
		}

		msg, err := events.Marshal(doc, watermill.NewUUID())
		require.NoError(t, err, "could not marshal document")

		cmp, err := events.UnmarshalDocument(msg)
		require.NoError(t, err, "could not unmarshal document")

		require.NotZero(t, cmp.FetchedAt)
		// TODO: deal with timestamp comparisons
		// require.Equal(t, doc, cmp, "unmarshaled and marshaled message do not match")
	})
}
