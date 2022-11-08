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

func TestCanonicalHeaders(t *testing.T) {
	headers := []string{
		fetch.HeaderUserAgent,
		fetch.HeaderAccept,
		fetch.HeaderAcceptLang,
		fetch.HeaderAcceptEncode,
		fetch.HeaderCacheControl,
		fetch.HeaderReferer,
		fetch.HeaderIfNoneMatch,
		fetch.HeaderIfModifiedSince,
		fetch.HeaderLastModified,
	}

	for _, header := range headers {
		require.Equal(t, header, http.CanonicalHeaderKey(header), "header does not match canonical header, change constant")
	}
}

// Helper function to create an http test server and set the fetch client.
func NewServer(t *testing.T, handler http.HandlerFunc) string {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	fetch.SetClient(server.Client())
	t.Logf("test server open at %s", server.URL)
	return server.URL
}

// Helper function to create an http TLS test server and set the fetch client.
func NewTLSServer(t *testing.T, handler http.HandlerFunc) string {
	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)
	fetch.SetClient(server.Client())
	t.Logf("test server open at %s", server.URL)
	return server.URL
}

// Helper function for the httptest server to return RSS test data.
func FixtureHandler(t *testing.T, path string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			t.Logf("could not open %s fixture: %s", path, err)
			errs := err.Error()
			rw.Header().Set("Content-Type", "text/plain")
			rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(errs)))
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(errs))
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			t.Logf("could not stat fixture: %s", err)
			errs := err.Error()
			rw.Header().Set("Content-Type", "text/plain")
			rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(errs)))
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(errs))
			return
		}

		// Set Headers
		rw.Header().Set("Content-Type", "text/xml")
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", d.Size()))
		rw.WriteHeader(http.StatusOK)

		if _, err = io.Copy(rw, f); err != nil {
			t.Logf("could not copy fixture data to resp: %s", err)
		}
	}
}
