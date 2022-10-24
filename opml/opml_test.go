package opml_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/rotationalio/baleen/opml"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Table driven tests for JSON and XML parsing
	paths := []string{
		"testdata/feedly.opml",
		"testdata/feedly.json",
	}

	for _, path := range paths {
		ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), "."))
		t.Run(ext, func(t *testing.T) {
			outline, err := opml.Load(path)
			require.NoError(t, err, "could not load outline")
			require.NotNil(t, outline.Body, "body was nil in the loaded outline")
			require.Len(t, outline.Body.Outlines, 11, "did not properly parse the outlines")
		})
	}
}

func TestURLs(t *testing.T) {
	outline, err := opml.Load("testdata/feedly.opml")
	require.NoError(t, err, "could not load fixture data")

	urls := outline.URLs()
	require.Len(t, urls, 10, "should have had all urls returned from feed that aren't blank")

	urls = outline.URLs("atom", "foo")
	require.Len(t, urls, 5, "should be able to filter on atom feeds")

	urls = outline.URLs("rss", "atom")
	require.Len(t, urls, 10, "should be able to filter all urls")

	urls = outline.URLs("foo", "bar", "baz")
	require.Len(t, urls, 0, "should be able to filter out all urls")
}
