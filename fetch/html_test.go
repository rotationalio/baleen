package fetch_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rotationalio/baleen/fetch"
	"github.com/stretchr/testify/require"
)

func TestHTMLResponse(t *testing.T) {
	// Create a test server serving rss2 data
	url := NewServer(t, FixtureHandler(t, "testdata/post.html"))

	// Fetch the RSS from the server
	fetcher := fetch.NewHTMLFetcher(url)
	html, err := fetcher.Fetch(context.Background())
	require.NoError(t, err)

	data, err := html.Extract()
	require.NoError(t, err)
	require.Len(t, data, 1048)

	require.Equal(t, "Hello World Post", html.Title())
	require.Equal(t, "Just a quick test post", html.Description())
}

func TestHTMLError(t *testing.T) {
	url := NewServer(t, func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
	})

	fetcher := fetch.NewHTMLFetcher(url)
	data, err := fetcher.Fetch(context.Background())
	require.Error(t, err)
	require.Nil(t, data)

	herr, ok := err.(fetch.HTTPError)
	require.True(t, ok, "expected an http error returned")
	require.Equal(t, http.StatusBadRequest, herr.Code)
	require.NotEmpty(t, herr.Status)
}
