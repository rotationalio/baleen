package baleen_test

import (
	"testing"

	. "github.com/kansaslabs/baleen"
)

const expectedVersion = "0.0"

func TestVersion(t *testing.T) {
	equals(t, expectedVersion, Version(false))
}
