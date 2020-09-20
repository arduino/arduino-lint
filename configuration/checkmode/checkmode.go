package checkmode

import "github.com/arduino/arduino-check/projects/projecttype"

type Type int

const (
	Permissive Type = iota
	LibraryManagerSubmission
	LibraryManagerIndexed
	Official
	Default
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
	}

	// This mode is always enabled
	checkModes[Default] = true

	return checkModes
}
