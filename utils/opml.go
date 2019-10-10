package utils

import (
	"encoding/xml"
)

// An OPML struct can be loaded from an OPML file
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// The Head of the opml
type Head struct {
	XMLName xml.Name `xml:"head"`
	Title   string   `xml:"title"`
}

// The Body of the opml
type Body struct {
	XMLName  xml.Name  `xml:"body"`
	Outlines []Outline `xml:"outline"`
}

// The Outline of the opml
type Outline struct {
	Text    string `xml:"text,attr"`
	Title   string `xml:"title,attr"`
	Type    string `xml:"type,attr"`
	XMLURL  string `xml:"xmlUrl,attr"`
	HTMLURL string `xml:"htmlUrl,attr"`
	Favicon string `xml:"rssfr-favicon,attr"`
}

/*
Author:  Rebecca Bilbro
Author:  Benjamin Bengfort
Created: Tue Oct 8 13:47:54 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: opml.go [c23ae33] rebeccabilbro@users.noreply.github.com $
*/
