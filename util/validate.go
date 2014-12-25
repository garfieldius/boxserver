package util

import (
	"regexp"
)

var validIdentifier *regexp.Regexp
var validVersion *regexp.Regexp

func init() {
	validIdentifier = regexp.MustCompile(`^[a-z0-9][a-z0-9_\-\.]+[a-z0-9]$`)
	validVersion    = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
}

func ValidKey(key string) bool {
	return validIdentifier.Match(([]byte)(key))
}

func ValidVersion(version string) bool {
	return validVersion.Match(([]byte)(version))
}
