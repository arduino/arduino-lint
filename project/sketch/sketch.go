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

/*
Package sketch provides functions specific to checking Arduino sketches.
See: https://arduino.github.io/arduino-cli/latest/sketch-specification/
*/
package sketch

import (
	"github.com/arduino/arduino-cli/arduino/globals"
	"github.com/arduino/go-paths-helper"
)

// HasMainFileValidExtension returns whether the file at the given path has a valid sketch main file extension.
// Sketches may contain source files with other extensions (e.g., .h, .cpp), but they are required to have at least one file with a main extension.
func HasMainFileValidExtension(filePath *paths.Path) bool {
	_, hasMainFileValidExtension := globals.MainFileValidExtensions[filePath.Ext()]
	return hasMainFileValidExtension
}

// HasSupportedExtension returns whether the file at the given path has any of the file extensions supported for source/header files of a sketch.
func HasSupportedExtension(filePath *paths.Path) bool {
	_, hasAdditionalFileValidExtensions := globals.AdditionalFileValidExtensions[filePath.Ext()]
	return hasAdditionalFileValidExtensions || HasMainFileValidExtension(filePath)
}
