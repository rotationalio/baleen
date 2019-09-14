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

Author:  Benjamin Bengfort
Created: Mon Apr 29 06:43:36 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: error.go [d6dba70] benjamin@bengfort.com $
*/
package fetch

import (
	"fmt"
	"net/http"
)

// HTTPError contains status information from the request and can be returned as error.
// This type of error is returned from the Fetcher when the server replies successfully
// but without a 200 status. The suggested use of this error is in a switch statement,
// e.g. something like: switch he := err.(type) {case fetch.HTTPError: ... default: ...}
type HTTPError struct {
	Code   int
	Status string
}

// Error implements the error interface and returns a string representation of the err.
func (e HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.Code, e.Status)
}

// NotModified returns true if the error is an HTTP 304
func (e HTTPError) NotModified() bool {
	return e.Code == http.StatusNotModified
}

// NotFound returns true if the error is an HTTP 404
func (e HTTPError) NotFound() bool {
	return e.Code == http.StatusNotFound
}
