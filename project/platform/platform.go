/*
Package packageindex provides functions specific to checking Arduino boards platforms.
See: https://arduino.github.io/arduino-cli/latest/platform-specification/
*/
package platform

import (
	"github.com/arduino/go-paths-helper"
)

var empty struct{}

var configurationFilenames = map[string]struct{}{
	"boards.txt":         empty,
	"boards.local.txt":   empty,
	"platform.txt":       empty,
	"platform.local.txt": empty,
	"programmers.txt":    empty,
}

// IsConfigurationFile returns whether the file at the given path has a boards platform configuration file filename.
func IsConfigurationFile(filePath *paths.Path) bool {
	_, isConfigurationFile := configurationFilenames[filePath.Base()]
	if isConfigurationFile {
		return true
	}
	return false
}

var requiredConfigurationFilenames = map[string]struct{}{
	// Arduino platforms must always contain a boards.txt
	"boards.txt": empty,
}

// IsRequiredConfigurationFile returns whether the file at the given path has the filename of a required boards platform configuration file.
func IsRequiredConfigurationFile(filePath *paths.Path) bool {
	_, isRequiredConfigurationFile := requiredConfigurationFilenames[filePath.Base()]
	if isRequiredConfigurationFile {
		return true
	}
	return false
}

// See: https://arduino.github.io/arduino-cli/latest/platform-specification/#platform-bundled-libraries
var bundledLibrariesFolderNames = map[string]struct{}{
	"libraries": empty,
}

func BundledLibrariesFolderNames() []string {
	folderNames := make([]string, 0, len(bundledLibrariesFolderNames))
	for folderName := range bundledLibrariesFolderNames {
		folderNames = append(folderNames, folderName)
	}

	return folderNames
}
