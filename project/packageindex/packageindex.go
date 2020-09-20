package packageindex

import (
	"regexp"

	"github.com/arduino/go-paths-helper"
)

var empty struct{}

var validExtensions = map[string]struct{}{
	".json": empty,
}

func HasValidExtension(filePath *paths.Path) bool {
	_, hasValidExtension := validExtensions[filePath.Ext()]
	if hasValidExtension {
		return true
	}
	return false
}

// Regular expressions for official and non-official package index filenames
// See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
var validFilenameRegex = map[bool]*regexp.Regexp{
	true:  regexp.MustCompile(`^package_(.+_)*index.json$`),
	false: regexp.MustCompile(`^package_(.+_)+index.json$`),
}

func HasValidFilename(filePath *paths.Path, officialCheckMode bool) bool {
	regex := validFilenameRegex[officialCheckMode]
	filename := filePath.Base()
	return regex.MatchString(filename)
}
