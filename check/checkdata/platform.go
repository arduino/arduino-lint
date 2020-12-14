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

package checkdata

import (
	"github.com/arduino/arduino-lint/project"
	"github.com/arduino/arduino-lint/project/platform/boardstxt"
	"github.com/arduino/go-properties-orderedmap"
)

// Initialize gathers the platform check data for the specified project.
func InitializeForPlatform(project project.Type) {
	boardsTxt, boardsTxtLoadError = boardstxt.Properties(ProjectPath())
}

var boardsTxt *properties.Map

// BoardsTxt returns the data from the boards.txt configuration file.
func BoardsTxt() *properties.Map {
	return boardsTxt
}

var boardsTxtLoadError error

// BoardsTxtLoadError returns the error output from loading the boards.txt configuration file.
func BoardsTxtLoadError() error {
	return boardsTxtLoadError
}
