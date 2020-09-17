package checklevel

import (
	"github.com/arduino/arduino-check/check/checkconfigurations"
	"github.com/arduino/arduino-check/configuration"
)

func CheckLevel(checkConfiguration checkconfigurations.Configuration) Level {
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
