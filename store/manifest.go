package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// OpenManifest either initializes a new leveldb database at the path of the string provided
// or opens an existing leveldb database found at that path.
func OpenManifest(path string) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic("Could not open manifest")
	}

	return db, err
}

/*
Author:  Rebecca Bilbro
Author:  Benjamin Bengfort
Created: Thu Oct 31 14:02:41 EDT 2019

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: manifest.go [] bilbro@gmail.com $
*/
