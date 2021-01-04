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

package projectdata

import (
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/platform/boardstxt"
	"github.com/arduino/arduino-lint/internal/rule/schema"
	"github.com/arduino/arduino-lint/internal/rule/schema/compliancelevel"
	"github.com/arduino/go-properties-orderedmap"
	"github.com/sirupsen/logrus"
)

// InitializeForPlatform gathers the platform rule data for the specified project.
func InitializeForPlatform(project project.Type) {
	boardsTxt, boardsTxtLoadError = boardstxt.Properties(ProjectPath())
	if boardsTxtLoadError != nil {
		logrus.Errorf("Error loading boards.txt from %s: %s", project.Path, boardsTxtLoadError)
		boardsTxtSchemaValidationResult = nil
	} else {
		boardsTxtSchemaValidationResult = boardstxt.Validate(boardsTxt)

		boardsTxtMenuIds = boardstxt.MenuIDs(boardsTxt)
		boardsTxtBoardIds = boardstxt.BoardIDs(boardsTxt)
	}
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

var boardsTxtSchemaValidationResult map[compliancelevel.Type]schema.ValidationResult

// BoardsTxtSchemaValidationResult returns the result of validating boards.txt against the JSON schema.
func BoardsTxtSchemaValidationResult() map[compliancelevel.Type]schema.ValidationResult {
	return boardsTxtSchemaValidationResult
}

var boardsTxtMenuIds []string

// BoardsTxtMenuIds returns the list of menu IDs present in the platform's boards.txt.
func BoardsTxtMenuIds() []string {
	return boardsTxtMenuIds
}

var boardsTxtBoardIds []string

// BoardsTxtMenuIds returns the list of board IDs present in the platform's boards.txt.
func BoardsTxtBoardIds() []string {
	return boardsTxtBoardIds
}
