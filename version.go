/*
Package baleen is the top level library of the baleen language ingestion service. This
library provides the primary components for running the service as a long running
background daemon including the main service itself, configuration and other utilities.
*/

package baleen

import "fmt"

// Version constants should be modified on new releases.
const (
	majorVersion = 0
	minorVersion = 0
	patchVersion = 0
	releaseLevel = "final"
	serialNumber = 0
)

// Version returns a string representation of the semantic package version number.
func Version(short bool) string {
	vers := fmt.Sprintf("%d.%d", majorVersion, minorVersion)
	if patchVersion > 0 {
		vers += fmt.Sprintf(".%d", minorVersion)
	}

	if !short {
		switch releaseLevel {
		case "alpha":
			vers += fmt.Sprintf("a%d", serialNumber)
		case "beta":
			vers += fmt.Sprintf("b%d", serialNumber)
		}
	}
	return vers
}

/*
Author:  Benjamin Bengfort
Author:  Rebecca Bilbro
Created: Thu Apr 25 18:32:19 2019 -0400

Copyright (C) 2019 Kansas Labs
For license information, see LICENSE.txt

ID: version.go [68a2562] benjamin@bengfort.com $
*/
