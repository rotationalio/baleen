/*
Package baleen_test provides testing for the functions in the baleen package.

Author:  Benjamin Bengfort
Author:  Rebecca Bilbro
Created: Thu Apr 25 18:32:19 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: version_test.go [68a2562] benjamin@bengfort.com $
*/
package baleen_test

import (
	"testing"

	. "github.com/kansaslabs/baleen"
)

const expectedVersion = "0.0"

func TestVersion(t *testing.T) {
	equals(t, expectedVersion, Version(false))
}
