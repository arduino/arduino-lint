// Package checklevel defines the level assigned to a check failure.
package checklevel

import (
	"fmt"

	"github.com/arduino/arduino-check/check/checkconfigurations"
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

// CheckLevel determines the check level assigned to failure of the given check under the current tool configuration.
func CheckLevel(checkConfiguration checkconfigurations.Type) (Type, error) {
	configurationCheckModes := configuration.CheckModes(checkConfiguration.ProjectType)
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
