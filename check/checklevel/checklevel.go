package checklevel

import (
	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/configuration"
)

//go:generate stringer -type=Type -linecomment
type Type int

// Line comments set the string for each level
const (
	Info    Type = iota // info
	Warning             // warning
	Error               // error
	Pass                // pass
	Notice              // notice
)

func CheckLevel(checkConfiguration checkconfigurations.Type) Type {
	configurationCheckModes := configuration.CheckModes(checkConfiguration.ProjectType)
	for _, promoteMode := range checkConfiguration.PromoteModes {
		if configurationCheckModes[promoteMode] == true {
			return Error
		}
	}

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

	return Pass
}
