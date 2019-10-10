/*
Package fetch provides a high-level interface for going out to get RSS and Atom feeds
from any source. Right now http and https requests are supported but future
implementations may also include authenticated fetchers, etc. Fetchers are intended to
synchronously make a single request to get the latest version the feed and to preserve
the state of the last request to the feed they are managing. They ensure that
connections are closed after each request use etag and last-modified headers to minimize
the amount of bandwidth required.

Basic Usage:

	fetcher := fetch.New("https://www.example.com/rss")
	feed, err := fetcher.Fetch()

For more on RSS hacking and bandwidth minimization see:
https://fishbowl.pastiche.org/2002/10/21/http_conditional_get_for_rss_hackers
*/
package fetch

import (
	"net"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

func init() {
	client = &http.Client{
		Timeout: 1 * time.Minute,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 45 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 45 * time.Second,
			DisableKeepAlives:   true,
			DisableCompression:  false,
		},
	}
}

// A package level http client for making requests. It is best practice to not use the
// default http.Client but to use your own with timeouts correctly specified. The
// package also admonishes us to only create one client for efficiency because the
// client is itself thread safe.The client is initialized by init() and can be modified
// using the SetDefaultClient function (e.g. for testing). All HTTP based fetchers
// should use this client.
var client *http.Client

// Fetcher provides a interface for anything that can get RSS data and provide it in a
// sequential fashion (e.g. without concurrency). The fetcher is the building block for
// larger subscription routines that periodically use the fetcher to retrieve data.
// Fetchers should therefore be treated as things that will only run inside of a single
// thread, whereas Subscription objects are things that may run concurrently.
type Fetcher interface {
	Fetch() (feed *gofeed.Feed, err error)
}

// New creates a new HTTP fetcher that can fetch rss feeds from the specified URL.
func New(url string) Fetcher {
	return &httpFetcher{
		url:    url,
		parser: gofeed.NewParser(),
	}
}

// SetDefaultClient allows you to specify an alternative http.Client to the default one
// used by all http based Fetchers in this package. Use this function to change the
// timeouts of the client or to set a test client.
func SetDefaultClient(c *http.Client) {
	client = c
}

//===========================================================================
// HTTP Fetcher
//===========================================================================

// The httpFetcher uses GET requests to retrieve data with a Baleen-specific http
// client. We avoid using gofeed.ParseURL because it is very simple and doesn't respect
// rate limits or etags, which are necessary for Baleen to run in continuous operation.
type httpFetcher struct {
	url      string         // the url of the RSS or atom feed
	parser   *gofeed.Parser // the universal feed parser for RSS and Atom feeds
	etag     string         // used for conditional http to minimize bandwidth
	modified string         // used for conditional http to minimize bandwidth
}

func (f *httpFetcher) Fetch() (feed *gofeed.Feed, err error) {
	var req *http.Request
	if req, err = f.newRequest(); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = client.Do(req); err != nil {
		return nil, err
	}

	// Close the body of the response reader when we're done.
	if rep != nil {
		defer func() {
			ce := rep.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
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
	if feed, err = f.parser.Parse(rep.Body); err != nil {
		return nil, err
	}

	// Get the eTag and last-modified from the response header if we've successfully
	// parsed the request and received a 200 response.
	f.etag = rep.Header.Get("ETag")
	f.modified = rep.Header.Get("Last-Modified")

	// Note the explicit return of err here, this is in case the Body.Close() returns
	// an error, which will supercede any other errors being returned.
	return feed, err
}

func (f *httpFetcher) newRequest() (req *http.Request, err error) {
	// Create the GET request
	if req, err = http.NewRequest("GET", f.url, nil); err != nil {
		return nil, err
	}

	// Be a good netizen and tell the server who we are and what we're doing
	// TODO: add actual version and system information to the user agent
	req.Header.Set("User-Agent", "Baleen/1.0")

	// Response control headers (request compressed response by default)
	// Note that compression and keep-alives are handled by our default client.
	req.Header.Set("Accept", "application/atom+xml,application/rdf+xml,application/rss+xml,application/x-netcdf,application/xml;q=0.9,text/xml;q=0.2,*/*;q=0.1")

	// Ask the server to refresh the cache if the content is an hour old
	req.Header.Set("Cache-Control", "max-age=3600")
	// Best practice is to leave the referer blank
	req.Header.Set("Referer", "")

	// Send the etag if an etag was sent from the server on a previous request.
	if f.etag != "" {
		req.Header.Set("If-None-Match", f.etag)
	}

	// Send the timestamp if last-modified was sent from the server previously.
	// TODO: make f.modifed a time.Time and add RFC 1123-compliant parsing.
	// SEE: https://github.com/kurtmckee/feedparser/blob/develop/feedparser/http.py#L113
	if f.modified != "" {
		req.Header.Set("If-Modified-Since", f.modified)
	}

	// RFC 3229 support
	req.Header.Set("A-IM", "feed")

	return req, nil
}

/*
Author:  Benjamin Bengfort
Author:  Rebecca Bilbro
Created: Mon Apr 29 06:43:36 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: fetch.go [d6dba70] benjamin@bengfort.com $
*/
