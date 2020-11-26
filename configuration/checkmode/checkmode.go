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

// Package checkmode defines the tool configuration options that affect checks.
package checkmode

import (
	"fmt"
	"strings"

	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/sirupsen/logrus"
)

// Type is the type for check modes.
type Type int

//go:generate stringer -type=Type -linecomment
const (
	Strict                   Type = iota // strict
	Specification                        // specification
	Permissive                           // permissive
	LibraryManagerSubmission             // --library-manager=submit
	LibraryManagerIndexed                // --library-manager=update
	Official                             // ARDUINO_CHECK_OFFICIAL
	All                                  // always
	Default                              // default
)

// ComplianceModeFromString parses the --compliance flag value and returns the corresponding check mode settings.
func ComplianceModeFromString(complianceModeString string) (bool, bool, bool, error) {
	switch strings.ToLower(complianceModeString) {
	case "strict":
		return true, false, false, nil
	case "specification":
		return false, true, false, nil
	case "permissive":
		return false, false, true, nil
	default:
		return false, false, false, fmt.Errorf("No matching compliance mode for string %s", complianceModeString)
	}
}

// LibraryManagerModeFromString parses the --library-manager flag value and returns the corresponding check mode settings.
func LibraryManagerModeFromString(libraryManagerModeString string) (bool, bool, error) {
	switch strings.ToLower(libraryManagerModeString) {
	case "submit":
		return true, false, nil
	case "update":
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

	// This mode is always enabled
	checkModes[All] = true

	return checkModes
}
