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

// Package rulemode defines the tool configuration options that affect rules.
package rulemode

import (
	"fmt"
	"strings"

	"github.com/arduino/arduino-lint/internal/project/projecttype"
	"github.com/sirupsen/logrus"
)

// Type is the type for rule modes.
type Type int

//go:generate go tool golang.org/x/tools/cmd/stringer -type=Type -linecomment
const (
	Strict                   Type = iota // strict
	Specification                        // specification
	Permissive                           // permissive
	LibraryManagerSubmission             // submit
	LibraryManagerIndexed                // update
	LibraryManagerIndexing               // ARDUINO_LINT_LIBRARY_MANAGER_INDEXING
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
	LibraryManagerIndexing:   empty,
	Official:                 empty,
	Default:                  empty,
}

// ComplianceModeFromString parses the --compliance flag value and returns the corresponding rule mode settings.
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

// LibraryManagerModeFromString parses the --library-manager flag value and returns the corresponding rule mode settings.
func LibraryManagerModeFromString(libraryManagerModeString string) (bool, bool, bool, error) {
	switch strings.ToLower(libraryManagerModeString) {
	case LibraryManagerSubmission.String():
		return true, false, false, nil
	case LibraryManagerIndexed.String():
		return false, true, false, nil
	case "false":
		return false, false, false, nil
	default:
		return false, false, false, fmt.Errorf("No matching Library Manager mode for string %s", libraryManagerModeString)
	}
}

// Modes merges the default rule mode values for the given superproject type with any user-specified rule mode settings.
func Modes(defaultRuleModes map[projecttype.Type]map[Type]bool, customRuleModes map[Type]bool, superprojectType projecttype.Type) map[Type]bool {
	ruleModes := make(map[Type]bool)

	for key, defaultValue := range defaultRuleModes[superprojectType] {
		customRuleModeValue, customRuleModeIsConfigured := customRuleModes[key]
		if customRuleModeIsConfigured {
			ruleModes[key] = customRuleModeValue
		} else {
			ruleModes[key] = defaultValue
		}
		logrus.Tracef("Rule mode option %s set to %t\n", key, ruleModes[key])
	}

	return ruleModes
}

// Compliance returns the tool configuration's compliance setting name.
func Compliance(ruleModes map[Type]bool) string {
	for key, value := range ruleModes {
		if value && (key == Strict || key == Specification || key == Permissive) {
			return key.String()
		}
	}

	panic(fmt.Errorf("Unrecognized compliance configuration"))
}

// LibraryManager returns the string identifier for the Library Manager configuration mode.
func LibraryManager(ruleModes map[Type]bool) string {
	for key, value := range ruleModes {
		if value && (key == LibraryManagerSubmission || key == LibraryManagerIndexed || key == LibraryManagerIndexing) {
			return key.String()
		}
	}

	return "false"
}
