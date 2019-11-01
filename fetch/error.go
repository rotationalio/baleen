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

/*
Author:  Benjamin Bengfort
Author:  Rebecca Bilbro
Created: Mon Apr 29 06:43:36 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: error.go [d6dba70] benjamin@bengfort.com $
*/
