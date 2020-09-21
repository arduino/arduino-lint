// Package configuration handles the configuration of the arduino-check tool.
package configuration

import (
	"github.com/arduino/arduino-check/configuration/checkmode"
	"github.com/arduino/arduino-check/project/projecttype"
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
	targetPath = paths.New("e:/electronics/arduino/libraries/arduino-check-test-library")

	// customCheckModes[checkmode.Permissive] = false
	// customCheckModes[checkmode.LibraryManagerSubmission] = false
	// customCheckModes[checkmode.LibraryManagerIndexed] = false
	// customCheckModes[checkmode.Official] = false
	// superprojectType = projecttype.All
	logrus.SetLevel(logrus.PanicLevel)

	logrus.WithFields(logrus.Fields{
		"superproject type filter": SuperprojectTypeFilter().String(),
		"recursive":                Recursive(),
		"projects path":            TargetPath().String(),
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

var targetPath *paths.Path

// TargetPath returns the projects search path.
func TargetPath() *paths.Path {
	return targetPath
}
