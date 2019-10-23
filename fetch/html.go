package fetch

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/k3a/html2text"
)

// GetContent retrieves the full content of articles using the url string provided in
// the RSS feed and returns a string representation of the content found on that page
func GetContent(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("no text retrieved from %s", url)
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("no text parsed from html retrieved from %s", url)
	}
	content := html2text.HTML2Text(string(html))

	return content
}

/*
Author:  Rebecca Bilbro
Author:  Benjamin Bengfort
Created: Sun Oct 20 15:26:21 EDT 2019

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: html.go [] bilbro@gmail.com $
*/
