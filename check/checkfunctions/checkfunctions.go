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

// Package checkfunctions contains the functions that implement each check.
package checkfunctions

import (
	"regexp"

	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/go-paths-helper"
)

// Type is the function signature for the check functions.
// The `output` result is the contextual information that will be inserted into the check's message template.
type Type func() (result checkresult.Type, output string)

// validProjectPathBaseName checks whether the provided library folder or sketch filename contains prohibited characters.
func validProjectPathBaseName(name string) bool {
	baseNameRegexp := regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")
	return baseNameRegexp.MatchString(name)
}

func containsMisspelledPathBaseName(pathList paths.PathList, correctBaseName string, misspellingQuery string) (*paths.Path, bool) {
	misspellingRegexp := regexp.MustCompile(misspellingQuery)
	for _, path := range pathList {
		if path.Base() == correctBaseName {
			return nil, false
		}

		if misspellingRegexp.MatchString(path.Base()) {
			return path, true
		}
	}

	return nil, false
}
