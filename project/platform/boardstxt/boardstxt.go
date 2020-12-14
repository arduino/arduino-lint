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
Package boardstxt provides functions specific to checking the boards.txt configuration files of Arduino boards platforms.
See: https://arduino.github.io/arduino-cli/latest/platform-specification/#boardstxt
*/
package boardstxt

import (
	"github.com/arduino/go-paths-helper"
	"github.com/arduino/go-properties-orderedmap"
)

// Properties parses the library.properties from the given path and returns the data.
func Properties(platformPath *paths.Path) (*properties.Map, error) {
	return properties.SafeLoadFromPath(platformPath.Join("boards.txt"))
}
