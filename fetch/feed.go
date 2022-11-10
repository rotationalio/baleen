package fetch

import (
	"context"
	"net/http"

	"github.com/mmcdole/gofeed"
)

// FeedFetcher provides a interface for anything that can get RSS data and provide it in
// a sequential fashion (e.g. without concurrency). The fetcher is the building block
// for larger subscription routines that periodically use the fetcher to retrieve data.
// FeedFetchers should therefore be treated as things that will only run inside of a
// single thread, whereas Subscription objects are things that may run concurrently.
type FeedFetcher struct {
	url      string         // the url of the RSS or atom feed
	parser   *gofeed.Parser // the universal feed parser for RSS and Atom feeds
	etag     string         // used for conditional http to minimize bandwidth
	modified string         // used for conditional http to minimize bandwidth
}

// NewFeedFetcher creates a new HTTP fetcher that can fetch rss feeds from the specified URL.
func NewFeedFetcher(url string) *FeedFetcher {
	return &FeedFetcher{
		url:    url,
		parser: gofeed.NewParser(),
	}
}

// The FeedFetcher uses GET requests to retrieve data with a Baleen-specific http
// client. We avoid using gofeed.ParseURL because it is very simple and doesn't respect
// rate limits or etags, which are necessary for Baleen to run in continuous operation.
func (f *FeedFetcher) Fetch(ctx context.Context) (feed *gofeed.Feed, err error) {
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
	// changed and that the feed is nil.
	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return nil, HTTPError{
			Status: rep.Status,
			Code:   rep.StatusCode,
		}
	}

	// Use the universal parser to parse the Atom or RSS feed
	// Note: Feeds with illegal character codes will not be successfully parsed & return nil here
	if feed, err = f.parser.Parse(rep.Body); err != nil {
		return nil, err
	}

	// Get the eTag and last-modified from the response header if we've successfully
	// parsed the request and received a 200 response.
	f.etag = rep.Header.Get(HeaderETag)
	f.modified = rep.Header.Get(HeaderLastModified)

	// Note the explicit return of err here, this is in case the Body.Close() returns
	// an error, which will supercede any other errors being returned.
	return feed, nil
}

func (f *FeedFetcher) ETag() string {
	return f.etag
}

func (f *FeedFetcher) Modified() string {
	return f.modified
}

func (f *FeedFetcher) newRequest(ctx context.Context) (req *http.Request, err error) {
	// Create the GET request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil); err != nil {
		return nil, err
	}

	// Be a good netizen and tell the server who we are and what we're doing
	req.Header.Set(HeaderUserAgent, userAgent)

	// Response control headers (request compressed response by default)
	// Note that compression and keep-alives are handled by our default client.
	req.Header.Set(HeaderAccept, acceptRSS)
	req.Header.Set(HeaderAcceptEncode, acceptEncode)

	// Ask the server to refresh the cache if the content is an hour old
	req.Header.Set(HeaderCacheControl, cacheControl)

	// Best practice is to leave the referer blank
	req.Header.Set(HeaderReferer, referer)

	// Send the etag if an etag was sent from the server on a previous request.
	if f.etag != "" {
		req.Header.Set(HeaderIfNoneMatch, f.etag)
	}

	// Send the timestamp if last-modified was sent from the server previously.
	// TODO: make f.modifed a time.Time and add RFC 1123-compliant parsing.
	// SEE: https://github.com/kurtmckee/feedparser/blob/develop/feedparser/http.py#L113
	if f.modified != "" {
		req.Header.Set(HeaderIfModifiedSince, f.modified)
	}

	// RFC 3229 support
	req.Header.Set(HeaderRFC3229, aimType)
	return req, nil
}
