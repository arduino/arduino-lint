// This file is part of Arduino Lint.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License, either
// version 3 of the License, or (at your option) any later version.
// This license covers the main part of Arduino Lint.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

/*
Package platform provides functions specific to linting Arduino boards platforms.
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
	return isConfigurationFile
}

var requiredConfigurationFilenames = map[string]struct{}{
	// Arduino platforms must always contain a boards.txt
	"boards.txt": empty,
}

// IsRequiredConfigurationFile returns whether the file at the given path has the filename of a required boards platform configuration file.
func IsRequiredConfigurationFile(filePath *paths.Path) bool {
	_, isRequiredConfigurationFile := requiredConfigurationFilenames[filePath.Base()]
	return isRequiredConfigurationFile
}

// See: https://arduino.github.io/arduino-cli/latest/platform-specification/#platform-bundled-libraries
var bundledLibrariesFolderNames = map[string]struct{}{
	"libraries": empty,
}

// BundledLibrariesFolderNames returns a list of supported names for the platform bundled libraries folder.
func BundledLibrariesFolderNames() []string {
	folderNames := make([]string, 0, len(bundledLibrariesFolderNames))
	for folderName := range bundledLibrariesFolderNames {
		folderNames = append(folderNames, folderName)
	}

	return folderNames
}
