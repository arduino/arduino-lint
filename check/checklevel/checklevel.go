package checklevel

import (
	"fmt"

	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
)

//go:generate stringer -type=Type -linecomment
type Type int

// Line comments set the string for each level
const (
	Info    Type = iota // info
	Warning             // warning
	Error               // error
	Notice              // notice
)

func CheckLevel(checkConfiguration checkconfigurations.Type) (Type, error) {
	configurationCheckModes := configuration.CheckModes(checkConfiguration.ProjectType)
	for _, errorMode := range checkConfiguration.ErrorModes {
		if configurationCheckModes[errorMode] == true {
			return Error, nil
		}
	}

	for _, warningMode := range checkConfiguration.WarningModes {
		if configurationCheckModes[warningMode] == true {
			return Warning, nil
		}
	}

	for _, infoMode := range checkConfiguration.InfoModes {
		if configurationCheckModes[infoMode] == true {
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
