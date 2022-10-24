/*
Package opml implements parsing support for Outline Processor Markup Language - an XML
format for creating outlines. It has since been adopted for other uses, the most common
being to exchange lists of web feeds between web feed aggregators. Baleen uses this
format to read a list of web feeds from a file on disk and to exchange lists of feeds
between feed aggregation tools and sites.
*/
package opml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// An OPML struct can be loaded from an OPML file using an XML processor.
type OPML struct {
	XMLName xml.Name `xml:"opml" json:"-"`
	Version string   `xml:"version,attr" json:"_version"`
	Head    Head     `xml:"head" json:"head"`
	Body    Body     `xml:"body" json:"body"`
}

// The Head of the OPML data.
type Head struct {
	XMLName xml.Name `xml:"head" json:"-"`
	Title   string   `xml:"title" json:"title"`
}

// The Body of the OPML data.
type Body struct {
	XMLName  xml.Name  `xml:"body" json:"-"`
	Outlines []Outline `xml:"outline" json:"outline"`
}

// The Outline of the OPML data - this contains the primary content for the web feed.
type Outline struct {
	Text    string `xml:"text,attr" json:"_text"`
	Title   string `xml:"title,attr" json:"_title"`
	Type    string `xml:"type,attr" json:"_type"`
	XMLURL  string `xml:"xmlUrl,attr" json:"_xmlUrl"`
	HTMLURL string `xml:"htmlUrl,attr" json:"_htmlUrl"`
	Favicon string `xml:"rssfr-favicon,attr" json:"-"`
}

// Load an outline from an XML file stored on disk.
func Load(path string) (outline *OPML, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	outline = &OPML{}
	switch ext := filepath.Ext(path); ext {
	case ".opml", ".xml":
		if err = xml.NewDecoder(f).Decode(outline); err != nil {
			return nil, err
		}
	case ".json":
		wrapper := make(map[string]*OPML)
		if err = json.NewDecoder(f).Decode(&wrapper); err != nil {
			return nil, err
		}

		var ok bool
		if outline, ok = wrapper["opml"]; !ok {
			return nil, errors.New("could not find opml json data in object")
		}
	default:
		return nil, fmt.Errorf("unknown extension %q", ext)
	}

	return outline, nil
}

// URLs implements Baleen's most common use for OPML: extracting all of the feed URLs
// from the outline documents and returning it as a slice of URL strings. This method
// can optionally filter by type. E.g. to specify only RSS use `o.URLs("rss")`, if no
// types are specified, all types are returned.
//
// When processing the outline URLs, this method will take the XMLURL first, and if it
// is empty it will return the HTMLURL. If neither URL contains data it will be skipped.
func (o *OPML) URLs(types ...string) []string {
	var filter map[string]struct{}
	if len(types) > 0 {
		filter = make(map[string]struct{})
		for _, t := range types {
			filter[t] = struct{}{}
		}
	}

	urls := make([]string, 0, len(o.Body.Outlines))
	for _, outline := range o.Body.Outlines {
		// Filter the outline by the specified type
		if filter != nil {
			if _, ok := filter[outline.Type]; !ok {
				continue
			}
		}

		switch {
		case outline.XMLURL != "":
			urls = append(urls, outline.XMLURL)
		case outline.HTMLURL != "":
			urls = append(urls, outline.HTMLURL)
		}
	}

	return urls
}
