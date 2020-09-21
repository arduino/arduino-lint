// Package library provides functions specific to checking Arduino libraries.
package library

import (
	"github.com/arduino/go-paths-helper"
)

var empty struct{}

// Reference: https://github.com/arduino/arduino-cli/blob/0.13.0/arduino/libraries/libraries.go#L167
var headerFileValidExtensions = map[string]struct{}{
	".h":   empty,
	".hpp": empty,
	".hh":  empty,
}

// HasHeaderFileValidExtension returns whether the file at the given path has a valid library header file extension.
func HasHeaderFileValidExtension(filePath *paths.Path) bool {
	_, hasHeaderFileValidExtension := headerFileValidExtensions[filePath.Ext()]
	if hasHeaderFileValidExtension {
		return true
	}
	return false
}

// See: https://arduino.github.io/arduino-cli/latest/library-specification/#library-metadata
var metadataFilenames = map[string]struct{}{
	"library.properties": empty,
}

// IsMetadataFile returns whether the file at the given path is an Arduino library metadata file.
func IsMetadataFile(filePath *paths.Path) bool {
	_, isMetadataFile := metadataFilenames[filePath.Base()]
	if isMetadataFile {
		return true
	}
	return false
}
