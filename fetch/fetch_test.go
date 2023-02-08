/*
Package fetch_test provides testing for the functions in the fetch package.
*/
package fetch_test

import (
	"bytes"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andybalholm/brotli"
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
		fetch.HeaderContentType,
		fetch.HeaderContentEncoding,
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

// Helper function for the httptest server to return HTML test data.
func FixtureHandler(t *testing.T, path string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			t.Logf("could not open %s fixture: %s", path, err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			t.Logf("could not stat fixture: %s", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
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

// Helper function for the httptest server to return compressed HTML test data.
func CompressedFixtureHandler(t *testing.T, path, compression string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			t.Logf("could not open %s fixture: %s", path, err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var (
			w   io.WriteCloser
			buf bytes.Buffer
		)
		switch compression {
		case "gzip", "frog":
			w = gzip.NewWriter(&buf)
		case "br":
			w = brotli.NewWriter(&buf)
		case "compress":
			w = lzw.NewWriter(&buf, lzw.MSB, 8)
		case "deflate":
			w = zlib.NewWriter(&buf)
		default:
			t.Logf("unknown compression format %q", compression)
			http.Error(rw, fmt.Sprintf("unknown compression format %q", compression), http.StatusInternalServerError)
			return
		}

		if _, err = io.Copy(w, f); err != nil {
			t.Logf("could not compress fixture: %s", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Close the compression writer to finalize compression
		if err = w.Close(); err != nil {
			t.Logf("could not compress fixture: %s", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set Headers
		rw.Header().Set("Content-Type", "text/xml")
		rw.Header().Set("Content-Encoding", compression)
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
		rw.WriteHeader(http.StatusOK)

		if _, err = io.Copy(rw, &buf); err != nil {
			t.Logf("could not copy compressed fixture data to resp: %s", err)
		}
	}
}
