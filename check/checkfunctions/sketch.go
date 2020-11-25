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

package checkfunctions

// The check functions for sketches.

import (
	"strings"

	"github.com/arduino/arduino-check/check/checkdata"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/project/sketch"
)

// ProhibitedCharactersInSketchFileName checks for prohibited characters in the sketch file names.
func ProhibitedCharactersInSketchFileName() (result checkresult.Type, output string) {
	directoryListing, _ := checkdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()

	foundInvalidSketchFileNames := []string{}
	for _, potentialSketchFile := range directoryListing {
		if sketch.HasSupportedExtension(potentialSketchFile) {
			if !validProjectPathBaseName(potentialSketchFile.Base()) {
				foundInvalidSketchFileNames = append(foundInvalidSketchFileNames, potentialSketchFile.Base())
			}
		}
	}

	if len(foundInvalidSketchFileNames) > 0 {
		return checkresult.Fail, strings.Join(foundInvalidSketchFileNames, ", ")
	}

	return checkresult.Pass, ""
}

// SketchFileNameGTMaxLength checks if the sketch file names exceed the maximum length.
func SketchFileNameGTMaxLength() (result checkresult.Type, output string) {
	directoryListing, _ := checkdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()

	foundTooLongSketchFileNames := []string{}
	for _, potentialSketchFile := range directoryListing {
		if sketch.HasSupportedExtension(potentialSketchFile) {
			if len(potentialSketchFile.Base())-len(potentialSketchFile.Ext()) > 63 {
				foundTooLongSketchFileNames = append(foundTooLongSketchFileNames, potentialSketchFile.Base())
			}
		}
	}

	if len(foundTooLongSketchFileNames) > 0 {
		return checkresult.Fail, strings.Join(foundTooLongSketchFileNames, ", ")
	}

	return checkresult.Pass, ""
}

// PdeSketchExtension checks for use of deprecated .pde sketch file extensions.
func PdeSketchExtension() (result checkresult.Type, output string) {
	directoryListing, _ := checkdata.ProjectPath().ReadDir()
	directoryListing.FilterOutDirs()
	pdeSketches := []string{}
	for _, filePath := range directoryListing {
		if filePath.Ext() == ".pde" {
			pdeSketches = append(pdeSketches, filePath.Base())
		}
	}

	if len(pdeSketches) > 0 {
		return checkresult.Fail, strings.Join(pdeSketches, ", ")
	}

	return checkresult.Pass, ""
}
