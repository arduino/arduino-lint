// Package checkmode defines the tool configuration options that affect checks.
package checkmode

import (
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/sirupsen/logrus"
)

// Type is the type for check modes.
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
		logrus.Tracef("Check mode option %s set to %t\n", key.String(), checkModes[key])
	}

	// This mode is always enabled
	checkModes[All] = true

	return checkModes
}
