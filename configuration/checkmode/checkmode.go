package checkmode

import (
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/sirupsen/logrus"
)

type Type int

//go:generate stringer -type=Type -linecomment
const (
	Permissive               Type = iota // --permissive
	LibraryManagerSubmission             // --library-manager=submit
	LibraryManagerIndexed                // --library-manager=update
	Official                             // ARDUINO_CHECK_OFFICIAL
	All                                  // always
	Default                              // default
)

func Modes(defaultCheckModes map[projecttype.Type]map[Type]bool, customCheckModes map[Type]bool, superprojectType projecttype.Type) map[Type]bool {
	checkModes := make(map[Type]bool)

	// Merge the default settings with any custom settings specified by the user
	for key, defaultValue := range defaultCheckModes[superprojectType] {
		customCheckModeValue, customCheckModeIsConfigured := customCheckModes[key]
		if customCheckModeIsConfigured {
			checkModes[key] = customCheckModeValue
		} else {
			checkModes[key] = defaultValue
		}
		logrus.Tracef("Check mode option %s set to %t\n", key.String(), checkModes[key])
	}

	// This mode is always enabled
	checkModes[All] = true

	return checkModes
}
