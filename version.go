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
