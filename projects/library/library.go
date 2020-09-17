package library

import (
	"github.com/arduino/go-paths-helper"
)

var empty struct{}

// reference: https://github.com/arduino/arduino-cli/blob/0.13.0/arduino/libraries/libraries.go#L167
var headerFileValidExtensions = map[string]struct{}{
	".h":   empty,
	".hpp": empty,
	".hh":  empty,
}

func HasHeaderFileValidExtension(filePath *paths.Path) bool {
	_, hasHeaderFileValidExtension := headerFileValidExtensions[filePath.Ext()]
	if hasHeaderFileValidExtension {
		return true
	}
	return false
}

var metadataFilenames = map[string]struct{}{
	"library.properties": empty,
}

func IsMetadataFile(filePath *paths.Path) bool {
	_, isMetadataFile := metadataFilenames[filePath.Base()]
	if isMetadataFile {
		return true
	}
	return false
}
