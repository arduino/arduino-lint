/*
Package packageindex provides functions specific to checking the package index files of the Arduino Boards Manager.
See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification
*/
package packageindex

import (
	"regexp"

	"github.com/arduino/go-paths-helper"
)

var empty struct{}

// Reference: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
var validExtensions = map[string]struct{}{
	".json": empty,
}

// HasValidExtension returns whether the file at the given path has a valid package index extension.
func HasValidExtension(filePath *paths.Path) bool {
	_, hasValidExtension := validExtensions[filePath.Ext()]
	return hasValidExtension
}

// Regular expressions for official and non-official package index filenames
// See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
var validFilenameRegex = map[bool]*regexp.Regexp{
	true:  regexp.MustCompile(`^package_(.+_)*index.json$`),
	false: regexp.MustCompile(`^package_(.+_)+index.json$`),
}

// HasValidFilename returns whether the file at the given path has a valid package index filename.
func HasValidFilename(filePath *paths.Path, officialCheckMode bool) bool {
	regex := validFilenameRegex[officialCheckMode]
	filename := filePath.Base()
	return regex.MatchString(filename)
}
