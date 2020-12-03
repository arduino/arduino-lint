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

// Package checklevel defines the level assigned to a check failure.
package checklevel

import (
	"fmt"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/check/checkresult"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
)

// Type is the type for the check levels.
//go:generate stringer -type=Type -linecomment
type Type int

// The line comments set the string for each level.
const (
	Info    Type = iota // info
	Warning             // warning
	Error               // error
	Notice              // notice
)

// CheckLevel determines the check level assigned to the given result of the given check under the current tool configuration.
func CheckLevel(checkConfiguration checkconfigurations.Type, checkResult checkresult.Type) (Type, error) {
	if checkResult != checkresult.Fail {
		return Notice, nil // Level provided by FailCheckLevel() is only relevant for failure result.
	}
	configurationCheckModes := configuration.CheckModes(checkConfiguration.ProjectType)
	return FailCheckLevel(checkConfiguration, configurationCheckModes)
}

// FailCheckLevel determines the level of a failed check for the given check modes.
func FailCheckLevel(checkConfiguration checkconfigurations.Type, configurationCheckModes map[checkmode.Type]bool) (Type, error) {
	for _, errorMode := range checkConfiguration.ErrorModes {
		if configurationCheckModes[errorMode] {
			return Error, nil
		}
	}

	for _, warningMode := range checkConfiguration.WarningModes {
		if configurationCheckModes[warningMode] {
			return Warning, nil
		}
	}

	for _, infoMode := range checkConfiguration.InfoModes {
		if configurationCheckModes[infoMode] {
			return Info, nil
		}
	}

	// Use default level
	for _, errorMode := range checkConfiguration.ErrorModes {
		if errorMode == checkmode.Default {
			return Error, nil
		}
	}

	for _, warningMode := range checkConfiguration.WarningModes {
		if warningMode == checkmode.Default {
			return Warning, nil
		}
	}

	for _, infoMode := range checkConfiguration.InfoModes {
		if infoMode == checkmode.Default {
			return Info, nil
		}
	}

	return Notice, fmt.Errorf("Check %s is incorrectly configured", checkConfiguration.ID)
}
