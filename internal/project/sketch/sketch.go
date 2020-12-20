// This file is part of arduino-lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

/*
Package sketch provides functions specific to checking Arduino sketches.
See: https://arduino.github.io/arduino-cli/latest/sketch-specification/
*/
package sketch

import (
	"fmt"

	"github.com/arduino/arduino-cli/arduino/globals"
	"github.com/arduino/go-paths-helper"
)

// HasMainFileValidExtension returns whether the file at the given path has a valid sketch main file extension.
// Sketches may contain source files with other extensions (e.g., .h, .cpp), but they are required to have at least one file with a main extension.
func HasMainFileValidExtension(filePath *paths.Path) bool {
	_, hasMainFileValidExtension := globals.MainFileValidExtensions[filePath.Ext()]
	return hasMainFileValidExtension
}

// ContainsMainSketchFile checks whether the provided path contains a file with valid main sketch file extension.
func ContainsMainSketchFile(searchPath *paths.Path) bool {
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
		if HasMainFileValidExtension(potentialHeaderFile) {
			return true
		}
	}

	return false
}

// HasSupportedExtension returns whether the file at the given path has any of the file extensions supported for source/header files of a sketch.
func HasSupportedExtension(filePath *paths.Path) bool {
	_, hasAdditionalFileValidExtensions := globals.AdditionalFileValidExtensions[filePath.Ext()]
	return hasAdditionalFileValidExtensions || HasMainFileValidExtension(filePath)
}

var empty struct{}

// See: https://arduino.github.io/arduino-cli/latest/sketch-specification/#metadata
var metadataFilenames = map[string]struct{}{
	"sketch.json": empty,
}

// MetadataPath returns the path of the sketch's metadata file.
func MetadataPath(sketchPath *paths.Path) *paths.Path {
	for metadataFileName := range metadataFilenames {
		metadataPath := sketchPath.Join(metadataFileName)
		exist, err := metadataPath.ExistCheck()
		if err != nil {
			panic(err)
		}

		if exist {
			return metadataPath
		}
	}

	return nil
}
