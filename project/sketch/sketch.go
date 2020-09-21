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
	if hasMainFileValidExtension {
		return true
	}
	return false
}
