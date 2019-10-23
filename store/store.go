/*
Package store contains functions for writing parsed feeds to cloud storage.
*/
package store

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spaolacci/murmur3"
)

// AWSCredentials stores the region and bucket needed to access an S3 bucket.
type AWSCredentials struct {
	Region string
	Bucket string
}

// A Document is the generic representation of an individual entry in a feed
type Document struct {
	FeedID       string
	LanguageCode string
	Year         int
	Month        string
	Day          int
	Title        string
	Description  string
	Content      string
	Link         string
}

// VerifyCredentials is a helper function that verifies the credentials are correct and
// the bucket exists, returning true if so else false.
func VerifyCredentials(creds *AWSCredentials) bool {
	if creds != nil {
		return true
	}
	return false
}

// GetSession is a helper function that returns an AWS session that can be reused for
// multiple writes. It takes in a pointer to an AWSCredentials struct and returns a pointer
// to an AWS connection and a nil, else a nil and the connection error.
func GetSession(creds *AWSCredentials) (sesh *session.Session, err error) {
	if VerifyCredentials(creds) {
		sesh, err = session.NewSession(&aws.Config{Region: &creds.Region})
		if err != nil {
			return nil, err
		}
		return sesh, nil
	}
	return nil, errors.New("invalid credentials")
}

// Upload a single Document to the specified S3 bucket using the established AWS session,
// setting file information including name (the feedID, language code, year, month, day, and
// hash of the content), the content size and type, and the encryption on the uploaded file.
func Upload(s *session.Session, doc Document, bucket string) error {
	// Get doc size and read the content into a buffer
	size := len(doc.Content)
	buffer := make([]byte, size)

	// Hash the contents of the file and use the hash to create a unique filename
	hasher := murmur3.New64()
	hasher.Write(buffer)
	hash := strconv.FormatInt(int64(hasher.Sum64()), 10)
	name := doc.LanguageCode + "/" + strconv.Itoa(doc.Year) + "/" + doc.Month + "/" + strconv.Itoa(doc.Day) + "/" + doc.FeedID + "/" + hash + ".html"

	fmt.Printf("\n storing %s to s3", name)

	// Put the object to the S3 bucket
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(name),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader([]byte(doc.Content)),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

/*
Author:  Rebecca Bilbro
Author:  Benjamin Bengfort
Created: Tue Oct 8 14:02:56 EDT 2019

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: store.go [] bilbro@gmail.com $
*/
