package checklevel

import (
	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/configuration/checkmode"
)

//go:generate stringer -type=Type -linecomment
type Type int

// Line comments set the string for each level
const (
	Pass    Type = iota // pass
	Info                // info
	Warning             // warning
	Error               // error
	Notice              // notice
)

func CheckLevel(checkConfiguration checkconfigurations.Type) Type {
	configurationCheckModes := configuration.CheckModes(checkConfiguration.ProjectType)
	for _, errorMode := range checkConfiguration.ErrorModes {
		if configurationCheckModes[errorMode] == true {
			return Error
		}
	}

	for _, warningMode := range checkConfiguration.WarningModes {
		if configurationCheckModes[warningMode] == true {
			return Warning
		}
	}

	for _, infoMode := range checkConfiguration.InfoModes {
		if configurationCheckModes[infoMode] == true {
			return Info
		}
	}

	for _, passMode := range checkConfiguration.PassModes {
		if configurationCheckModes[passMode] == true {
			return Pass
		}
	}

	// Use default level
	for _, errorMode := range checkConfiguration.ErrorModes {
		if errorMode == checkmode.Default {
			return Error
		}
	}

	for _, warningMode := range checkConfiguration.WarningModes {
		if warningMode == checkmode.Default {
			return Warning
		}
	}

	for _, infoMode := range checkConfiguration.InfoModes {
		if infoMode == checkmode.Default {
			return Info
		}
	}

	for _, passMode := range checkConfiguration.PassModes {
		if passMode == checkmode.Default {
			return Pass
		}
	}

	// TODO: this should return an error
	return Pass
}
