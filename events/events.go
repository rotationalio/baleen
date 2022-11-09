/*
Package events provides data serialization for Baleen-specific events using message
pack - a binary JSON compatible serialization format. Message pack is slightly larger
than protocol buffers or other serialization formats but can be simpler to implement.
*/
package events

//go:generate msgp

import (
	"time"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
)

// Types specifies the event types for Ensign
const (
	TypeSubscription = "Subscription"
	TypeFeedSync     = "FeedSync"
	TypeFeedItem     = "FeedItem"
	TypeDocument     = "Document"
)

// Versions specifies the numeric version for each event type
const (
	VersionSubscription uint32 = 1
	VersionFeedSync     uint32 = 1
	VersionFeedItem     uint32 = 1
	VersionDocument     uint32 = 1
)

// TypedEvents can return their type for Ensign serialization
type TypedEvent interface {
	Type() *api.Type
}

type Subscription struct {
	FeedID   string `msg:"feed_id,omitempty"` // a unique ID for the feed (optional)
	Title    string `msg:"title"`             // the title of the subscription
	FeedType string `msg:"feed_type"`         // either rss or atom
	FeedURL  string `msg:"feed_url"`          // the url to the feed (xmlURL in OPML)
	SiteURL  string `msg:"site_url"`          // the url to the site (htmlURL in OPML)
}

var _ TypedEvent = Subscription{}

type FeedSync struct {
	FeedID       string    `msg:"feed_id"`
	ETag         string    `msg:"etag"`
	LastModified string    `msg:"last_modified"`
	Active       bool      `msg:"active"`
	StatusCode   int       `msg:"status_code"`
	Error        string    `msg:"error"`
	SyncedAt     time.Time `msg:"synced_at"`
	FeedItems    int64     `msg:"feed_items"`
	Title        string    `msg:"title"`
	Description  string    `msg:"description"`
	Link         string    `msg:"link"`
	Links        []string  `msg:"links"`
	FeedLink     string    `msg:"feed_link"`
	Updated      string    `msg:"updated"`
	Published    string    `msg:"published"`
	Language     string    `msg:"language"`
	Copyright    string    `msg:"copyright"`
	Generator    string    `msg:"generator"`
	Categories   []string  `msg:"categories"`
	FeedType     string    `msg:"feed_type"`
	FeedVersion  string    `msg:"feed_version"`
}

var _ TypedEvent = FeedSync{}

type FeedItem struct {
	FeedID      string   `msg:"feed_id"`
	Title       string   `msg:"title"`
	Description string   `msg:"description"`
	Content     string   `msg:"content"`
	Link        string   `msg:"link"`
	Updated     string   `msg:"updated"`
	Published   string   `msg:"published"`
	GUID        string   `msg:"guid"`
	Authors     []string `msg:"authors"`
	Image       string   `msg:"image"`
	Categories  []string `msg:"categories"`
	Enclosures  []string `msg:"enclosures"`
}

var _ TypedEvent = FeedItem{}

type Document struct {
	ETag         string    `msg:"etag,omitempty"`
	LastModified string    `msg:"last_modified,omitempty"`
	Active       bool      `msg:"active"`
	StatusCode   int       `msg:"status_code,omitempty"`
	Error        string    `msg:"error,omitempty"`
	FetchedAt    time.Time `msg:"fetched_at"`
	FeedID       string    `msg:"feed_id"`
	Language     string    `msg:"language"`
	Year         int       `msg:"year"`
	Month        string    `msg:"month"`
	Day          int       `msg:"day"`
	Title        string    `msg:"title"`
	Description  string    `msg:"description"`
	Content      []byte    `msg:"content"`
	Encoding     string    `msg:"encoding"`
	Link         string    `msg:"link"`
}

var _ TypedEvent = Document{}

func (Subscription) Type() *api.Type {
	return &api.Type{
		Name:    TypeSubscription,
		Version: VersionSubscription,
	}
}

func (FeedSync) Type() *api.Type {
	return &api.Type{
		Name:    TypeFeedSync,
		Version: VersionFeedSync,
	}
}

func (FeedItem) Type() *api.Type {
	return &api.Type{
		Name:    TypeFeedItem,
		Version: VersionFeedItem,
	}
}

func (Document) Type() *api.Type {
	return &api.Type{
		Name:    TypeDocument,
		Version: VersionDocument,
	}
}
