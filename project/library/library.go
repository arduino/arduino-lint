// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

// Package library provides functions specific to checking Arduino libraries.
package library

import (
	"fmt"

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
	return hasHeaderFileValidExtension
}

// ContainsHeaderFile checks whether the provided path contains a file with valid header extension.
func ContainsHeaderFile(searchPath *paths.Path) bool {
	if searchPath.NotExist() {
		panic(fmt.Sprintf("Error: provided path %s does not exist.", searchPath))
	}
	if searchPath.IsNotDir() {
		panic(fmt.Sprintf("Error: provided path %s is not a directory.", searchPath))
	}

	directoryListing, err := searchPath.ReadDir()
	if err != nil {
		panic(err)
	}

	directoryListing.FilterOutDirs()
	for _, potentialHeaderFile := range directoryListing {
		if HasHeaderFileValidExtension(potentialHeaderFile) {
			return true
		}
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

// ContainsMetadataFile checks whether the provided path contains an Arduino library metadata file.
func ContainsMetadataFile(searchPath *paths.Path) bool {
	if searchPath.NotExist() {
		panic(fmt.Sprintf("Error: provided path %s does not exist.", searchPath))
	}
	if searchPath.IsNotDir() {
		panic(fmt.Sprintf("Error: provided path %s is not a directory.", searchPath))
	}

	directoryListing, err := searchPath.ReadDir()
	if err != nil {
		panic(err)
	}

	directoryListing.FilterOutDirs()
	for _, potentialMetadataFile := range directoryListing {
		if IsMetadataFile(potentialMetadataFile) {
			return true
		}
	}

	return false
}

// See: https://arduino.github.io/arduino-cli/latest/library-specification/#library-examples
var examplesFolderValidNames = map[string]struct{}{
	"examples": empty,
}

// Only "examples" is specification-compliant, but apparently "example" is also supported
// See: https://github.com/arduino/arduino-cli/blob/0.13.0/arduino/libraries/loader.go#L153
var examplesFolderSupportedNames = map[string]struct{}{
	"examples": empty,
	"example":  empty,
}

// ExamplesFolderSupportedNames returns a slice of supported examples folder names
func ExamplesFolderSupportedNames() []string {
	folderNames := make([]string, 0, len(examplesFolderSupportedNames))
	for folderName := range examplesFolderSupportedNames {
		folderNames = append(folderNames, folderName)
	}

	return folderNames
}
