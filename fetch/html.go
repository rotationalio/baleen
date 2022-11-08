package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HTMLFetcher is an interface for fetching the full HTML associated with a feed item
type HTMLFetcher struct {
	url string // the url of the article full text
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
func (f *HTMLFetcher) Fetch(ctx context.Context) (raw []byte, err error) {
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
	html, err := io.ReadAll(rep.Body)
	if err != nil {
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
