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

// Package configuration handles the configuration of the arduino-check tool.
package configuration

import (
	"os"

	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/result/outputformat"
	"github.com/arduino/go-paths-helper"
	"github.com/sirupsen/logrus"
)

// TODO: will it be possible to use init() instead?
// Initialize sets up the tool configuration according to defaults and user-specified options.
func Initialize() {
	setDefaults()
	// TODO configuration according to command line input
	// TODO validate target path value, exit if not found
	// TODO support multiple paths
	// TODO validate output format input

	targetPath = paths.New("e:/electronics/arduino/libraries/arduino-check-test-library")

	// customCheckModes[checkmode.Permissive] = false
	// customCheckModes[checkmode.LibraryManagerSubmission] = false
	// customCheckModes[checkmode.LibraryManagerIndexed] = false
	// customCheckModes[checkmode.Official] = false
	// superprojectType = projecttype.All

	outputFormat = outputformat.JSON
	//reportFilePath = paths.New("report.json")

	logrus.SetLevel(logrus.PanicLevel)

	logrus.WithFields(logrus.Fields{
		"superproject type filter": SuperprojectTypeFilter(),
		"recursive":                Recursive(),
		"projects path":            TargetPath(),
	}).Debug("Configuration initialized")
}

var customCheckModes = make(map[checkmode.Type]bool)

// CheckModes returns the check modes configuration for the given project type.
func CheckModes(superprojectType projecttype.Type) map[checkmode.Type]bool {
	return checkmode.Modes(defaultCheckModes, customCheckModes, superprojectType)
}

var superprojectTypeFilter projecttype.Type

// SuperprojectType returns the superproject type filter configuration.
func SuperprojectTypeFilter() projecttype.Type {
	return superprojectTypeFilter
}

var recursive bool

// Recursive returns the recursive project search configuration value.
func Recursive() bool {
	return recursive
}

var outputFormat outputformat.Type

// OutputFormat returns the tool output format configuration value.
func OutputFormat() outputformat.Type {
	return outputFormat
}

var reportFilePath *paths.Path

// ReportFilePath returns the path to save the report file at.
func ReportFilePath() *paths.Path {
	return reportFilePath
}

var targetPath *paths.Path

// TargetPath returns the projects search path.
func TargetPath() *paths.Path {
	return targetPath
}

// SchemasPath returns the path to the folder containing the JSON schemas.
func SchemasPath() *paths.Path {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return paths.New(executablePath).Parent().Join("etc", "schemas")
}
