package fetch

import (
	"bytes"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/brotli"
)

// HTMLFetcher is an interface for fetching the full HTML associated with a feed item
type HTMLFetcher struct {
	url string // the url of the article full text
}

// HTML is an in-memory materialized view of an HTML document fetched by the HTMLFetcher.
// It has helper methods to decode and parse the contents of the response, particularly
// if that response is compressed or encoded in non UTF-8 string encoding.
type HTML struct {
	content     *bytes.Buffer
	ctype       string
	encoding    string
	title       string
	description string
}

// NewHTMLFetcher creates a new HTML fetcher that can fetch the full HTML from the specified URL.
func NewHTMLFetcher(url string) *HTMLFetcher {
	return &HTMLFetcher{
		url: url,
	}
}

// The HTMLFetcher uses GET requests to retrieve the html containing the full text
// of articles of feeds with a Baleen-specific http client.
// TODO: return an HTML file instead of simply raw bytes (including document data).
func (f *HTMLFetcher) Fetch(ctx context.Context) (html *HTML, err error) {
	var req *http.Request
	if req, err = f.newRequest(ctx); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = client.Do(req); err != nil {
		return nil, err
	}

	// Close the body of the response reader when we're done.
	if rep != nil && rep.Body != nil {
		defer rep.Body.Close()
	}

	// Check the status code of the response; note that 304 means not modified, but we
	// are still returning a 304 error to signal to the Subscription that nothing has
	// changed and that the post is nil.
	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return nil, HTTPError{
			Status: rep.Status,
			Code:   rep.StatusCode,
		}
	}

	// Materialize the HTML content from the body
	buf := make([]byte, 0, rep.ContentLength)
	html = &HTML{
		content:  bytes.NewBuffer(buf),
		ctype:    rep.Header.Get(HeaderContentType),
		encoding: rep.Header.Get(HeaderContentEncoding),
	}

	if written, err := io.Copy(html.content, rep.Body); err != nil || written != rep.ContentLength {
		return nil, fmt.Errorf("no text parsed from html retrieved from %s: %w", f.url, err)
	}
	return html, nil
}

func (f *HTMLFetcher) newRequest(ctx context.Context) (req *http.Request, err error) {
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil); err != nil {
		return nil, err
	}

	req.Header.Set(HeaderUserAgent, userAgent)
	req.Header.Set(HeaderAccept, acceptHTML)
	req.Header.Set(HeaderAcceptLang, acceptLang)
	req.Header.Set(HeaderAcceptEncode, acceptEncode)
	req.Header.Set(HeaderCacheControl, cacheControl)
	req.Header.Set(HeaderReferer, referer)

	return req, nil
}

// Extract handles compression and content encoding from the response.
func (h *HTML) Extract() (_ []byte, err error) {
	var reader io.ReadCloser
	if reader, err = h.extract(); err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (h *HTML) extract() (io.ReadCloser, error) {
	buf := bytes.NewBuffer(nil)
	tee := io.TeeReader(h.content, buf)
	h.content = buf

	switch h.encoding {
	case gzipEncode:
		return gzip.NewReader(tee)
	case brotliEncode:
		return io.NopCloser(brotli.NewReader(tee)), nil
	case lzwEncode:
		// TODO: what values to use for order and width?
		return lzw.NewReader(tee, lzw.MSB, 8), nil
	case zlibEncode:
		return zlib.NewReader(tee)
	case "", "identity":
		return io.NopCloser(tee), nil
	default:
		return nil, fmt.Errorf("unknown content encoding %q", h.encoding)
	}
}

func (h *HTML) Title() string {
	if h.title == "" {
		if err := h.parse(); err != nil {
			panic(err)
		}
	}
	return h.title
}

func (h *HTML) Description() string {
	if h.description == "" {
		h.parse()
	}
	return h.description
}

func (h *HTML) parse() (err error) {
	var reader io.ReadCloser
	if reader, err = h.extract(); err != nil {
		return err
	}
	defer reader.Close()

	var tree *goquery.Document
	if tree, err = goquery.NewDocumentFromReader(reader); err != nil {
		return err
	}

	h.title = tree.Find("title").Contents().Text()
	tree.Find("meta").EachWithBreak(func(index int, item *goquery.Selection) bool {
		if item.AttrOr("name", "") == "description" {
			h.description = item.AttrOr("content", "")
			return h.description == ""
		}
		return true
	})

	return nil
}
