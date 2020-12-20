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

// Package checkmode defines the tool configuration options that affect checks.
package checkmode

import (
	"fmt"
	"strings"

	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/sirupsen/logrus"
)

// Type is the type for check modes.
type Type int

//go:generate stringer -type=Type -linecomment
const (
	Strict                   Type = iota // strict
	Specification                        // specification
	Permissive                           // permissive
	LibraryManagerSubmission             // submit
	LibraryManagerIndexed                // update
	Official                             // ARDUINO_LINT_OFFICIAL
	Default                              // default
)

var empty struct{}

// Types provides an iterator and validator for Type.
var Types = map[Type]struct{}{
	Strict:                   empty,
	Specification:            empty,
	Permissive:               empty,
	LibraryManagerSubmission: empty,
	LibraryManagerIndexed:    empty,
	Official:                 empty,
	Default:                  empty,
}

// ComplianceModeFromString parses the --compliance flag value and returns the corresponding check mode settings.
func ComplianceModeFromString(complianceModeString string) (bool, bool, bool, error) {
	switch strings.ToLower(complianceModeString) {
	case Strict.String():
		return true, false, false, nil
	case Specification.String():
		return false, true, false, nil
	case Permissive.String():
		return false, false, true, nil
	default:
		return false, false, false, fmt.Errorf("No matching compliance mode for string %s", complianceModeString)
	}
}

// LibraryManagerModeFromString parses the --library-manager flag value and returns the corresponding check mode settings.
func LibraryManagerModeFromString(libraryManagerModeString string) (bool, bool, error) {
	switch strings.ToLower(libraryManagerModeString) {
	case LibraryManagerSubmission.String():
		return true, false, nil
	case LibraryManagerIndexed.String():
		return false, true, nil
	case "false":
		return false, false, nil
	default:
		return false, false, fmt.Errorf("No matching Library Manager mode for string %s", libraryManagerModeString)
	}
}

// Modes merges the default check mode values for the given superproject type with any user-specified check mode settings.
func Modes(defaultCheckModes map[projecttype.Type]map[Type]bool, customCheckModes map[Type]bool, superprojectType projecttype.Type) map[Type]bool {
	checkModes := make(map[Type]bool)

	for key, defaultValue := range defaultCheckModes[superprojectType] {
		customCheckModeValue, customCheckModeIsConfigured := customCheckModes[key]
		if customCheckModeIsConfigured {
			checkModes[key] = customCheckModeValue
		} else {
			checkModes[key] = defaultValue
		}
		logrus.Tracef("Check mode option %s set to %t\n", key, checkModes[key])
	}

	return checkModes
}

func Compliance(checkModes map[Type]bool) string {
	for key, value := range checkModes {
		if value && (key == Strict || key == Specification || key == Permissive) {
			return key.String()
		}
	}

	panic(fmt.Errorf("Unrecognized compliance configuration"))
}

// LibraryManager returns the string identifier for the Library Manager configuration mode.
func LibraryManager(checkModes map[Type]bool) string {
	for key, value := range checkModes {
		if value && (key == LibraryManagerSubmission || key == LibraryManagerIndexed) {
			return key.String()
		}
	}

	return "false"
}
