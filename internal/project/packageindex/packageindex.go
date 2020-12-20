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
Package packageindex provides functions specific to checking the package index files of the Arduino Boards Manager.
See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification
*/
package packageindex

import (
	"fmt"
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

func Find(folderPath *paths.Path) (*paths.Path, error) {
	exist, err := folderPath.ExistCheck()
	if !exist {
		return nil, fmt.Errorf("Error opening path %s: %s", folderPath, err)
	}

	if folderPath.IsNotDir() {
		return folderPath, nil
	}

	directoryListing, err := folderPath.ReadDir()
	if err != nil {
		return nil, err
	}

	directoryListing.FilterOutDirs()
	for _, potentialPackageIndexFile := range directoryListing {
		if HasValidFilename(potentialPackageIndexFile, true) {
			return potentialPackageIndexFile, nil
		}
	}
	for _, potentialPackageIndexFile := range directoryListing {
		if HasValidExtension(potentialPackageIndexFile) {
			return potentialPackageIndexFile, nil
		}
	}

	return nil, fmt.Errorf("No package index file found in %s", folderPath)
}
