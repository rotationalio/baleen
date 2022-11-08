package fetch_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rotationalio/baleen/fetch"
	"github.com/stretchr/testify/require"
)

func TestRSSResponse(t *testing.T) {
	// Create a test server serving rss2 data
	url := NewServer(t, FixtureHandler(t, "testdata/rss2.xml"))

	// Fetch the RSS from the server
	fetcher := fetch.NewFeedFetcher(url)
	feed, err := fetcher.Fetch(context.Background())
	require.NoError(t, err)
	require.Equal(t, feed.FeedType, "rss")
	require.Equal(t, feed.Title, "Sample Feed")
	require.Equal(t, len(feed.Items), 1)
}

func TestAtomResponse(t *testing.T) {
	// Create a test server serving atom1 data
	url := NewServer(t, FixtureHandler(t, "testdata/atom1.xml"))

	// Fetch the Atom from the server
	fetcher := fetch.NewFeedFetcher(url)
	feed, err := fetcher.Fetch(context.Background())
	require.NoError(t, err)
	require.Equal(t, feed.FeedType, "atom")
	require.Equal(t, feed.Title, "Sample Feed")
	require.Equal(t, len(feed.Items), 1)
}

func TestSendETag(t *testing.T) {
	// Make one reqeust that gets an etag response
	// subsequent request should contain etag (and respond with 304)
	handler := func(rw http.ResponseWriter, req *http.Request) {
		// If request has an etag, then send not modified
		if etag := req.Header.Get("If-None-Match"); etag == "ABCDEFG" {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.Header().Set("ETag", "ABCDEFG")
		FixtureHandler(t, "testdata/atom1.xml")(rw, req)
	}
	url := NewServer(t, handler)

	fetcher := fetch.NewFeedFetcher(url)

	// The first fetch should return the feed
	feed, err := fetcher.Fetch(context.Background())
	require.NoError(t, err)
	require.Equal(t, feed.Title, "Sample Feed")

	// The second fetch should return 304
	feed, err = fetcher.Fetch(context.Background())
	he, ok := err.(fetch.HTTPError)
	require.True(t, ok, "did not return an HTTPError on detection of etag")
	require.Equal(t, he.Code, http.StatusNotModified)
	require.True(t, feed == nil, "feed is not nil")
}

func TestSendLastModified(t *testing.T) {
	// Make one reqeust that gets an last-modified response
	// subsequent request should contain last-modified (and respond with 304)

	// Start a local test HTTP server and close when test is done
	handler := func(rw http.ResponseWriter, req *http.Request) {
		// If request has an etag, then send not modified
		if modified := req.Header.Get("If-Modified-Since"); modified == "Wed, 21 Oct 2015 07:28:00 GMT" {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.Header().Set("Last-Modified", "Wed, 21 Oct 2015 07:28:00 GMT")
		FixtureHandler(t, "testdata/rss2.xml")(rw, req)
	}
	url := NewServer(t, handler)

	fetcher := fetch.NewFeedFetcher(url)

	// The first fetch should return the feed
	feed, err := fetcher.Fetch(context.Background())
	require.NoError(t, err)
	require.Equal(t, feed.Title, "Sample Feed")

	// The second fetch should return 304
	feed, err = fetcher.Fetch(context.Background())
	he, ok := err.(fetch.HTTPError)
	require.True(t, ok, "did not return an HTTPError on detection of last modified")
	require.Equal(t, he.Code, http.StatusNotModified)
	require.True(t, feed == nil, "feed is not nil")
}

func TestFeedError(t *testing.T) {
	url := NewServer(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
	})

	fetcher := fetch.NewFeedFetcher(url)
	feed, err := fetcher.Fetch(context.Background())
	require.Error(t, err)
	require.Nil(t, feed)

	herr, ok := err.(fetch.HTTPError)
	require.True(t, ok, "expected an http error returned")
	require.Equal(t, http.StatusBadRequest, herr.Code)
	require.NotEmpty(t, herr.Status)
}
