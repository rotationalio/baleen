/*
Package fetch_test provides testing for the functions in the fetch package.
*/
package fetch_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rotationalio/baleen/fetch"
	"github.com/stretchr/testify/require"
)

// Helper function for the httptest server to return RSS test data.
func serveTestdata(t *testing.T, path string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			t.Errorf("could not open testdata data: %s", err)
		}

		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			t.Errorf("could not stat RSS data: %s", err)
		}

		// Set Headers
		rw.Header().Set("Content-Type", "text/xml")
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", d.Size()))
		rw.WriteHeader(200)

		if _, err = io.Copy(rw, f); err != nil {
			t.Errorf("could not copy RSS data to resp: %s", err)
		}

	}
}

func TestRSSResponse(t *testing.T) {
	// Start a local test HTTP server and close when test is done
	server := httptest.NewServer(serveTestdata(t, "testdata/rss2.xml"))
	defer server.Close()

	// Set the default client to the test server client.
	fetch.SetDefaultClient(server.Client())

	// Fetch the RSS from the server
	fetcher := fetch.NewFeedFetcher(server.URL)
	feed, err := fetcher.Fetch()
	require.NoError(t, err)
	require.Equal(t, feed.FeedType, "rss")
	require.Equal(t, feed.Title, "Sample Feed")
	require.Equal(t, len(feed.Items), 1)
}

func TestAtomResponse(t *testing.T) {
	// Start a local test HTTP server and close when test is done
	server := httptest.NewServer(serveTestdata(t, "testdata/atom1.xml"))
	defer server.Close()

	// Set the default client to the test server client.
	fetch.SetDefaultClient(server.Client())

	// Fetch the Atom from the server
	fetcher := fetch.NewFeedFetcher(server.URL)
	feed, err := fetcher.Fetch()
	require.NoError(t, err)
	require.Equal(t, feed.FeedType, "atom")
	require.Equal(t, feed.Title, "Sample Feed")
	require.Equal(t, len(feed.Items), 1)
}

func TestSendETag(t *testing.T) {
	// Make one reqeust that gets an etag response
	// subsequent request should contain etag (and respond with 304)

	// Start a local test HTTP server and close when test is done
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// If request has an etag, then send not modified
		if etag := req.Header.Get("If-None-Match"); etag == "ABCDEFG" {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.Header().Set("ETag", "ABCDEFG")
		serveTestdata(t, "testdata/atom1.xml")(rw, req)
	}))
	defer server.Close()

	// Set the default client to the test server client.
	fetch.SetDefaultClient(server.Client())

	fetcher := fetch.NewFeedFetcher(server.URL)

	// The first fetch should return the feed
	feed, err := fetcher.Fetch()
	require.NoError(t, err)
	require.Equal(t, feed.Title, "Sample Feed")

	// The second fetch should return 304
	feed, err = fetcher.Fetch()
	he, ok := err.(fetch.HTTPError)
	require.True(t, ok, "did not return an HTTPError on detection of etag")
	require.Equal(t, he.Code, http.StatusNotModified)
	require.True(t, feed == nil, "feed is not nil")
}

func TestSendLastModified(t *testing.T) {
	// Make one reqeust that gets an last-modified response
	// subsequent request should contain last-modified (and respond with 304)

	// Start a local test HTTP server and close when test is done
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// If request has an etag, then send not modified
		if modified := req.Header.Get("If-Modified-Since"); modified == "Wed, 21 Oct 2015 07:28:00 GMT" {
			rw.WriteHeader(http.StatusNotModified)
			return
		}

		rw.Header().Set("Last-Modified", "Wed, 21 Oct 2015 07:28:00 GMT")
		serveTestdata(t, "testdata/rss2.xml")(rw, req)
	}))
	defer server.Close()

	// Set the default client to the test server client.
	fetch.SetDefaultClient(server.Client())

	fetcher := fetch.NewFeedFetcher(server.URL)

	// The first fetch should return the feed
	feed, err := fetcher.Fetch()
	require.NoError(t, err)
	require.Equal(t, feed.Title, "Sample Feed")

	// The second fetch should return 304
	feed, err = fetcher.Fetch()
	he, ok := err.(fetch.HTTPError)
	require.True(t, ok, "did not return an HTTPError on detection of last modified")
	require.Equal(t, he.Code, http.StatusNotModified)
	require.True(t, feed == nil, "feed is not nil")
}
