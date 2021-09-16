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

package projectdata

import (
	"github.com/arduino/arduino-lint/internal/project"
	"github.com/arduino/arduino-lint/internal/project/platform/boardstxt"
	"github.com/arduino/arduino-lint/internal/project/platform/platformtxt"
	"github.com/arduino/arduino-lint/internal/project/platform/programmerstxt"
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
		boardsTxtVisibleBoardIds = boardstxt.VisibleBoardIDs(boardsTxt)
	}

	programmersTxtExists = ProjectPath().Join("programmers.txt").Exist()

	programmersTxt, programmersTxtLoadError = programmerstxt.Properties(ProjectPath())
	if programmersTxtLoadError != nil {
		logrus.Tracef("Error loading programmers.txt from %s: %s", project.Path, programmersTxtLoadError)
		programmersTxtSchemaValidationResult = nil
	} else {
		programmersTxtSchemaValidationResult = programmerstxt.Validate(programmersTxt)

		programmersTxtProgrammerIds = programmerstxt.ProgrammerIDs(programmersTxt)
	}

	platformTxtExists = ProjectPath().Join("platform.txt").Exist()

	platformTxt, platformTxtLoadError = platformtxt.Properties(ProjectPath())
	if platformTxtLoadError != nil {
		logrus.Tracef("Error loading platform.txt from %s: %s", project.Path, platformTxtLoadError)
		platformTxtSchemaValidationResult = nil
		platformTxtPluggableDiscoveryNames = nil
		platformTxtUserProvidedFieldNames = nil
		platformTxtToolNames = nil
	} else {
		platformTxtSchemaValidationResult = platformtxt.Validate(platformTxt)

		platformTxtPluggableDiscoveryNames = platformtxt.PluggableDiscoveryNames(platformTxt)
		platformTxtUserProvidedFieldNames = platformtxt.UserProvidedFieldNames(platformTxt)
		platformTxtToolNames = platformtxt.ToolNames(platformTxt)
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

// BoardsTxtBoardIds returns the list of board IDs present in the platform's boards.txt.
func BoardsTxtBoardIds() []string {
	return boardsTxtBoardIds
}

var boardsTxtVisibleBoardIds []string

// BoardsTxtVisibleBoardIds returns the list of IDs for visible boards present in the platform's boards.txt.
func BoardsTxtVisibleBoardIds() []string {
	return boardsTxtVisibleBoardIds
}

var programmersTxtExists bool

// ProgrammersTxtExists returns whether the platform contains a programmer.txt file.
func ProgrammersTxtExists() bool {
	return programmersTxtExists
}

var programmersTxt *properties.Map

// ProgrammersTxt returns the data from the programmers.txt configuration file.
func ProgrammersTxt() *properties.Map {
	return programmersTxt
}

var programmersTxtLoadError error

// ProgrammersTxtLoadError returns the error output from loading the programmers.txt configuration file.
func ProgrammersTxtLoadError() error {
	return programmersTxtLoadError
}

var programmersTxtSchemaValidationResult map[compliancelevel.Type]schema.ValidationResult

// ProgrammersTxtSchemaValidationResult returns the result of validating programmers.txt against the JSON schema.
func ProgrammersTxtSchemaValidationResult() map[compliancelevel.Type]schema.ValidationResult {
	return programmersTxtSchemaValidationResult
}

var programmersTxtProgrammerIds []string

// ProgrammersTxtProgrammerIds returns the list of board IDs present in the platform's programmers.txt.
func ProgrammersTxtProgrammerIds() []string {
	return programmersTxtProgrammerIds
}

var platformTxtExists bool

// PlatformTxtExists returns whether the platform contains a platform.txt file.
func PlatformTxtExists() bool {
	return platformTxtExists
}

var platformTxt *properties.Map

// PlatformTxt returns the data from the platform.txt configuration file.
func PlatformTxt() *properties.Map {
	return platformTxt
}

var platformTxtLoadError error

// PlatformTxtLoadError returns the error output from loading the platform.txt configuration file.
func PlatformTxtLoadError() error {
	return platformTxtLoadError
}

var platformTxtSchemaValidationResult map[compliancelevel.Type]schema.ValidationResult

// PlatformTxtSchemaValidationResult returns the result of validating platform.txt against the JSON schema.
func PlatformTxtSchemaValidationResult() map[compliancelevel.Type]schema.ValidationResult {
	return platformTxtSchemaValidationResult
}

var platformTxtPluggableDiscoveryNames []string

// PlatformTxtPluggableDiscoveryNames returns the list of pluggable discoveries present in the platform's platform.txt.
func PlatformTxtPluggableDiscoveryNames() []string {
	return platformTxtPluggableDiscoveryNames
}

var platformTxtUserProvidedFieldNames map[string][]string

// PlatformTxtUserProvidedFieldNames returns the list of user provided field names present in the platform's platform.txt, mapped by board name.
func PlatformTxtUserProvidedFieldNames() map[string][]string {
	return platformTxtUserProvidedFieldNames
}

var platformTxtToolNames []string

// PlatformTxtToolNames returns the list of tools present in the platform's platform.txt.
func PlatformTxtToolNames() []string {
	return platformTxtToolNames
}
